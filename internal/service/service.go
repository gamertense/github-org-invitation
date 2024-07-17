package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"my-app/internal/model"
	"net/http"
	"os"
)

// SendInvitation sends an invitation to join a team using emails from a CSV file
func SendInvitation(orgName, teamName string) error {
	githubToken := os.Getenv("GITHUB_TOKEN")
	githubApiUrl := os.Getenv("GITHUB_API_URL")

	teamsURL := fmt.Sprintf("%s/orgs/%s/teams", githubApiUrl, orgName)

	req, _ := http.NewRequest("GET", teamsURL, nil)
	req.Header.Add("Authorization", "Bearer "+githubToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var teams []model.Team
	json.Unmarshal(body, &teams)

	var teamID int
	for _, team := range teams {
		if team.Name == teamName {
			teamID = team.ID
			break
		}
	}

	filePath := os.Getenv("EMAIL_LIST_PATH")
	emails, err := extractEmailsFromCSV(filePath)
	if err != nil {
		return err
	}
	// Loop through the emails and send an invitation to each one
	for _, email := range emails {
		invitation := model.Invitation{
			Role:    "direct_member",
			TeamIDs: []int{teamID},
			Email:   email,
		}
		// print the invitation
		fmt.Println("Sending invitations to:")
		fmt.Println(invitation)

		jsonData, _ := json.Marshal(invitation)

		invitationsURL := fmt.Sprintf("%s/orgs/%s/invitations", githubApiUrl, orgName)

		req, _ = http.NewRequest("POST", invitationsURL, bytes.NewBuffer(jsonData))
		req.Header.Add("Authorization", "Bearer "+githubToken)
		req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
		req.Header.Add("Accept", "application/vnd.github+json")

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 201 {
			return fmt.Errorf("error sending invitation: %s", resp.Status)
		}
	}

	return nil
}

func extractEmailsFromCSV(filePath string) ([]string, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new reader
	reader := csv.NewReader(file)

	var emails []string
	// Read the CSV file line by line
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// Extract an email from the last column
		emails = append(emails, record[len(record)-1])
	}

	// Remove the first element (header)
	emails = emails[1:]

	return emails, nil
}
