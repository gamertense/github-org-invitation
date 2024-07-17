package model

type Team struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Invitation struct {
	Role    string `json:"role"`
	TeamIDs []int  `json:"teamIds"`
	Email   string `json:"email"`
}
