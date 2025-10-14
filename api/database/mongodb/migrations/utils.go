package migrations

import (
	"strings"
)

// isCollectionExistsError checks if the error is due to collection already existing
func isCollectionExistsError(err error) bool {
	if err == nil {
		return false
	}
	
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "collection already exists") || 
		   strings.Contains(errMsg, "already exists")
}
