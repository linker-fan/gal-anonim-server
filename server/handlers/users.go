package handlers

import (
	"linker-fan/gal-anonim-server/server/database"
	"linker-fan/gal-anonim-server/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username  string `json:"username"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

//Register handler validates json and inserts new user into the database
//@author hyperxpizza
func Register(c *gin.Context) {
	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	//validate input data
	err := utils.ValidateUsername(request.Username)
	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	//check if username is already taken
	err = database.CheckIfUsernameExists(request.Username)
	if err != nil {
		if err.Error() == "Username already taken" {
			c.Status(http.StatusConflict)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	//check if passwords are the same
	if request.Password1 != request.Password2 {
		c.Status(http.StatusNotAcceptable)
		return
	}

	//validate password
	err = utils.ValidatePassword(request.Password1)
	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	//generate password hash
	passwordHash, err := utils.HashPassword(request.Password1)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	err = database.InsertUser(request.Username, passwordHash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
	return

}

func Login(c *gin.Context) {

}
