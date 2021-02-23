package xgo_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/glassonion1/xgo"
)

func TestGCPToken(t *testing.T) {
	t.Run("tests generate gcp token", func(t *testing.T) {
		path := filepath.Join("testdata/test_jws", "credentials.json")
		token, err := xgo.GenGCPAccessToken(path, "test audience")
		if err != nil {
			t.Errorf("test failed: %+v", err)
		}
		if len(token) == 0 {
			t.Errorf("token size 0: %+v", err)
		}

		fmt.Printf("token: %s", token)
	})
}
