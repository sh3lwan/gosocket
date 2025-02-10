package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sh3lwan/gosocket/config"
)


// JWT Header
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// JWT Claims
type Claims struct {
	Sub string `json:"sub"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}

// Verify and Decode JWT Token
// VerifyJWT verifies the JWT and returns the payload if valid
func VerifyJWT(token string) (*Claims, error) {
    secretKey := config.GetEnv("JWT_SECRET")

	// Split the token into header, payload, and signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	headerB64 := parts[0]
	payloadB64 := parts[1]
	signatureB64 := parts[2]

	// Decode the header
	headerJSON, err := base64.RawURLEncoding.DecodeString(headerB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %v", err)
	}

	var header Header
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, fmt.Errorf("failed to parse header: %v", err)
	}

	// Ensure the algorithm is HMAC SHA-256
	if header.Alg != "HS256" {
		return nil, fmt.Errorf("unsupported algorithm: %s", header.Alg)
	}

	// Decode the payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	var payload Claims
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %v", err)
	}

	// Check if the token has expired
	if time.Now().Unix() > payload.Exp {
		return nil, fmt.Errorf("token has expired")
	}

	// Recreate the signature
	signingInput := fmt.Sprintf("%s.%s", headerB64, payloadB64)
    fmt.Println("signingInput", signingInput)
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(signingInput))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
    fmt.Println("expectedSignature", expectedSignature, signatureB64)

	// Compare the signatures
	if signatureB64 != expectedSignature {
		return nil, fmt.Errorf("invalid signature")
	}
    fmt.Println("payload", &payload)

	return &payload, nil
}

func base64URLDecode(data string) ([]byte, error) {
	if len(data)%4 != 0 {
		data += strings.Repeat("=", 4-len(data)%4)
	}
	return base64.URLEncoding.DecodeString(data)
}
// Create HMAC SHA-256 Signature
func createHMAC(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

