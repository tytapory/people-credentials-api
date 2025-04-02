package genderize

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetGender(name string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Genderize API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Genderize API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Genderize API response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse Genderize API response: %v", err)
	}

	gender, ok := result["gender"].(string)
	if !ok {
		return "", fmt.Errorf("invalid or missing 'gender' in Genderize API response")
	}

	return gender, nil
}
