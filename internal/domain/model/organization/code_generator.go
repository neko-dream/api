package organization

import (
	"github.com/expr-lang/expr/parser/utils"
	"github.com/neko-dream/server/internal/domain/messages"
)

// ValidateOrganizationCode
// 四文字以上の英数字で構成される。
func ValidateOrganizationCode(code string) error {
	if len(code) < 4 {
		return messages.OrganizationCodeTooShort
	}
	for _, r := range code {
		if !utils.IsAlphaNumeric(r) {
			return messages.OrganizationCodeTooShort
		}
	}
	return nil
}
