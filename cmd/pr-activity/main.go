package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

type reviewEvent struct {
	id         int
	name       string
	createDate time.Time
}

func main() {
	// Define command-line flags
	username := flag.String("username", "", "GitHub username")
	filename := flag.String("filename", "", "CSV filename")
	flag.Parse()

	// Check if the username and filename are provided
	if *username == "" || *filename == "" {
		fmt.Println("Usage: go run main.go -username <username> -filename <filename>")
		return
	}

	// Construct the URL to fetch the user's events
	url := fmt.Sprintf("https://api.github.com/users/%s/events?per_page=100", *username)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Set the user agent header to identify our request
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "bearer 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")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Send the HTTP request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Parse the response body into a JSON object
	var events []map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	// Open the CSV file for writing
	file, err := os.Create(*filename)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	err = writer.Write([]string{"PR Number", "Title", "Created Date", "Merge Date"})
	if err != nil {
		fmt.Println("Error writing CSV header:", err)
		return
	}

	// Loop through the events and find the PR events
	fmt.Printf("Received [%d] events", len(events))
	reviewEvents := make([]*reviewEvent, 0)

	for _, event := range events {
		eventType, ok := event["type"].(string)
		if !ok {
			continue
		}

		if eventType == "PullRequestReviewEvent" {
			prData, ok := event["payload"].(map[string]interface{})["pull_request"].(map[string]interface{})
			if !ok {
				continue
			}

			prTitle, ok := prData["title"].(string)
			if !ok {
				continue
			}
			prNumber, ok := prData["number"].(float64)
			if !ok {
				continue
			}

			createdDateRaw, ok := prData["created_at"].(string)
			if !ok {
				continue
			}

			createdDate, err := time.Parse(time.RFC3339, createdDateRaw)
			if err != nil {
				fmt.Println(err)
			}

			r := &reviewEvent{
				id:         int(prNumber),
				name:       prTitle,
				createDate: createdDate,
			}

			reviewEvents = append(reviewEvents, r)

			fmt.Println("--------------", len(reviewEvents))

			fmt.Printf("Writing %s\n", prTitle)
			// Write the PR data to the CSV file
			/**err = writer.Write([]string{fmt.Sprintf("%d", int(prNumber)), prTitle, createdDate, mergedDate})
			if err != nil {
				fmt.Println("Error writing PR data to CSV:", err)
				return
			}**/
		}
	}

	fmt.Printf("PR activity saved to %s\n", *filename)
}
