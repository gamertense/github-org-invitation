# GitHub Organization Invitation Service

This is a simple Go backend service that sends an invitation to a user for a GitHub organization.

## Prerequisites
- Go 1.22.2
- A GitHub account with a personal access token
    - The token should have the `admin:org` scope
    - The token should be stored in an environment variable named `GITHUB_TOKEN`
- A GitHub organization with at least one team
- A user to invite to the organization
    - The user should have an email address associated with their GitHub account

## How to run

1. Clone the repository
2. Create a `.env` file in the root directory by following the `.env.sample` file.
3. Run `go run main.go` in the root directory.
4. Send a POST request to `http://localhost:8080/send-invitation` with the following JSON payload:
    
    ```json
    {
        "orgName": "<organization>",
        "teamName": "<team>"
    }
    ```

    Replace `<organization>` with the name of the organization and `<team>` with the name of the team you want to invite the user to.