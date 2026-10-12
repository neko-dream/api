package crypto

import (
	"context"
	"strings"
	"testing"

	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testConfig = config.Config{
	ENCRYPTION_VERSION: Version1,
	ENCRYPTION_SECRET:  "12345678901234567890123456789012", // 32バイト
}

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name        string
		config      config.Config
		expectError bool
	}{
		{
			name:        "有効なバージョン",
			config:      testConfig,
			expectError: false,
		},
		{
			name: "無効なバージョン",
			config: config.Config{
				ENCRYPTION_VERSION: "v999",
				ENCRYPTION_SECRET:  testConfig.ENCRYPTION_SECRET,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewEncryptor(&tt.config)
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

func TestGCMEncryptor_String(t *testing.T) {
	enc, err := NewEncryptor(&testConfig)
	require.NoError(t, err)
	ctx := context.Background()

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
			encrypted, err := enc.EncryptString(ctx, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, encrypted)

			// 復号化
			decrypted, err := enc.DecryptString(ctx, encrypted)
			require.NoError(t, err)

			if tt.wantSame {
				assert.Equal(t, tt.input, decrypted)
			} else {
				assert.NotEqual(t, tt.input, decrypted)
			}
		})
	}
}

func TestCBCEncryptor_Int(t *testing.T) {
	enc, err := NewEncryptor(&testConfig)
	require.NoError(t, err)
	ctx := context.Background()

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
			encrypted, err := enc.EncryptInt(ctx, tt.input)
			require.NoError(t, err)
			assert.NotEmpty(t, encrypted)

			// 復号化
			decrypted, err := enc.DecryptInt(ctx, encrypted)
			require.NoError(t, err)
			assert.Equal(t, tt.input, decrypted)
		})
	}
}

func TestCBCEncryptor_InvalidFormat(t *testing.T) {
	enc, err := NewEncryptor(&testConfig)
	require.NoError(t, err)
	ctx := context.Background()

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
			_, err := enc.DecryptString(ctx, tt.input)
			assert.Error(t, err)
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
			}
		})
	}
}

func TestGetEncryptorFromCiphertext(t *testing.T) {
	// 正常系: 有効な暗号文からEncryptorを取得
	enc, err := NewEncryptor(&testConfig)
	require.NoError(t, err)

	encrypted, err := enc.EncryptString(context.Background(), "test")
	require.NoError(t, err)

	decrypter, err := GetEncryptorFromCiphertext(encrypted, []byte(testConfig.ENCRYPTION_SECRET))
	assert.NoError(t, err)
	assert.NotNil(t, decrypter)

	// 異常系: 不正なフォーマット
	_, err = GetEncryptorFromCiphertext("invalid", []byte(testConfig.ENCRYPTION_SECRET))
	assert.ErrorIs(t, err, ErrInvalidFormat)
}

func BenchmarkEncryptor_String(b *testing.B) {
	enc, err := NewEncryptor(&testConfig)
	require.NoError(b, err)
	ctx := context.Background()

	benchmarks := []struct {
		name  string
		input string
	}{
		{"小さい文字列(10B)", "0123456789"},
		{"中サイズ文字列(100B)", strings.Repeat("0123456789", 10)},
		{"大きい文字列(1KB)", strings.Repeat("0123456789", 100)},
		{"とても大きい文字列(10MB)", strings.Repeat("0123456789", 1000000)},
	}

	for _, bm := range benchmarks {
		b.Run("暗号化_"+bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := enc.EncryptString(ctx, bm.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// 復号のベンチマーク用に事前に暗号化しておく
		encrypted, err := enc.EncryptString(ctx, bm.input)
		require.NoError(b, err)

		b.Run("復号_"+bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := enc.DecryptString(ctx, encrypted)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkEncryptor_Int(b *testing.B) {
	enc, err := NewEncryptor(&testConfig)
	require.NoError(b, err)
	ctx := context.Background()

	benchmarks := []struct {
		name  string
		input int64
	}{
		{"小さい整数", 42},
		{"中サイズ整数", 1234567890},
		{"大きい整数", 9223372036854775807},
		{"負の整数", -9223372036854775808},
	}

	for _, bm := range benchmarks {
		b.Run("暗号化_"+bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := enc.EncryptInt(ctx, bm.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		encrypted, err := enc.EncryptInt(ctx, bm.input)
		require.NoError(b, err)

		b.Run("復号_"+bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := enc.DecryptInt(ctx, encrypted)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
