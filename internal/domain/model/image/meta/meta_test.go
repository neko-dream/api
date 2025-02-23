package meta

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/h2non/filetype/types"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// ValidExtensionメソッドのテスト
func TestValidExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		formats  []string
		ext      string
		expected bool
	}{
		{"Valid JPEG", []string{"jpeg", "png"}, "image/jpeg", true},
		{"Invalid GIF", []string{"jpeg", "png"}, "image/gif", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rule := ImageValidationRule{
				allowedFormats: tt.formats,
			}
			meta := ImageMeta{
				Extension: types.MIME{Value: tt.ext},
			}
			ctx := context.Background()
			assert.Equal(t, tt.expected, rule.ValidExtension(ctx, meta))
		})
	}
}

// ValidFileSizeメソッドのテスト
func TestValidFileSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		maxSize  int
		size     int
		expected bool
	}{
		{"Valid Size", 4194304, 4000000, true},
		{"Invalid Size", 4194304, 5000000, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rule := ImageValidationRule{
				maxFileSize: tt.maxSize,
			}
			meta := ImageMeta{
				Size: tt.size,
			}
			ctx := context.Background()
			assert.Equal(t, tt.expected, rule.ValidFileSize(ctx, meta))
		})
	}
}

// ValidBoundsメソッドのテスト
func TestValidBounds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		width    *int
		height   *int
		metaW    int
		metaH    int
		expected bool
	}{
		{"Valid Bounds", lo.ToPtr(300), lo.ToPtr(300), 300, 300, true},
		{"Invalid Width", lo.ToPtr(300), lo.ToPtr(300), 400, 300, false},
		{"Invalid Height", lo.ToPtr(300), lo.ToPtr(300), 300, 400, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rule := ImageValidationRule{
				maxFileSize: 0,
				width:       tt.width,
				height:      tt.height,
			}
			meta := ImageMeta{
				Width:  tt.metaW,
				Height: tt.metaH,
			}
			ctx := context.Background()
			assert.Equal(t, tt.expected, rule.ValidBounds(ctx, meta))
		})
	}
}

// ValidAspectRatioメソッドのテスト
func TestValidAspectRatio(t *testing.T) {
	t.Parallel()

	minAspectRatio := 1.0
	maxAspectRatio := 2.0

	tests := []struct {
		name     string
		minRatio *float64
		maxRatio *float64
		metaW    int
		metaH    int
		expected bool
	}{
		{"Valid Aspect Ratio", &minAspectRatio, &maxAspectRatio, 200, 100, true},
		{"Invalid Aspect Ratio", &minAspectRatio, &maxAspectRatio, 50, 100, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rule := ImageValidationRule{
				minAspectRatio: tt.minRatio,
				maxAspectRatio: tt.maxRatio,
			}
			meta := ImageMeta{
				Width:  tt.metaW,
				Height: tt.metaH,
			}
			ctx := context.Background()
			assert.Equal(t, tt.expected, rule.ValidAspectRatio(ctx, meta))
		})
	}
}

// GetImageSize関数のテスト
func TestGetImageSize(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	file := bytes.NewReader([]byte("12345678901234"))

	size, err := GetImageSize(ctx, file)
	assert.NoError(t, err)
	assert.Equal(t, 14, size)
}

// GetExtension関数のテスト
func TestGetExtension(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	file := bytes.NewReader([]byte("\xff\xd8\xff"))

	ext, err := GetExtension(ctx, file)
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", ext.Value)
}

// GetBounds関数のテスト
func TestGetBounds(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := image.NewRGBA(image.Rect(0, 0, 100, 100))
	file := new(bytes.Buffer)
	err := png.Encode(file, s)
	if err != nil {
		t.Fatal(err)
	}

	width, height, err := GetBounds(ctx, file)
	assert.NoError(t, err)
	assert.Equal(t, 100, width)
	assert.Equal(t, 100, height)
}

// NewImageForProfile関数のテスト
func TestNewImageForProfile(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := shared.NewUUID[user.User]()

	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	for x := 0; x < 100; x++ {
		for y := 0; y < 50; y++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255}) // 青色で塗りつぶし
		}
	}
	f, _ := os.Create("dummy.png")
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	f, _ = os.Open("dummy.png")

	meta, err := NewImageForProfile(ctx, userID, f)
	assert.NoError(t, err)
	assert.NotNil(t, meta)
	assert.Equal(t, "users/"+userID.String()+".jpg", meta.Key)
	assert.Equal(t, 100, meta.Width)
	assert.Equal(t, 50, meta.Height)
	assert.Equal(t, "image/png", meta.Extension.Value)

	os.Remove("dummy.png")
}

// Validateメソッドのテスト
func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		rule     ImageValidationRule
		meta     ImageMeta
		expected []string
	}{
		{
			name: "Valid Image",
			rule: ImageValidationRule{
				allowedFormats: []string{"jpeg", "png"},
				maxFileSize:    4194304,
				width:          lo.ToPtr(300),
				height:         lo.ToPtr(300),
			},
			meta: ImageMeta{
				Extension: types.MIME{Value: "image/png"},
				Size:      4000000,
				Width:     300,
				Height:    300,
			},
			expected: nil,
		},
		{
			name: "Invalid Format",
			rule: ImageValidationRule{
				allowedFormats: []string{"jpeg", "png"},
			},
			meta: ImageMeta{
				Extension: types.MIME{Value: "image/gif"},
			},
			expected: []string{"サポートされていないフォーマットです。"},
		},
		{
			name: "Invalid File Size",
			rule: ImageValidationRule{
				maxFileSize: 4194304,
			},
			meta: ImageMeta{
				Size: 5000000,
			},
			expected: []string{"ファイルサイズが大きすぎます。"},
		},
		{
			name: "Invalid Bounds",
			rule: ImageValidationRule{
				width:  lo.ToPtr(300),
				height: lo.ToPtr(300),
			},
			meta: ImageMeta{
				Width:  400,
				Height: 300,
			},
			expected: []string{"画像のサイズが大きすぎます。"},
		},
		{
			name: "Invalid Aspect Ratio",
			rule: ImageValidationRule{
				minAspectRatio: lo.ToPtr(1.0),
				maxAspectRatio: lo.ToPtr(2.0),
			},
			meta: ImageMeta{
				Width:  50,
				Height: 100,
			},
			expected: []string{"アスペクト比が不正です。"},
		},
		{
			name: "Multiple Errors",
			rule: ImageValidationRule{
				allowedFormats: []string{"jpeg"},
				maxFileSize:    4194304,
				width:          lo.ToPtr(300),
				height:         lo.ToPtr(300),
			},
			meta: ImageMeta{
				Extension: types.MIME{Value: "image/png"},
				Size:      5000000,
				Width:     400,
				Height:    300,
			},
			expected: []string{
				"サポートされていないフォーマットです。",
				"ファイルサイズが大きすぎます。",
				"画像のサイズが大きすぎます。",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			err := tt.meta.Validate(ctx, tt.rule)
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				for _, expectedErr := range tt.expected {
					assert.Contains(t, err.Error(), expectedErr)
				}
			}
		})
	}
}
