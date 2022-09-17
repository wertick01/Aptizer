package authorization

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"aptizer.com/internal/app/db"
	"aptizer.com/internal/app/handlers/middleware"
	"aptizer.com/internal/app/handlers/responses"
	"aptizer.com/internal/app/models"
	"aptizer.com/internal/app/processors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Authoriser struct {
	processor *processors.Processor
}

// NewAuthoriser - Creating new copy of Authoriser srtuct.
func NewAuthoriser(processor *processors.Processor) *Authoriser {
	authoriser := new(Authoriser)
	authoriser.processor = processor
	return authoriser
}

// Login - Login in app.
func (m *Authoriser) Login(c *gin.Context) {
	var creds *models.UserAuth

	if err := c.ShouldBindJSON(&creds); err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DecodeErr, http.StatusBadRequest)
		return
	}

	user, err := m.processor.FindByPhone(creds.UserPhone)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
		return
	}

	if !db.ComparePassword(user.Hash, creds.Password) {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", models.WrongPassword), models.WrongPassword, http.StatusInternalServerError)
		return
	}
	user.Hash = ""

	m.CreateToken(user, c)
	rt, err := GenerateRefreshToken()
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.RefreshTokenError, http.StatusInternalServerError)
		return
	}

	rt.UserID = user.UserID
	rt.UserAgent = c.Request.UserAgent()

	if err = m.processor.SetRefrToken(rt); err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"result":        "OK",
			"userdata":      user,
			"refresh_token": rt.RefreshToken,
		},
	)
}

// CreateToken - Generating JWT token.
func (m *Authoriser) CreateToken(user *models.User, c *gin.Context) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &models.Claims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(middleware.MySigningKey)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.JWTTokenError, http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

var RefrTest string

// Refresh - Finding user, checking and generating refresh-token for him.
func (m *Authoriser) Refresh(c *gin.Context) {
	var refr *models.RefreshToken

	if err := c.ShouldBindJSON(&refr); err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DecodeErr, http.StatusBadRequest)
		return
	}
	refr.UserAgent = c.Request.UserAgent()

	refr, err := m.processor.CheckRefrToken(refr)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
		return
	}

	user, err := m.processor.FindUser(refr.UserID)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
		return
	}

	m.CreateToken(user, c)
	refr, err = GenerateRefreshToken()
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.RefreshTokenError, http.StatusInternalServerError)
		return
	}

	refr.UserID = user.UserID
	refr.UserAgent = c.Request.UserAgent()

	err = m.processor.SetRefrToken(refr)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"result":        "OK",
			"userdata":      user,
			"refresh_token": refr.RefreshToken,
		},
	)
}

// GenerateRefreshToken - Generating new refresh token for user.
func GenerateRefreshToken() (*models.RefreshToken, error) {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	if _, err := r.Read(b); err != nil {
		return nil, err
	}

	tn := time.Now().Add(RefreshTokenTTL)
	rt := &models.RefreshToken{
		RefreshToken: fmt.Sprintf("%x", b),
		ExpiresIn:    tn.Unix(),
	}
	RefrTest = rt.RefreshToken

	return rt, nil
}

const RefreshTokenTTL = 1440 * time.Hour
