package domain

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenerateRandomPassword(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Generate random password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := GenerateRandomPassword()

			// Should not return error
			if err != nil {
				t.Errorf("GenerateRandomPassword() error = %v, want nil", err)
				return
			}

			// Should return 6 character password
			if len(password) != 6 {
				t.Errorf("GenerateRandomPassword() password length = %d, want 6", len(password))
			}

			// Should contain only digits
			for _, char := range password {
				if char < '0' || char > '9' {
					t.Errorf("GenerateRandomPassword() contains non-digit character: %c", char)
				}
			}
		})
	}
}

func TestGenerateRandomPassword_Uniqueness(t *testing.T) {
	passwords := make(map[string]bool)

	for i := 0; i < 100; i++ {
		password, err := GenerateRandomPassword()
		if err != nil {
			t.Fatalf("GenerateRandomPassword() error = %v", err)
		}
		passwords[password] = true
	}

	// Should have generated at least some unique passwords
	// (Very unlikely to get 100 identical 6-digit passwords)
	if len(passwords) < 10 {
		t.Errorf("GenerateRandomPassword() generated too few unique passwords: %d", len(passwords))
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "123456",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "Long password (within bcrypt limit)",
			password: strings.Repeat("a", 70), // bcrypt limit is 72 bytes
			wantErr:  false,
		},
		{
			name:     "Password exceeding bcrypt limit",
			password: strings.Repeat("a", 100), // exceeds 72 bytes
			wantErr:  true,
		},
		{
			name:     "Special characters",
			password: "!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Hash should not be empty
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}

				// Hash should be different from original password
				if hash == tt.password {
					t.Error("HashPassword() returned same as input password")
				}

				// Hash should start with bcrypt identifier
				if !strings.HasPrefix(hash, "$2") {
					t.Error("HashPassword() hash doesn't start with bcrypt identifier")
				}

				// Should be able to verify the hash
				err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
				if err != nil {
					t.Errorf("HashPassword() generated invalid hash: %v", err)
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// Generate a known hash for testing
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to generate hash for testing: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "Wrong password",
			password: "wrongpassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Empty hash",
			password: password,
			hash:     "",
			want:     false,
		},
		{
			name:     "Invalid hash",
			password: password,
			hash:     "invalid_hash",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPasswordHash(tt.password, tt.hash)
			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordWorkflow(t *testing.T) {
	// Test complete password workflow: generate -> hash -> verify
	t.Run("Complete workflow", func(t *testing.T) {
		// Generate random password
		password, err := GenerateRandomPassword()
		if err != nil {
			t.Fatalf("GenerateRandomPassword() error = %v", err)
		}

		// Hash the password
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword() error = %v", err)
		}

		// Verify the password
		if !CheckPasswordHash(password, hash) {
			t.Error("CheckPasswordHash() failed to verify generated and hashed password")
		}

		// Verify wrong password fails
		if CheckPasswordHash("wrongpassword", hash) {
			t.Error("CheckPasswordHash() incorrectly verified wrong password")
		}
	})
}

func BenchmarkGenerateRandomPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenerateRandomPassword()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "testpassword123"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCheckPasswordHash(b *testing.B) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPasswordHash(password, hash)
	}
}
