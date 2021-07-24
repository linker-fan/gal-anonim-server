package handlers

import (
	"database/sql"
	"linker-fan/gal-anonim-server/server/database"
	"linker-fan/gal-anonim-server/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type CreateRoomRequest struct {
	Name      string `json:"name"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

func CreateRoomHandler(c *gin.Context) {
	var request CreateRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err := utils.ValidateRoomName(request.Name)
	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if request.Password1 != request.Password2 {
		c.Status(http.StatusConflict)
		return
	}

	err = utils.ValidatePassword(request.Password1)
	if err != nil {
		c.Status(http.StatusNotAcceptable)
		return
	}

	passwordHash, err := utils.HashPassword(request.Password1)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	uniqueRoomID := xid.New()
	ownerID, exists := c.Get("id")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	err = database.InsertRoom(uniqueRoomID.String(), request.Name, passwordHash, ownerID.(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"unique_room_id": uniqueRoomID.String(),
	})
}

func DeleteRoomHandler(c *gin.Context) {
	uniqueRoomID := c.Param("uniqueRoomID")
	if uniqueRoomID == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	err := database.CheckIfUniqueRoomIDExists(uniqueRoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	err = database.ChceckIfUserIsOwnerOfTheRoom(uniqueRoomID, userID.(int))
	if err != nil {
		if err.Error() == "Not the owner" {
			isAdmin, exists := c.Get("is_admin")
			if !exists {
				c.Status(http.StatusInternalServerError)
				return
			}

			if isAdmin == false {
				c.Status(http.StatusUnauthorized)
			}
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	err = database.DeleteRoom(uniqueRoomID)
	c.Status(http.StatusOK)
}

func GetRoomMembersHandler(c *gin.Context) {

}

func UpdateRoomDataHandler(c *gin.Context) {

}
