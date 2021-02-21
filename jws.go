package xgo

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"golang.org/x/oauth2/jws"
)

// Verify verifies the token from given the public key
func Verify(token string, key *rsa.PublicKey) error {
	return jws.Verify(token, key)
}

// Claim is jws claim struct
type Claim struct {
	PrivateKeyID string                 `json:"private_key_id"`
	PrivateKey   string                 `json:"private_key"`
	ClientEmail  string                 `json:"client_email"`
	TokenURI     string                 `json:"token_uri"`
	Audience     string                 `json:"-"`
	CustomClaims map[string]interface{} `json:"-"`
}

// Token gets a jws token string from claim
func (c *Claim) Token(expire time.Duration) (string, error) {
	rsaKey, err := c.privateKey()
	if err != nil {
		return "", err
	}

	iat := time.Now()
	exp := iat.Add(expire)
	jwt := &jws.ClaimSet{
		Iss:           c.ClientEmail,
		Sub:           c.ClientEmail,
		Aud:           c.Audience,
		Iat:           iat.Unix(),
		Exp:           exp.Unix(),
		PrivateClaims: c.CustomClaims,
	}
	jwsHeader := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
		KeyID:     c.PrivateKeyID,
	}

	return jws.Encode(jwsHeader, jwt, rsaKey)
}

// PublicKey gets a public key object from claim
func (c *Claim) PublicKey() (*rsa.PublicKey, error) {
	privateKey, err := c.privateKey()
	if err != nil {
		return nil, err
	}

	rsaPublicKey := privateKey.Public()
	parsed, ok := rsaPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is invalid")
	}
	return parsed, nil
}

func (c *Claim) privateKey() (*rsa.PrivateKey, error) {
	key := []byte(c.PrivateKey)
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("private key should be a PEM or plain PKSC1 or PKCS8; parse error: %v", err)
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is invalid")
	}
	return parsed, nil
}
