package handler

import (
	"my-app/internal/model"
	"net/http"

	"my-app/internal/service"

	"github.com/gin-gonic/gin"
)

// SendInvitation is an HTTP handler for sending invitations
func SendInvitation(c *gin.Context) {
	var req model.InvitationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := service.SendInvitation(req.OrgName, req.TeamName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent successfully"})
}

// FetchUsernameByEmail is an HTTP handler for fetching a username by email
func FetchUsernameByEmail(c *gin.Context) {
	email := c.Query("email")

	user, err := service.FetchUsernameByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
