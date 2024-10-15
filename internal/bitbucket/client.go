package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const bitbucketAPIURL = "https://api.bitbucket.org/2.0/repositories/"

type Repository struct {
	Name  string `json:"name"`
	Links struct {
		Clone []struct {
			Href string `json:"href"`
		} `json:"clone"`
	} `json:"links"`
}

type Response struct {
	Values []Repository `json:"values"`
}

type Client struct {
	apiToken string
}

func NewClient(apiToken string) *Client {
	return &Client{apiToken: apiToken}
}

func (c *Client) FetchRepositories(projectKey string) ([]Repository, error) {
	url := fmt.Sprintf("%s%s", bitbucketAPIURL, projectKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repositories: %s", resp.Status)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Values, nil
}
