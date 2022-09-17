package middleware

import (
	"fmt"
	"net/http"

	"aptizer.com/internal/app/handlers/responses"
	"aptizer.com/internal/app/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var MySigningKey = []byte("johenews")

// CheckJWTToken - JWT token validation function.
func CheckJWTToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.ErrorCookieToken, http.StatusUnauthorized)
				return
			}
			responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.ErrorWrongJWTToken, http.StatusBadRequest)
			return
		}

		tknStr := token.Value
		claims := &models.Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return MySigningKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.ErrorParseToken, http.StatusUnauthorized)
				return
			}
			responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.ErrorParseToken, http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.NonValidToken, http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
