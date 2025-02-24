package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testKey = []byte("12345678901234567890123456789012") // 32バイト

func TestNewEncrypter(t *testing.T) {
	tests := []struct {
		name        string
		version     Version
		expectError bool
	}{
		{
			name:        "有効なバージョン",
			version:     Version1,
			expectError: false,
		},
		{
			name:        "無効なバージョン",
			version:     "v999",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewEncrypter(tt.version, testKey)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, enc)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, enc)
			}
		})
	}
}

func TestCBCEncrypter_String(t *testing.T) {
	enc, err := NewEncrypter(Version1, testKey)
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		wantSame bool
	}{
		{
			name:     "通常の文字列",
			input:    "Hello, World!",
			wantErr:  false,
			wantSame: true,
		},
		{
			name:     "空文字列",
			input:    "",
			wantErr:  false,
			wantSame: true,
		},
		{
			name:     "日本語文字列",
			input:    "こんにちは世界",
			wantErr:  false,
			wantSame: true,
		},
		{
			name:     "長い文字列",
			input:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			wantErr:  false,
			wantSame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 暗号化
			encrypted, err := enc.EncryptString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, encrypted)

			// 復号化
			decrypted, err := enc.DecryptString(encrypted)
			require.NoError(t, err)

			if tt.wantSame {
				assert.Equal(t, tt.input, decrypted)
			} else {
				assert.NotEqual(t, tt.input, decrypted)
			}
		})
	}
}

func TestCBCEncrypter_Int(t *testing.T) {
	enc, err := NewEncrypter(Version1, testKey)
	require.NoError(t, err)

	tests := []struct {
		name  string
		input int64
	}{
		{
			name:  "正の整数",
			input: 12345,
		},
		{
			name:  "負の整数",
			input: -12345,
		},
		{
			name:  "ゼロ",
			input: 0,
		},
		{
			name:  "最大値",
			input: 9223372036854775807,
		},
		{
			name:  "最小値",
			input: -9223372036854775808,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 暗号化
			encrypted, err := enc.EncryptInt(tt.input)
			require.NoError(t, err)
			assert.NotEmpty(t, encrypted)

			// 復号化
			decrypted, err := enc.DecryptInt(encrypted)
			require.NoError(t, err)
			assert.Equal(t, tt.input, decrypted)
		})
	}
}

func TestCBCEncrypter_InvalidFormat(t *testing.T) {
	enc, err := NewEncrypter(Version1, testKey)
	require.NoError(t, err)

	tests := []struct {
		name      string
		input     string
		wantError error
	}{
		{
			name:      "不正なフォーマット（区切りなし）",
			input:     "invalidformat",
			wantError: ErrInvalidFormat,
		},
		{
			name:      "不正なフォーマット（区切り不足）",
			input:     "v1.data",
			wantError: ErrInvalidFormat,
		},
		{
			name:      "不正なバージョン",
			input:     "v2.ZGF0YQ==.aXY=",
			wantError: ErrUnsupportedVersion,
		},
		{
			name:      "不正なBase64",
			input:     "v1.!!!.aXY=",
			wantError: nil, // エラーの種類は特定しない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := enc.DecryptString(tt.input)
			assert.Error(t, err)
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
			}
		})
	}
}

func TestGetEncrypterFromCiphertext(t *testing.T) {
	// 正常系: 有効な暗号文からEncrypterを取得
	enc, err := NewEncrypter(Version1, testKey)
	require.NoError(t, err)

	encrypted, err := enc.EncryptString("test")
	require.NoError(t, err)

	decrypter, err := GetEncrypterFromCiphertext(encrypted, testKey)
	assert.NoError(t, err)
	assert.NotNil(t, decrypter)

	// 異常系: 不正なフォーマット
	_, err = GetEncrypterFromCiphertext("invalid", testKey)
	assert.ErrorIs(t, err, ErrInvalidFormat)
}
