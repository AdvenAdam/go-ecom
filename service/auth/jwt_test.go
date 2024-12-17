package auth

import "testing"

func TestCreateJWTToken(t *testing.T) {
	t.Run("fail if cant create token", func(t *testing.T) {
		secret := []byte("secret")
		token, err := CreateJWTToken(secret, 1)
		if err != nil {
			t.Errorf(" failed to create token: %v", err)
		}
		if token == "" {
			t.Error("token is empty")
		}
	})
}
