package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("test")
	if err != nil {
		t.Error(err)
	}

	if hash == "" {
		t.Error("hash is empty")
	}

	if hash == "test" {
		t.Error("hash is not hashed")
	}
}

func TestComparePassword(t *testing.T) {
	hash, err := HashPassword("test")
	if err != nil {
		t.Errorf("failed to hash password: %v", err)
	}

	if !ComparePassword(hash, []byte("test")) {
		t.Error("expected password to match hash")
	}

	if ComparePassword(hash, []byte("test1")) {
		t.Error("expected password to not match hash")
	}

}
