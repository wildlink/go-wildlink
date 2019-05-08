package wildlink

import (
	"fmt"
)

type APIError struct {
	ErrorMessage string
}

func (e APIError) Error() string {
	if e.ErrorMessage != "" {
		return fmt.Sprintf("wildlink:  %v", e.ErrorMessage)
	}
	return ""
}

func (e APIError) Empty() bool {
	return e.ErrorMessage == ""
}

func relevantError(httpError error, apiError APIError) error {
	if httpError != nil {
		return httpError
	}
	if apiError.Empty() {
		return nil
	}
	return apiError
}
