package xgo_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/glassonion1/xgo"
)

func TestJWS(t *testing.T) {
	t.Run("tests claim", func(t *testing.T) {
		path := filepath.Join("testdata/test_jws", "credentials.json")
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("test failed: %+v", err)
		}

		claim := xgo.Claim{
			Audience: "test audience",
		}
		if err := json.Unmarshal(raw, &claim); err != nil {
			t.Errorf("test failed: %+v", err)
		}
		token, err := claim.Token(time.Hour)
		if err != nil {
			t.Errorf("test failed: %+v", err)
		}
		if len(token) == 0 {
			t.Errorf("token size 0: %+v", err)
		}

		fmt.Printf("token: %s", token)

		pk, err := claim.PublicKey()
		if err != nil {
			t.Errorf("test failed: %+v", err)
		}
		err = xgo.VerifyJWSToken(token, pk)
		if err != nil {
			t.Errorf("test failed: %+v", err)
		}
	})
}
