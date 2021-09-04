package handlers

import (
	"database/sql"
	"linker-fan/gal-anonim-server/server/auth"
	"linker-fan/gal-anonim-server/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

//Register handler validates json and inserts new user into the database
//@author hyperxpizza
func (a *API) Register(c *gin.Context) {
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

	err = utils.ValidateEmail(request.Email)
	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	//check if username is already taken
	err = a.dw.CheckIfUsernameExists(request.Username)
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

	//insert user into the database
	err = a.dw.InsertUser(request.Username, passwordHash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
	return

}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Login function validates username and password given as json through request body
//If the username and password are valid and matching, sets http Cookie with jwt token
//@author hyperxpizza
func (a *API) Login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	//check if username exists in the database and get users password and id
	id, passwordHash, isAdmin, err := a.dw.GetIDAndPasswordByUsername(request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	//check if password from request matches passowordHash from the database
	err = utils.CompareHashAndPassword(passwordHash, request.Password)
	if err != nil {
		c.Status(http.StatusUnauthorized) // or status conflict?
		return
	}

	//generate jwt token
	tokenString, expTime, err := auth.GenerateJWTToken(request.Username, id, isAdmin)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	//TODO:
	//not sure if set cookie or return tokenString?
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwtToken",
		Expires:  *expTime,
		Value:    tokenString,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

//MeHandler returns data set in the context from jwt token
//@author hyperxpizza
func (a *API) MeHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	isAdmin, exists := c.Get("is_admin")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       id.(int),
		"username": username.(string),
		"is_admin": isAdmin.(bool),
	})
	return
}

//RefreshTokenHandler creates a new token and sets a new, valid cookie
//@author hyperxpizza
func (a *API) RefreshTokenHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	isAdmin, exists := c.Get("is_admin")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	tokenString, expTime, err := auth.GenerateJWTToken(username.(string), id.(int), isAdmin.(bool))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	//TODO:
	//not sure if set cookie or return tokenString?
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwtToken",
		Expires:  *expTime,
		Value:    tokenString,
		Secure:   false,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

type SetPinRequest struct {
	Pin string `json:"pin"`
}

func (a *API) SetPinHandler(c *gin.Context) {
	var request SetPinRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := utils.ValidatePin(request.Pin); err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if err := a.dw.SetPin(request.Pin, id.(int)); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	return
}
