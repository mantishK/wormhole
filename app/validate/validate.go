package validate

import (
	"errors"
	"net/url"
	"strings"
)

//checks if the data contains the keys. Used for POST, PUT and DELETE requests
func RequiredData(data map[string]interface{}, keys []string) (int, error) {
	if data == nil {
		return 0, errors.New("data not set")
	}
	for count, value := range keys {
		_, ok := data[value]
		if !ok {
			return count, errors.New(value + " required ")
		}
	}
	return -1, nil
}

func RequiredParams(data url.Values, keys []string) (int, error) {
	if data == nil {
		return 0, errors.New("data not set")
	}
	for count, value := range keys {
		_, ok := data[value]
		if !ok {
			return count, errors.New(value + " required ")
		}
	}
	return -1, nil
}

func Password(password string) bool {
	if len(password) < 8 {
		return false
	}
	if !strings.ContainsAny(password, "1234567890") {
		return false
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return false
	}
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}
	if strings.ContainsAny(password, " ") {
		return false
	}
	return true
}

func UserName(username string) bool {
	if strings.Contains(username, " ") {
		return false
	}
	return true
}

//Validate URLs
func Url(url string) bool {
	return true
}
