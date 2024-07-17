package model

// InvitationRequest represents the expected structure of the request body
type InvitationRequest struct {
	OrgName  string `json:"orgName" example:"DecaturMakers"`
	TeamName string `json:"teamName" example:"Administrators"`
}
