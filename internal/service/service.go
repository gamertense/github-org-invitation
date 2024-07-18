package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		log.Println("Sending invitations to:")
		log.Println(invitation)

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

		// If the user has been added to the orgnization, the status code won't be 201. Comment this out to continue the process.
		if resp.StatusCode != 201 {
			return fmt.Errorf("error sending invitation: %s", resp.Status)
		}

		bodyBytes, err := getRespBodyBytes(resp)
		if err != nil {
			return err
		}

		// Get GitHub username from response body.
		var result struct {
			Login string `json:"login"`
		}
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return err
		}
		file, err := os.OpenFile("./data/usernames.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		// Save username to a file in this format: email,username\n
		if _, err := file.WriteString(fmt.Sprintf("%s,%s\n", email, result.Login)); err != nil {
			return err
		}
	}

	return nil
}

// Get body bytes from the response and print it
func getRespBodyBytes(resp *http.Response) ([]byte, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	log.Println("response body:")
	log.Println(string(bodyBytes))

	return bodyBytes, nil
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

// FetchUsernameByEmail fetches a GitHub username by email
func FetchUsernameByEmail(email string) (string, error) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	githubApiUrl := os.Getenv("GITHUB_API_URL")

	searchURL := fmt.Sprintf("%s/search/users?q=%s", githubApiUrl, email)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+githubToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Assuming a simplified version of parsing the response
	// You'll need to define a struct that matches the GitHub API response format for user search
	var result struct {
		Items []struct {
			Login string `json:"login"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Items) == 0 {
		return "", fmt.Errorf("no GitHub user found with email %s", email)
	}

	// Returning the first matched username
	return result.Items[0].Login, nil
}
