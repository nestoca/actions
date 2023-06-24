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

type PageResponse struct {
	Version Version `json:"version"`
}

func Publish(opts Opts) error {
	currentVersion, err := getCurrentVersion(opts)
	if err != nil {
		return fmt.Errorf("getting current version: %w", err)
	}
	return updatePage(opts, currentVersion+1)
}

func getCurrentVersion(opts Opts) (int, error) {
	body := ""
	response, err := sendRequest(opts, "GET", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return 0, fmt.Errorf("sending GET request: %w", err)
	}

	var pageResponse PageResponse
	err = json.NewDecoder(bytes.NewReader(response)).Decode(&pageResponse)
	if err != nil {
		return 0, err
	}
	return pageResponse.Version.Number, nil
}

func updatePage(opts Opts, newVersion int) error {
	file, err := os.ReadFile(opts.PageFile)
	if err != nil {
		return fmt.Errorf("reading page file: %w", err)
	}

	updateRequest := UpdatePageRequest{
		Version: Version{Number: newVersion},
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

	// Update page
	_, err = sendRequest(opts, "PUT", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	fmt.Println("Page updated successfully.")
	return nil
}

func sendRequest(opts Opts, method string, bodyReader *bytes.Buffer) ([]byte, error) {
	// Create HTTP request
	endpoint := fmt.Sprintf("%s/wiki/rest/api/content/%s", opts.BaseUrl, opts.PageID)
	req, err := http.NewRequest(method, endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(opts.Username, opts.ApiToken)

	// Send request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer response.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code: %d\n%s", response.StatusCode, responseBody)
	}
	return responseBody, nil
}
