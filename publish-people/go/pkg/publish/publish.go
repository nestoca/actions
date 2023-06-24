package publish

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Opts struct {
	BaseUrl   string
	PageID    string
	SpaceKey  string
	Username  string
	ApiToken  string
	PageTitle string
	PageFile  string
}

type UpdatePageRequest struct {
	Version Version `json:"version"`
	Title   string  `json:"title"`
	Type    string  `json:"type"`
	Body    Body    `json:"body"`
}

type Version struct {
	Number int `json:"number"`
}

type Body struct {
	Storage Storage `json:"storage"`
}

type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

func Publish(opts Opts) error {
	file, err := os.ReadFile(opts.PageFile)
	if err != nil {
		return fmt.Errorf("reading page file: %w", err)
	}

	updateRequest := UpdatePageRequest{
		Version: Version{Number: 2}, // Increase the version number to update the page
		Title:   opts.PageTitle,
		Type:    "page",
		Body: Body{
			Storage: Storage{
				Value:          string(file),
				Representation: "storage",
			},
		},
	}
	requestBody, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("marshalling update request: %w", err)
	}

	// Create the HTTP request
	endpoint := fmt.Sprintf("%s/wiki/rest/api/content/%s", opts.BaseUrl, opts.PageID)
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(opts.Username, opts.ApiToken)

	// Send request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer response.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	if response.StatusCode == http.StatusOK {
		fmt.Println("Page updated successfully.")
	} else {
		return fmt.Errorf("unexpected response code: %d\n%s", response.StatusCode, responseBody)
	}
	return nil
}
