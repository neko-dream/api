package organization

import (
	"testing"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/stretchr/testify/assert"
)

func TestValidateOrganizationCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr error
	}{
		{
			name:    "正常系: 4文字の英数字",
			code:    "ABC1",
			wantErr: nil,
		},
		{
			name:    "正常系: 8文字の英数字",
			code:    "CODE1234",
			wantErr: nil,
		},
		{
			name:    "正常系: 大文字小文字混在",
			code:    "AbCd123",
			wantErr: nil,
		},
		{
			name:    "異常系: 3文字（短すぎる）",
			code:    "ABC",
			wantErr: messages.OrganizationCodeTooShort,
		},
		{
			name:    "異常系: 空文字",
			code:    "",
			wantErr: messages.OrganizationCodeTooShort,
		},
		{
			name:    "異常系: 日本語を含む",
			code:    "ABC日本",
			wantErr: messages.OrganizationCodeInvalid,
		},
		{
			name:    "異常系: 特殊文字を含む（ハイフン）",
			code:    "ABC-123",
			wantErr: messages.OrganizationCodeInvalid,
		},
		{
			name:    "異常系: 特殊文字を含む（アンダースコア）",
			code:    "ABC_123",
			wantErr: messages.OrganizationCodeInvalid,
		},
		{
			name:    "異常系: スペースを含む",
			code:    "ABC 123",
			wantErr: messages.OrganizationCodeInvalid,
		},
		{
			name:    "境界値: ちょうど4文字",
			code:    "ABCD",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrganizationCode(tt.code)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
