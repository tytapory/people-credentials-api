package nationalize

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetNationality(name string) (string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Nationalize API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Nationalize API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Nationalize API response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse Nationalize API response: %v", err)
	}

	countries, ok := result["country"].([]interface{})
	if !ok || len(countries) == 0 {
		return "", fmt.Errorf("no valid country data found in Nationalize API response")
	}

	first := countries[0].(map[string]interface{})
	cid, ok := first["country_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid or missing 'country_id' in Nationalize API response")
	}

	return cid, nil
}
