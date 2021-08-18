package handlers

import (
	"database/sql"
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

func (a *API) CreateRoomHandler(c *gin.Context) {
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

	roomID, err := a.dw.InsertRoom(uniqueRoomID.String(), request.Name, passwordHash, ownerID.(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	//add owner as the first member of the room
	err = a.dw.InsertMember(roomID, ownerID.(int))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	a.wsServer.CreateRoom(uniqueRoomID.String(), false)

	c.JSON(http.StatusCreated, gin.H{
		"unique_room_id": uniqueRoomID.String(),
	})
}

func (a *API) DeleteRoomHandler(c *gin.Context) {
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

	err := a.dw.CheckIfUniqueRoomIDExists(uniqueRoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	err = a.dw.ChceckIfUserIsOwnerOfTheRoom(uniqueRoomID, userID.(int))
	if err != nil {
		if err.Error() == "Not the owner" {
			isAdmin, exists := c.Get("is_admin")
			if !exists {
				c.Status(http.StatusInternalServerError)
				return
			}

			if isAdmin == false {
				c.Status(http.StatusUnauthorized)
				return
			}
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	err = a.dw.DeleteRoom(uniqueRoomID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	a.wsServer.DeleteRoom(uniqueRoomID)

	c.Status(http.StatusOK)
}

func (a *API) GetRoomMembersHandler(c *gin.Context) {

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

	//check if room with this unique id exists
	err := a.dw.CheckIfUniqueRoomIDExists(uniqueRoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	//check if user sending this request is a member of this room
	err = a.dw.CheckIfUserIsAMemberOfASpecificRoom(uniqueRoomID, userID.(int))
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusConflict)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	//get members of this room
	members, err := a.dw.GetRoomMembers(uniqueRoomID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usernames": members,
	})
}

func (a *API) UpdateRoomDataHandler(c *gin.Context) {

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

	err := a.dw.ChceckIfUserIsOwnerOfTheRoom(uniqueRoomID, userID.(int))
	if err != nil {
		if err.Error() == "Not the owner" {
			isAdmin, exists := c.Get("is_admin")
			if !exists {
				c.Status(http.StatusInternalServerError)
				return
			}

			if isAdmin == false {
				c.Status(http.StatusUnauthorized)
				return
			}
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	//get queries
	newName := c.Query("name")
	newPassword1 := c.Query("password1")
	newPassword2 := c.Query("password2")

	if newName == "" && newPassword1 == "" && newPassword2 == "" {
		c.Status(http.StatusNotModified)
		return
	}

	if newName != "" {
		err := utils.ValidateRoomName(newName)
		if err != nil {
			c.Status(http.StatusNotAcceptable)
			return
		}

		err = a.dw.UpdateRoomName(newName, uniqueRoomID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		if newPassword1 != "" && newPassword2 != "" {
			c.Status(http.StatusOK)
			return
		}
	}

	if newPassword1 != "" && newPassword2 != "" {
		if newPassword1 != newPassword2 {
			c.Status(http.StatusConflict)
			return
		}

		err = utils.ValidatePassword(newPassword1)
		if err != nil {
			c.Status(http.StatusNotAcceptable)
			return
		}

		hashedPassword, err := utils.HashPassword(newPassword1)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		err = a.dw.UpdateRoomPassword(hashedPassword, uniqueRoomID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusNotModified)
	return

}

type MemberRequest struct {
	Username string `json:"username"`
}

func (a *API) AddMemberToTheRoomHandler(c *gin.Context) {
	var request MemberRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

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

	err := a.dw.ChceckIfUserIsOwnerOfTheRoom(uniqueRoomID, userID.(int))
	if err != nil {
		if err.Error() == "Not the owner" {
			isAdmin, exists := c.Get("is_admin")
			if !exists {
				c.Status(http.StatusInternalServerError)
				return
			}

			if isAdmin == false {
				c.Status(http.StatusUnauthorized)
				return
			}
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	//get users id
	userToAddID, err := a.dw.GetUserIDByUsername(request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	//get room id by uniqueRoomID
	roomID, err := a.dw.GetRoomIDByUniqueRoomID(uniqueRoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	err = a.dw.InsertMember(roomID, userToAddID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	return
}

func (a *API) RemoveMemberFromTheRoomHandler(c *gin.Context) {
	var request MemberRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

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

	username, exists := c.Get("username")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	if username == request.Username {
		c.Status(http.StatusForbidden)
		return
	}

	err := a.dw.ChceckIfUserIsOwnerOfTheRoom(uniqueRoomID, userID.(int))
	if err != nil {
		if err.Error() == "Not the owner" {
			isAdmin, exists := c.Get("is_admin")
			if !exists {
				c.Status(http.StatusInternalServerError)
				return
			}

			if isAdmin == false {
				c.Status(http.StatusUnauthorized)
				return
			}
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	userToRemoveID, err := a.dw.GetUserIDByUsername(request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	roomID, err := a.dw.GetRoomIDByUniqueRoomID(uniqueRoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	err = a.dw.CheckIfUserIsAMemberOfASpecificRoom(uniqueRoomID, userToRemoveID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusConflict)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	err = a.dw.DeleteMember(roomID, userToRemoveID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	//room := wsServer.FindRoomByID(uniqueRoomID)

	c.Status(http.StatusOK)
}

func (a *API) LeaveRoomHandler(c *gin.Context) {

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

	roomID, err := a.dw.GetRoomIDByUniqueRoomID(uniqueRoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	err = a.dw.CheckIfUserIsAMemberOfASpecificRoom(uniqueRoomID, userID.(int))
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusConflict)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	//if user that want to leave the room is an admin, delete the whole room and it's members
	err = a.dw.ChceckIfUserIsOwnerOfTheRoom(uniqueRoomID, userID.(int))
	if err != nil {
		if err.Error() == "Not the owner" {
			isAdmin, exists := c.Get("is_admin")
			if !exists {
				c.Status(http.StatusInternalServerError)
				return
			}

			if isAdmin.(bool) {
				err := a.dw.DeleteRoom(uniqueRoomID)
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return
				}

				err = a.dw.DeleteAllRoomMembers(roomID)
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return
				}

				c.Status(http.StatusOK)
				return
			} else {
				err := a.dw.DeleteMember(roomID, userID.(int))
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return
				}
				c.Status(http.StatusOK)
				return
			}
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	err = a.dw.DeleteRoom(uniqueRoomID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	err = a.dw.DeleteAllRoomMembers(roomID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	return
}
