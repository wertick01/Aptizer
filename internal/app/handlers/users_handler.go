package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"aptizer.com/internal/app/handlers/authorization"
	"aptizer.com/internal/app/handlers/responses"
	"aptizer.com/internal/app/models"
	"github.com/gin-gonic/gin"
)

// Create - Launches the processes of authorization check, request unmarshalling,
// user creation, generates tokens and response for him.
func (handler *Handler) Create(c *gin.Context) {
	var newUser *models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DecodeErr, http.StatusBadRequest)
		return
	}

	user, err := handler.Processor.CreateUser(newUser)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
		return
	}

	handler.Authorizer.CreateToken(user, c)
	rt, err := authorization.GenerateRefreshToken()
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.RefreshTokenError, http.StatusInternalServerError)
		return
	}

	rt.UserAgent = c.Request.UserAgent()
	rt.UserID = user.UserID

	if err = handler.Processor.SetRefrToken(rt); err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"result": "OK",
			"data":   user,
		},
	)
}

// List - Launches the processes of authorization check, gets a list of users
// and generates a response for the admin.
func (handler *Handler) List(c *gin.Context) {
	list, err := handler.Processor.ListUsers()

	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"result": "OK",
			"data":   list,
		},
	)
}

// Find - Launches the processes of authorization check, returns the user
// by id (userid) and generates a response for the admin.
func (handler *Handler) Find(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.ErrorParsingID, http.StatusBadRequest)
		return
	}

	user, err := handler.Processor.FindUser(id)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"result":   "OK",
			"userdata": user,
		},
	)
}

// FindByPhone - Launches the processes of authorization check, request unmarshalling,
// returns the user by phone number (userphone) and generates a response for the admin/user.
func (handler *Handler) FindByPhone(c *gin.Context) {
	var wantedUser *models.User

	if err := c.ShouldBindJSON(&wantedUser); err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DecodeErr, http.StatusBadRequest)
		return
	}

	user, err := handler.Processor.FindByPhone(wantedUser.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
			return
		}
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.UserNotFound, http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"result":   "OK",
			"userdata": user,
		},
	)
}

// Change - Launches the processes of authorization check, request unmarshalling,
// making changes to the user's data and generating a response for him.
func (handler *Handler) Change(c *gin.Context) {
	var changeUser *models.User

	if err := c.ShouldBindJSON(&changeUser); err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DecodeErr, http.StatusBadRequest)
		return
	}

	user, err := handler.Processor.UpdateUser(changeUser)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"result":   "OK",
			"userdata": user,
		},
	)
}

// Delete - Launches the processes of authorization check, deletes user's data
// from database by id (userid) and generating a response for him.
func (handler *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.ErrorParsingID, http.StatusBadRequest)
		return
	}

	deleteduser, err := handler.Processor.DeleteUser(id)
	if err != nil {
		responses.WrapGinErrorWithStatus(c, fmt.Errorf("signed string error: %s", err), models.DatabaseQueryError, http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"result":   "OK",
			"userdata": deleteduser,
		},
	)
}
