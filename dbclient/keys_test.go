package dbclient_test

import (
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/dbclient"
)

func TestConstants(t *testing.T) {
	t.Run("AttrPK", func(t *testing.T) {
		if dbclient.AttrPK != "pk" {
			t.Errorf("AttrPK = %q, want %q", dbclient.AttrPK, "pk")
		}
	})

	t.Run("AttrSK", func(t *testing.T) {
		if dbclient.AttrSK != "sk" {
			t.Errorf("AttrSK = %q, want %q", dbclient.AttrSK, "sk")
		}
	})

	t.Run("PrefixAccount", func(t *testing.T) {
		if dbclient.PrefixAccount != "ACCOUNT#" {
			t.Errorf("PrefixAccount = %q, want %q", dbclient.PrefixAccount, "ACCOUNT#")
		}
	})

	t.Run("PrefixUser", func(t *testing.T) {
		if dbclient.PrefixUser != "USER#" {
			t.Errorf("PrefixUser = %q, want %q", dbclient.PrefixUser, "USER#")
		}
	})

	t.Run("SKMeta", func(t *testing.T) {
		if dbclient.SKMeta != "META#" {
			t.Errorf("SKMeta = %q, want %q", dbclient.SKMeta, "META#")
		}
	})
}

func TestAccountPK(t *testing.T) {
	tests := []struct {
		accountID string
		want      string
	}{
		{"abc", "ACCOUNT#abc"},
		{"123", "ACCOUNT#123"},
		{"", "ACCOUNT#"},
	}
	for _, tt := range tests {
		t.Run(tt.accountID, func(t *testing.T) {
			got := dbclient.AccountPK(tt.accountID)
			if got != tt.want {
				t.Errorf("AccountPK(%q) = %q, want %q", tt.accountID, got, tt.want)
			}
		})
	}
}

func TestUserPK(t *testing.T) {
	tests := []struct {
		userID string
		want   string
	}{
		{"xyz", "USER#xyz"},
		{"456", "USER#456"},
		{"", "USER#"},
	}
	for _, tt := range tests {
		t.Run(tt.userID, func(t *testing.T) {
			got := dbclient.UserPK(tt.userID)
			if got != tt.want {
				t.Errorf("UserPK(%q) = %q, want %q", tt.userID, got, tt.want)
			}
		})
	}
}
