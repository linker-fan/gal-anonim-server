package handlers

import (
	"bytes"
	"fmt"
	"io"
	"linker-fan/gal-anonim-server/server/hub"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *API) FileUploadHandler(c *gin.Context) {
	uniqueRoomID := c.Param("uniqueRoomID")
	if uniqueRoomID == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	//find client
	client := a.wsServer.FindClientByID(id.(int))
	if client == nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	room := a.wsServer.FindRoomByID(uniqueRoomID)
	if room == nil {
		c.Status(http.StatusBadRequest)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	byteData := &bytes.Buffer{}
	n, err := io.Copy(byteData, io.LimitReader(file, a.fs.MaxUploadSize+10))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if n > a.fs.MaxUploadSize {
		c.Status(http.StatusNotAcceptable)
		return
	}

	err = a.fs.CreateBucketIfDoesNotExist(uniqueRoomID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	err = a.fs.StoreFile(uniqueRoomID, byteData, header.Filename)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	m := hub.Message{
		Action:  hub.SendMessageAction,
		Message: fmt.Sprintf("%s/%s", uniqueRoomID, header.Filename),
		Target:  room,
		Sender:  client,
	}

	room.Broadcast <- &m

	c.JSON(http.StatusOK, gin.H{
		"url": fmt.Sprintf("%s/%s", uniqueRoomID, header.Filename),
	})
	return
}

func (a *API) DownloadFileHandler(c *gin.Context) {

}
