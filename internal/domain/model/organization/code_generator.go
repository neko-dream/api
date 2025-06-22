package organization

import (
	"slices"
	"unicode"

	"github.com/neko-dream/server/internal/domain/messages"
)

// ValidateOrganizationCode
// 四文字以上の英数字で構成される。
func ValidateOrganizationCode(code string) error {
	if len(code) < 4 {
		return messages.OrganizationCodeTooShort
	}
	for _, r := range code {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !slices.Contains([]rune{'_', '-'}, r) {
			return messages.OrganizationCodeInvalid
		}
		if r > 127 {
			return messages.OrganizationCodeInvalid
		}
	}
	return nil
}
