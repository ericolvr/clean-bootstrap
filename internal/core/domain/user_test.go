package domain

import (
	"testing"
)

func TestUser_ValidateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid user",
			user: &User{
				Name:     "John Doe",
				Mobile:   "11999999999",
				UserType: UserTypeUser,
				Password: "123456",
				Status:   true,
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			user: &User{
				Name:     "",
				Mobile:   "11999999999",
				UserType: UserTypeUser,
				Password: "123456",
				Status:   true,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "Empty mobile",
			user: &User{
				Name:     "John Doe",
				Mobile:   "",
				UserType: UserTypeUser,
				Password: "123456",
				Status:   true,
			},
			wantErr: true,
			errMsg:  "mobile is required",
		},
		{
			name: "Invalid user type - zero",
			user: &User{
				Name:     "John Doe",
				Mobile:   "11999999999",
				UserType: 0,
				Password: "123456",
				Status:   true,
			},
			wantErr: true,
			errMsg:  "invalid user type",
		},
		{
			name: "Invalid user type - negative",
			user: &User{
				Name:     "John Doe",
				Mobile:   "11999999999",
				UserType: -1,
				Password: "123456",
				Status:   true,
			},
			wantErr: true,
			errMsg:  "invalid user type",
		},
		{
			name: "Empty password",
			user: &User{
				Name:     "John Doe",
				Mobile:   "11999999999",
				UserType: UserTypeUser,
				Password: "",
				Status:   true,
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "Valid master user",
			user: &User{
				Name:     "Master User",
				Mobile:   "11888888888",
				UserType: UserTypeMaster,
				Password: "masterpass",
				Status:   true,
			},
			wantErr: false,
		},
		{
			name: "Valid admin user",
			user: &User{
				Name:     "Admin User",
				Mobile:   "11777777777",
				UserType: UserTypeAdmin,
				Password: "adminpass",
				Status:   true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.ValidateUser()

			if (err != nil) != tt.wantErr {
				t.Errorf("User.ValidateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("User.ValidateUser() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUser_ValidateUserUpdate(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid user update",
			user: &User{
				ID:       1,
				Name:     "John Doe Updated",
				Mobile:   "11999999999",
				UserType: UserTypeUser,
				Status:   true,
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			user: &User{
				ID:       1,
				Name:     "",
				Mobile:   "11999999999",
				UserType: UserTypeUser,
				Status:   true,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "Empty mobile",
			user: &User{
				ID:       1,
				Name:     "John Doe",
				Mobile:   "",
				UserType: UserTypeUser,
				Status:   true,
			},
			wantErr: true,
			errMsg:  "mobile is required",
		},
		{
			name: "Invalid user type - zero",
			user: &User{
				ID:       1,
				Name:     "John Doe",
				Mobile:   "11999999999",
				UserType: 0,
				Status:   true,
			},
			wantErr: true,
			errMsg:  "invalid user type",
		},
		{
			name: "Invalid user type - negative",
			user: &User{
				ID:       1,
				Name:     "John Doe",
				Mobile:   "11999999999",
				UserType: -1,
				Status:   true,
			},
			wantErr: true,
			errMsg:  "invalid user type",
		},
		{
			name: "Valid update without password (password not required for updates)",
			user: &User{
				ID:       1,
				Name:     "John Doe",
				Mobile:   "11999999999",
				UserType: UserTypeUser,
				Password: "", // Password not required for updates
				Status:   true,
			},
			wantErr: false,
		},
		{
			name: "Update with all valid user types",
			user: &User{
				ID:       1,
				Name:     "Admin User",
				Mobile:   "11888888888",
				UserType: UserTypeAdmin,
				Status:   false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.ValidateUserUpdate()

			if (err != nil) != tt.wantErr {
				t.Errorf("User.ValidateUserUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("User.ValidateUserUpdate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUserType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		userType UserType
		expected UserType
	}{
		{
			name:     "UserTypeMaster",
			userType: UserTypeMaster,
			expected: 1,
		},
		{
			name:     "UserTypeAdmin",
			userType: UserTypeAdmin,
			expected: 1,
		},
		{
			name:     "UserTypeUser",
			userType: UserTypeUser,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.userType != tt.expected {
				t.Errorf("UserType %s = %v, want %v", tt.name, tt.userType, tt.expected)
			}
		})
	}
}

func TestUser_ValidationEdgeCases(t *testing.T) {
	t.Run("Whitespace only name", func(t *testing.T) {
		user := &User{
			Name:     "   ",
			Mobile:   "11999999999",
			UserType: UserTypeUser,
			Password: "123456",
		}

		// Current implementation doesn't trim whitespace, so this should pass
		// If you want to add trimming, you'd need to update the validation logic
		err := user.ValidateUser()
		if err != nil {
			t.Errorf("User.ValidateUser() with whitespace name error = %v, want nil", err)
		}
	})

	t.Run("Very long name", func(t *testing.T) {
		longName := make([]byte, 1000)
		for i := range longName {
			longName[i] = 'a'
		}

		user := &User{
			Name:     string(longName),
			Mobile:   "11999999999",
			UserType: UserTypeUser,
			Password: "123456",
		}

		err := user.ValidateUser()
		if err != nil {
			t.Errorf("User.ValidateUser() with long name error = %v, want nil", err)
		}
	})

	t.Run("Very long mobile", func(t *testing.T) {
		longMobile := make([]byte, 100)
		for i := range longMobile {
			longMobile[i] = '9'
		}

		user := &User{
			Name:     "John Doe",
			Mobile:   string(longMobile),
			UserType: UserTypeUser,
			Password: "123456",
		}

		err := user.ValidateUser()
		if err != nil {
			t.Errorf("User.ValidateUser() with long mobile error = %v, want nil", err)
		}
	})
}

func TestUser_StructFields(t *testing.T) {
	t.Run("User struct has all required fields", func(t *testing.T) {
		user := &User{
			ID:       123,
			Name:     "Test User",
			Mobile:   "11999999999",
			UserType: UserTypeUser,
			Password: "password123",
			Status:   true,
		}

		// Test that all fields are accessible
		if user.ID != 123 {
			t.Errorf("User.ID = %v, want 123", user.ID)
		}
		if user.Name != "Test User" {
			t.Errorf("User.Name = %v, want 'Test User'", user.Name)
		}
		if user.Mobile != "11999999999" {
			t.Errorf("User.Mobile = %v, want '11999999999'", user.Mobile)
		}
		if user.UserType != UserTypeUser {
			t.Errorf("User.UserType = %v, want %v", user.UserType, UserTypeUser)
		}
		if user.Password != "password123" {
			t.Errorf("User.Password = %v, want 'password123'", user.Password)
		}
		if user.Status != true {
			t.Errorf("User.Status = %v, want true", user.Status)
		}
	})
}
