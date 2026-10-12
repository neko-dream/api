package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthProviderName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    AuthProviderName
		expectError bool
	}{
		{
			name:        "GOOGLEを正しく変換できる",
			input:       "GOOGLE",
			expected:    ProviderGoogle,
			expectError: false,
		},
		{
			name:        "googleを正しく変換できる（小文字）",
			input:       "google",
			expected:    ProviderGoogle,
			expectError: false,
		},
		{
			name:        "LINEを正しく変換できる",
			input:       "LINE",
			expected:    ProviderLine,
			expectError: false,
		},
		{
			name:        "lineを正しく変換できる（小文字）",
			input:       "line",
			expected:    ProviderLine,
			expectError: false,
		},
		{
			name:        "DEVを正しく変換できる",
			input:       "DEV",
			expected:    ProviderDEV,
			expectError: false,
		},
		{
			name:        "PASSWORDを正しく変換できる",
			input:       "PASSWORD",
			expected:    ProviderPassword,
			expectError: false,
		},
		{
			name:        "GUESTを正しく変換できる",
			input:       "GUEST",
			expected:    ProviderGuest,
			expectError: false,
		},
		{
			name:        "guestを正しく変換できる（小文字）",
			input:       "guest",
			expected:    ProviderGuest,
			expectError: false,
		},
		{
			name:        "不正なプロバイダー名はエラー",
			input:       "INVALID",
			expected:    "",
			expectError: true,
		},
		{
			name:        "空文字はエラー",
			input:       "",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewAuthProviderName(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, AuthProviderName(""), result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestAuthProviderName_IsGuest(t *testing.T) {
	tests := []struct {
		name     string
		provider AuthProviderName
		expected bool
	}{
		{
			name:     "GUESTはtrue",
			provider: ProviderGuest,
			expected: true,
		},
		{
			name:     "GOOGLEはfalse",
			provider: ProviderGoogle,
			expected: false,
		},
		{
			name:     "LINEはfalse",
			provider: ProviderLine,
			expected: false,
		},
		{
			name:     "PASSWORDはfalse",
			provider: ProviderPassword,
			expected: false,
		},
		{
			name:     "DEVはfalse",
			provider: ProviderDEV,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.provider.IsGuest())
		})
	}
}

func TestAuthProviderName_IssuerURI(t *testing.T) {
	tests := []struct {
		name     string
		provider AuthProviderName
		expected IssuerURI
	}{
		{
			name:     "GOOGLEのIssuerURI",
			provider: ProviderGoogle,
			expected: GoogleIssuerURI,
		},
		{
			name:     "LINEのIssuerURI",
			provider: ProviderLine,
			expected: LineIssuerURI,
		},
		{
			name:     "DEVは空のIssuerURI",
			provider: ProviderDEV,
			expected: "",
		},
		{
			name:     "PASSWORDは空のIssuerURI",
			provider: ProviderPassword,
			expected: "",
		},
		{
			name:     "GUESTは空のIssuerURI",
			provider: ProviderGuest,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.provider.IssuerURI())
		})
	}
}

func TestAuthProviderName_String(t *testing.T) {
	tests := []struct {
		name     string
		provider AuthProviderName
		expected string
	}{
		{
			name:     "GOOGLE",
			provider: ProviderGoogle,
			expected: "GOOGLE",
		},
		{
			name:     "LINE",
			provider: ProviderLine,
			expected: "LINE",
		},
		{
			name:     "DEV",
			provider: ProviderDEV,
			expected: "DEV",
		},
		{
			name:     "PASSWORD",
			provider: ProviderPassword,
			expected: "PASSWORD",
		},
		{
			name:     "GUEST",
			provider: ProviderGuest,
			expected: "GUEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.provider.String())
		})
	}
}

func TestIssuerURI_String(t *testing.T) {
	tests := []struct {
		name     string
		issuer   IssuerURI
		expected string
	}{
		{
			name:     "GoogleIssuerURI",
			issuer:   GoogleIssuerURI,
			expected: "https://accounts.google.com",
		},
		{
			name:     "LineIssuerURI",
			issuer:   LineIssuerURI,
			expected: "https://access.line.me",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.issuer.String())
		})
	}
}
