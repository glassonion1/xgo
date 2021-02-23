package xgo

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2/google"
)

// GenGCPAccessToken generates GCP access token from the service account file
func GenGCPAccessToken(saKeyfile, audience string) (string, error) {
	sa, err := ioutil.ReadFile(saKeyfile)
	if err != nil {
		return "", fmt.Errorf("Could not read service account file: %v", err)
	}
	ts, err := google.JWTAccessTokenSourceFromJSON(sa, audience)
	if err != nil {
		return "", fmt.Errorf("Could not parse service account JSON: %v", err)
	}
	token, err := ts.Token()
	if err != nil {
		return "", err
	}
	if !token.Valid() {
		return "", fmt.Errorf("invalid token")
	}
	return token.AccessToken, nil
}
