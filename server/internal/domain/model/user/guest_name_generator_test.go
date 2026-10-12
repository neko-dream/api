package user

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateGuestDisplayID(t *testing.T) {
	t.Run("8桁の数字を生成する", func(t *testing.T) {
		id := GenerateGuestDisplayID()
		assert.Len(t, id, 8)
	})

	t.Run("数字のみで構成される", func(t *testing.T) {
		id := GenerateGuestDisplayID()
		matched, err := regexp.MatchString(`^\d{8}$`, id)
		assert.NoError(t, err)
		assert.True(t, matched, "期待: 8桁の数字, 実際: %s", id)
	})

	t.Run("10000000以上の値を生成する", func(t *testing.T) {
		// 100回生成して全て10000000以上であることを確認
		for i := 0; i < 100; i++ {
			id := GenerateGuestDisplayID()
			assert.GreaterOrEqual(t, id, "10000000", "生成されたID %s は10000000未満", id)
		}
	})

	t.Run("毎回異なる値を生成する", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			id := GenerateGuestDisplayID()
			ids[id] = true
		}
		// 100回生成して、少なくとも90種類以上のユニークな値が生成されることを期待
		assert.GreaterOrEqual(t, len(ids), 90, "ランダム性が不足: 100回中%d種類のみ生成", len(ids))
	})
}

func TestGenerateGuestDisplayName(t *testing.T) {
	t.Run("ゲストで始まる名前を生成する", func(t *testing.T) {
		name := GenerateGuestDisplayName()
		assert.True(t, len(name) > len("ゲスト"), "名前が短すぎる: %s", name)
		assert.Equal(t, "ゲスト", name[:len("ゲスト")])
	})

	t.Run("ゲスト+6文字の英数字で構成される", func(t *testing.T) {
		name := GenerateGuestDisplayName()
		// "ゲスト" (9バイト) + 6文字の英数字
		prefix := "ゲスト"
		suffix := name[len(prefix):]
		assert.Len(t, suffix, 6, "サフィックスの長さが6文字でない: %s", suffix)

		matched, err := regexp.MatchString(`^[a-z0-9]{6}$`, suffix)
		assert.NoError(t, err)
		assert.True(t, matched, "サフィックスが英数字(小文字)でない: %s", suffix)
	})

	t.Run("毎回異なる値を生成する", func(t *testing.T) {
		names := make(map[string]bool)
		for i := 0; i < 100; i++ {
			name := GenerateGuestDisplayName()
			names[name] = true
		}
		assert.GreaterOrEqual(t, len(names), 90, "ランダム性が不足: 100回中%d種類のみ生成", len(names))
	})
}
