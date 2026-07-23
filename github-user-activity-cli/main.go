package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Event struct {
	Type  string `json:"type"`
	Actor struct {
		Login string `json:"login"`
	} `json:"actor"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload json.RawMessage `json:"payload"`
	Created string          `json:"created_at"`
}

type PushPayload struct {
	Ref string `json:"ref"`
}

type CreatePayload struct {
	RefType string `json:"ref_type"`
	Ref     string `json:"ref"`
}

type WatchPayload struct {
	Action string `json:"action"`
}

type IssuesPayload struct {
	Action string `json:"action"`
	Issue  struct {
		Number int `json:"number"`
	} `json:"issue"`
}

type PullRequestPayload struct {
	Action      string `json:"action"`
	PullRequest struct {
		Number int `json:"number"`
	} `json:"pull_request"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: github-user-activity-cli <username>")
		os.Exit(1)
	}
	username := os.Args[1]
	resp, err := http.Get("https://api.github.com/users/" + username + "/events")
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			fmt.Println("user not found")
			os.Exit(1)
		} else {
			fmt.Println("internal server error")
			os.Exit(1)
		}
	}

	var events []Event

	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			var payload PushPayload
			if err := json.Unmarshal(event.Payload, &payload); err != nil {
				fmt.Println("error:", err)
				continue
			}
			fmt.Printf("pushed to %s, %s\n", event.Repo.Name, payload.Ref)
		case "CreateEvent":
			var payload CreatePayload
			if err := json.Unmarshal(event.Payload, &payload); err != nil {
				fmt.Println("error decoding create payload:", err)
				continue
			}
			if payload.RefType == "repository" {
				fmt.Printf("Created repository %s\n", event.Repo.Name)
			} else {
				fmt.Printf("Created %s %s in %s\n", payload.RefType, payload.Ref, event.Repo.Name)
			}

		case "WatchEvent":
			fmt.Printf("Starred %s\n", event.Repo.Name)

		case "IssuesEvent":
			var payload IssuesPayload
			if err := json.Unmarshal(event.Payload, &payload); err != nil {
				fmt.Println("error decoding issues payload:", err)
				continue
			}
			fmt.Printf("%s issue #%d in %s\n", payload.Action, payload.Issue.Number, event.Repo.Name)

		case "PullRequestEvent":
			var payload PullRequestPayload
			if err := json.Unmarshal(event.Payload, &payload); err != nil {
				fmt.Println("error decoding pull request payload:", err)
				continue
			}
			fmt.Printf("%s pull request #%d in %s\n", payload.Action, payload.PullRequest.Number, event.Repo.Name)

		default:
			fmt.Printf("%s on %s\n", event.Type, event.Repo.Name)
		}
	}
}
