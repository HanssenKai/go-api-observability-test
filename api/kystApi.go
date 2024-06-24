package api

import (
	"io"
	"net/http"
)

func fetchLocations() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://kystdatahuset.no/ws/api/location/all", nil)
	if err != nil {
		return "", err
	}

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
