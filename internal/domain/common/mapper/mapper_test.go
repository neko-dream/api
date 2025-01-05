package mapper_test

import (
	"fmt"
	"testing"

	"github.com/neko-dream/server/internal/domain/common/mapper"
)

// -----------------------------------------
// サンプル用の構造体
// -----------------------------------------
type Source struct {
	ID      int
	Name    string
	Enabled bool
}

type Destination struct {
	ID      int
	Name    string
	Enabled bool
}

// int→string に変換してみるための構造体
type Source2 struct {
	ID   int
	Note string
}

type Destination2 struct {
	ID   string
	Note string
}

// ポインタ型を含んだ構造体
type SourcePtr struct {
	PtrValue *int
}

type DestinationPtr struct {
	PtrValue *int
}

// -----------------------------------------
// テストケース
// -----------------------------------------

// 単純コピーのテスト
func TestMap_SimpleCopy(t *testing.T) {
	src := Source{
		ID:      1,
		Name:    "test",
		Enabled: true,
	}

	dest, err := mapper.Map[Source, Destination](src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dest.ID != src.ID {
		t.Errorf("expected %d, got %d", src.ID, dest.ID)
	}
	if dest.Name != src.Name {
		t.Errorf("expected %s, got %s", src.Name, dest.Name)
	}
	if dest.Enabled != src.Enabled {
		t.Errorf("expected %v, got %v", src.Enabled, dest.Enabled)
	}
}

// カスタムコンバータのテスト (int→stringなど)
func TestMap_WithConverter(t *testing.T) {
	src := Source2{
		ID:   123,
		Note: "hello",
	}

	dest, err := mapper.MapWithConverters[Source2, Destination2](
		src,
		mapper.Converter(func(i int) (string, error) {
			return fmt.Sprintf("Value_%d", i), nil
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// IDはカスタムコンバータにより "Value_123" に変換されるはず
	expectedID := "Value_123"
	if dest.ID != expectedID {
		t.Errorf("expected %s, got %s", expectedID, dest.ID)
	}

	// Note は単純コピー
	if dest.Note != src.Note {
		t.Errorf("expected %s, got %s", src.Note, dest.Note)
	}
}

// ポインタ型を含むコピーのテスト
func TestMap_PointerField(t *testing.T) {
	// ポインタ型のゼロ値やnilコピーが意図通りか検証する
	val := 10
	src := SourcePtr{
		PtrValue: &val,
	}

	dest, err := mapper.Map[SourcePtr, DestinationPtr](src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dest.PtrValue == nil {
		t.Error("expected ptrValue to not be nil")
		return
	}
	if *dest.PtrValue != val {
		t.Errorf("expected %d, got %d", val, *dest.PtrValue)
	}
}
