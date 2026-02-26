package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"knowledge-capsule/app/models"
)

var secretKey []byte

func InitJWTSecret(secret string) {
	secretKey = []byte(secret)
}

// GenerateJWT creates a simple JWT token manually.
func GenerateJWT(userID, email, role string, expiry time.Duration) (string, error) {
	header := base64.URLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	if role == "" {
		role = models.RoleUser
	}

	claims := models.TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Iat:    time.Now().Unix(),
		Exp:    time.Now().Add(expiry).Unix(),
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payload := base64.URLEncoding.EncodeToString(payloadBytes)

	unsignedToken := header + "." + payload
	signature := signToken(unsignedToken)
	return unsignedToken + "." + signature, nil
}

// VerifyJWT validates and parses a JWT token.
func VerifyJWT(token string) (*models.TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	unsignedToken := parts[0] + "." + parts[1]
	expectedSig := signToken(unsignedToken)

	if !hmac.Equal([]byte(expectedSig), []byte(parts[2])) {
		return nil, errors.New("invalid signature")
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims models.TokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, err
	}

	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// signToken creates HMAC signature for JWT.
func signToken(unsignedToken string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(unsignedToken))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
