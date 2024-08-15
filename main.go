package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Event struct {
	Type      string `bson:"type" json:"type"`
	Repo      Repo `bson:"repo" json:"repo"`
	Payload   Payload `bson:"payload" json:"payload"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Repo struct {
	Name string `bson:"name" json:"name"`
	URL  string `bson:"url" json:"url"`
}

type Payload struct {
	Size int `bson:"size" json:"size"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Username is required as a positional argument.")
	}

	username := os.Args[1]
	
    url := getUrl(username)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", resp.StatusCode)
	}


	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(events) == 0 {
		fmt.Println("No events found.")
		return
	}

	var messages strings.Builder
	for _, event := range events {
		messages.WriteString(formatEventMessage(event) + "\n")
	}

	fmt.Print(messages.String())
}

func getUrl(username string) string {
	return "https://api.github.com/users/" + username + "/events?per_page=5&page=1"
}

func formatEventMessage(event Event) string {
	msg := ""

	switch event.Type {
	case "PushEvent":
		msg = fmt.Sprintf("Pushed %d commits to %s", event.Payload.Size, event.Repo.Name)
	case "CreateEvent":
		msg = fmt.Sprintf("A Git branch or tag is created: %s", event.Repo.Name)
	case "DeleteEvent":
		msg = fmt.Sprintf("A Git branch or tag is deleted: %s", event.Repo.Name)
	case "ForkEvent":
		msg = fmt.Sprintf("A user forks a repository: %s", event.Repo.Name)
	case "GollumEvent":
		msg = fmt.Sprintf("A wiki page is created or updated: %s", event.Repo.Name)
	case "IssueCommentEvent":
		msg = fmt.Sprintf("Activity related to an issue or pull request comment in: %s", event.Repo.Name)
	case "IssuesEvent":
		msg = fmt.Sprintf("An issue is opened or closed in: %s", event.Repo.Name)
	case "PullRequestEvent":
		msg = fmt.Sprintf("A pull request is opened, closed, or merged in: %s", event.Repo.Name)
	case "WatchEvent":
		msg = fmt.Sprintf("A user starred a repository: %s", event.Repo.Name)
	case "ReleaseEvent":
		msg = fmt.Sprintf("A release is published in: %s", event.Repo.Name)
	default:
		msg = fmt.Sprintf("Received an unknown event type: %s", event.Type)
	}

	return msg
}
