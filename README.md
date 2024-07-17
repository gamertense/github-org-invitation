# github-org-invitation

This is a simple Go backend service that sends an invitation to a GitHub organization to a user.

## How to run

1. Clone the repository
2. Run `go run main.go` in the root directory
3. Send a POST request to `http://localhost:8080/send-invitation` with the following JSON payload
    
    ```json
    {
        "orgName": "<organization>",
        "teamName": "<team>",
    }
    ```

    Replace `<organization>` with the name of the organization and `<team>` with the name of the team you want to invite the user to.