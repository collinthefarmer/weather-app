package fetch

import (
	"weather/internal/validation"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func JSON[T validation.Validates](url string, into *T) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error contacting IP-API: %w", err)
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("non-200 status code returned from endpoint: %v", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(body, into); err != nil {
		return fmt.Errorf("error decodings JSON from IP-API response body: %w", err)
	}

	if problems, err := (*into).Validate(); err != nil {
		return fmt.Errorf("error validating response: %s", problems)
	}

	return nil
}
