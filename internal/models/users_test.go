package models

import (
	"snippetbox/internal/assert"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestUserModelExists(t *testing.T) {
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Invalid ID",
			userID: 2,
			want:   false,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDb(t)
			m := UserModel{db}

			exists, err := m.Exists(tt.userID)

			assert.Equal(t, tt.want, exists)
			assert.NilError(t, err)
		})
	}
}
