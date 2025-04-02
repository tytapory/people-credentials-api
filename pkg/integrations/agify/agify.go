package agify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetAge(name string) (int, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("Agify API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Agify API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read Agify API response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse Agify API response: %v", err)
	}

	age, ok := result["age"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid or missing 'age' in Agify API response")
	}

	return int(age), nil
}
