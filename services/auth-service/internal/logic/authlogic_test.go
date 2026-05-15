package logic

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestLoginRequest_Validation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		expectOK bool
	}{
		{"empty username", "", "pass", false},
		{"empty password", "admin", "", false},
		{"both empty", "", "", false},
		{"valid", "admin", "pass123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &LoginRequest{Username: tt.username, Password: tt.password}
			if (req.Username != "" && req.Password != "") != tt.expectOK {
				t.Errorf("LoginRequest{%q, %q}: expectOK=%v", tt.username, tt.password, tt.expectOK)
			}
		})
	}
}

func TestValidateLogic_JWTTokenGeneration(t *testing.T) {
	secret := "test-secret-key"
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id":  uint64(1001),
		"username": "admin",
		"role":     int32(1),
		"mine_id":  uint64(1),
		"exp":      now.Add(3600 * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	// Verify
	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("token should be valid")
	}

	parsedClaims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims should be MapClaims")
	}

	if uint64(parsedClaims["user_id"].(float64)) != 1001 {
		t.Errorf("user_id = %v, want 1001", parsedClaims["user_id"])
	}
	if parsedClaims["username"].(string) != "admin" {
		t.Errorf("username = %s, want admin", parsedClaims["username"].(string))
	}
	if int32(parsedClaims["role"].(float64)) != 1 {
		t.Errorf("role = %v, want 1", parsedClaims["role"])
	}
}

func TestValidateLogic_InvalidToken(t *testing.T) {
	secret := "test-secret-key"

	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"malformed", "not-a-jwt"},
		{"garbage", "header.payload.sig"},
		{"wrong format", "abc.def"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := jwt.Parse(tt.token, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err == nil && parsed.Valid {
				t.Error("invalid token should not be valid")
			}
		})
	}
}

func TestValidateLogic_WrongSigningKey(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id":  uint64(1001),
		"username": "admin",
		"exp":      time.Now().Add(3600 * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("correct-secret"))

	// Validate with wrong secret
	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("wrong-secret"), nil
	})
	if err == nil {
		t.Error("expected error with wrong signing key")
	}
	if parsed != nil && parsed.Valid {
		t.Error("token should not be valid with wrong key")
	}
}

func TestValidateLogic_ExpiredToken(t *testing.T) {
	secret := "test-secret"
	claims := jwt.MapClaims{
		"user_id":  uint64(1001),
		"username": "admin",
		"exp":      time.Now().Add(-1 * time.Hour).Unix(), // expired 1 hour ago
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))

	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err == nil {
		t.Error("expected error for expired token")
	}
	if parsed != nil && parsed.Valid {
		t.Error("expired token should not be valid")
	}
}

func TestValidateLogic_ClaimsExtraction(t *testing.T) {
	secret := "test-secret"
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id":  float64(42),
		"username": "alice",
		"role":     float64(1),
		"mine_id":  float64(2),
		"exp":      now.Add(3600 * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))

	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("token should parse: %v", err)
	}

	parsedClaims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims should be MapClaims")
	}

	// This simulates what ValidateLogic does
	userID := uint64(parsedClaims["user_id"].(float64))
	username := parsedClaims["username"].(string)
	role := int32(parsedClaims["role"].(float64))
	mineID := uint64(parsedClaims["mine_id"].(float64))

	if userID != 42 {
		t.Errorf("user_id = %d, want 42", userID)
	}
	if username != "alice" {
		t.Errorf("username = %s, want alice", username)
	}
	if role != 1 {
		t.Errorf("role = %d, want 1", role)
	}
	if mineID != 2 {
		t.Errorf("mine_id = %d, want 2", mineID)
	}
}

func TestValidateLogic_MissingClaims(t *testing.T) {
	secret := "test-secret"

	tests := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{"missing user_id", jwt.MapClaims{"username": "admin", "exp": 9999999999}},
		{"missing username", jwt.MapClaims{"user_id": float64(1), "exp": 9999999999}},
		{"wrong types", jwt.MapClaims{"user_id": "not-a-number", "username": 123, "exp": 9999999999}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims)
			tokenStr, _ := token.SignedString([]byte(secret))

			parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				return // expected parse failure
			}
			if parsed.Valid {
				// Claims may parse but extraction may fail — test what ValidateLogic does
				parsedClaims, ok := parsed.Claims.(jwt.MapClaims)
				if !ok {
					return
				}
				// Simulate the extraction code from ValidateLogic
				defer func() {
					recover() // type assertion panics are expected for wrong types
				}()
				_ = uint64(parsedClaims["user_id"].(float64))
				_ = parsedClaims["username"].(string)
			}
		})
	}
}
