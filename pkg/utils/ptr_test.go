package utils_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
)

func TestMain(m *testing.M) {
	m.Run()
}

type CustomString string
type CustomInt int
type CustomFloat float64

func TestToNullableSQL_String(t *testing.T) {
	t.Run("有効な文字列", func(t *testing.T) {
		result := utils.ToNullableSQL[sql.NullString]("hello")
		expected := sql.NullString{String: "hello", Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("空文字列", func(t *testing.T) {
		result := utils.ToNullableSQL[sql.NullString]("")
		expected := sql.NullString{String: "", Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("文字列ポインタ", func(t *testing.T) {
		str := "hello"
		result := utils.ToNullableSQL[sql.NullString](&str)
		expected := sql.NullString{String: "hello", Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("nil文字列ポインタ", func(t *testing.T) {
		var str *string
		result := utils.ToNullableSQL[sql.NullString](str)
		expected := sql.NullString{String: "", Valid: false}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("カスタム文字列型", func(t *testing.T) {
		custom := CustomString("custom")
		result := utils.ToNullableSQL[sql.NullString](custom)
		expected := sql.NullString{String: "custom", Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}

func TestToNullableSQL_Int(t *testing.T) {
	t.Run("有効なint32", func(t *testing.T) {
		result := utils.ToNullableSQL[sql.NullInt32](42)
		expected := sql.NullInt32{Int32: 42, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("有効なint8", func(t *testing.T) {
		var val int8 = 8
		result := utils.ToNullableSQL[sql.NullInt16](val)
		expected := sql.NullInt16{Int16: 8, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("有効なint16", func(t *testing.T) {
		var val int16 = 16
		result := utils.ToNullableSQL[sql.NullInt16](val)
		expected := sql.NullInt16{Int16: 16, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("有効なint64", func(t *testing.T) {
		var val int64 = 64
		result := utils.ToNullableSQL[sql.NullInt64](val)
		expected := sql.NullInt64{Int64: 64, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("nil整数ポインタ", func(t *testing.T) {
		var val *int
		result := utils.ToNullableSQL[sql.NullInt32](val)
		expected := sql.NullInt32{Int32: 0, Valid: false}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("カスタム整数型", func(t *testing.T) {
		custom := CustomInt(123)
		result := utils.ToNullableSQL[sql.NullInt32](custom)
		expected := sql.NullInt32{Int32: 123, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}

func TestToNullableSQL_Bool(t *testing.T) {
	t.Run("有効なbool_true", func(t *testing.T) {
		result := utils.ToNullableSQL[sql.NullBool](true)
		expected := sql.NullBool{Bool: true, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("有効なbool_false", func(t *testing.T) {
		result := utils.ToNullableSQL[sql.NullBool](false)
		expected := sql.NullBool{Bool: false, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("nilブールポインタ", func(t *testing.T) {
		var val *bool
		result := utils.ToNullableSQL[sql.NullBool](val)
		expected := sql.NullBool{Bool: false, Valid: false}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}

func TestToNullableSQL_Float(t *testing.T) {
	t.Run("有効なfloat32", func(t *testing.T) {
		var val float32 = 3.14
		result := utils.ToNullableSQL[sql.NullFloat64](val)
		expected := sql.NullFloat64{Float64: float64(val), Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("有効なfloat64", func(t *testing.T) {
		val := 3.14159
		result := utils.ToNullableSQL[sql.NullFloat64](val)
		expected := sql.NullFloat64{Float64: 3.14159, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("nil浮動小数点ポインタ", func(t *testing.T) {
		var val *float64
		result := utils.ToNullableSQL[sql.NullFloat64](val)
		expected := sql.NullFloat64{Float64: 0, Valid: false}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("カスタム浮動小数点型", func(t *testing.T) {
		custom := CustomFloat(2.71)
		result := utils.ToNullableSQL[sql.NullFloat64](custom)
		expected := sql.NullFloat64{Float64: 2.71, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}

func TestToNullableSQL_Time(t *testing.T) {
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)

	t.Run("有効な時刻", func(t *testing.T) {
		result := utils.ToNullableSQL[sql.NullTime](testTime)
		expected := sql.NullTime{Time: testTime, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("時刻ポインタ", func(t *testing.T) {
		result := utils.ToNullableSQL[sql.NullTime](&testTime)
		expected := sql.NullTime{Time: testTime, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("nil時刻ポインタ", func(t *testing.T) {
		var val *time.Time
		result := utils.ToNullableSQL[sql.NullTime](val)
		expected := sql.NullTime{Time: time.Time{}, Valid: false}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}

func TestToNullableSQL_UUID(t *testing.T) {
	testUUID := uuid.New()

	t.Run("有効なUUID", func(t *testing.T) {
		result := utils.ToNullableSQL[uuid.NullUUID](testUUID)
		expected := uuid.NullUUID{UUID: testUUID, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("UUIDポインタ", func(t *testing.T) {
		result := utils.ToNullableSQL[uuid.NullUUID](&testUUID)
		expected := uuid.NullUUID{UUID: testUUID, Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("nilUUIDポインタ", func(t *testing.T) {
		var val *uuid.UUID
		result := utils.ToNullableSQL[uuid.NullUUID](val)
		expected := uuid.NullUUID{UUID: uuid.UUID{}, Valid: false}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}

func TestToNullableSQL_EdgeCases(t *testing.T) {
	t.Run("nilインターフェース", func(t *testing.T) {
		var val interface{}
		result := utils.ToNullableSQL[sql.NullString](val)
		expected := sql.NullString{String: "", Valid: false}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("未サポート型", func(t *testing.T) {
		type UnsupportedType struct{ Value string }
		val := UnsupportedType{Value: "test"}
		result := utils.ToNullableSQL[sql.NullString](val)
		expected := sql.NullString{} // ゼロ値
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("二重ポインタ", func(t *testing.T) {
		str := "hello"
		strPtr := &str
		result := utils.ToNullableSQL[sql.NullString](&strPtr)
		expected := sql.NullString{String: "hello", Valid: true}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("nil二重ポインタ", func(t *testing.T) {
		var str **string
		result := utils.ToNullableSQL[sql.NullString](str)
		expected := sql.NullString{String: "", Valid: false}
		if result != expected {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}

// ベンチマークテスト
func BenchmarkToNullableSQL_String(b *testing.B) {
	str := "benchmark test string"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullString](str)
	}
}

func BenchmarkToNullableSQL_StringPointer(b *testing.B) {
	str := "benchmark test string"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullString](&str)
	}
}

func BenchmarkToNullableSQL_NilStringPointer(b *testing.B) {
	var str *string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullString](str)
	}
}

func BenchmarkToNullableSQL_Int32(b *testing.B) {
	val := int32(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullInt32](val)
	}
}

func BenchmarkToNullableSQL_Int64(b *testing.B) {
	val := int64(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullInt64](val)
	}
}

func BenchmarkToNullableSQL_Bool(b *testing.B) {
	val := true
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullBool](val)
	}
}

func BenchmarkToNullableSQL_Float64(b *testing.B) {
	val := 3.14159
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullFloat64](val)
	}
}

func BenchmarkToNullableSQL_Time(b *testing.B) {
	val := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullTime](val)
	}
}

func BenchmarkToNullableSQL_UUID(b *testing.B) {
	val := uuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[uuid.NullUUID](val)
	}
}

func BenchmarkToNullableSQL_CustomType(b *testing.B) {
	val := CustomString("custom test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullString](val)
	}
}

func BenchmarkToNullableSQL_UnsupportedType(b *testing.B) {
	type UnsupportedType struct{ Value string }
	val := UnsupportedType{Value: "test"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullString](val)
	}
}

func BenchmarkToNullableSQL_NestedPointer(b *testing.B) {
	str := "--------------------------------------------------------------------"
	strPtr := lo.ToPtr(lo.ToPtr(lo.ToPtr(lo.ToPtr(lo.ToPtr(lo.ToPtr(lo.ToPtr(str)))))))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ToNullableSQL[sql.NullString](strPtr)
	}
}

// パフォーマンス比較用
func BenchmarkToNullableSQL_AllTypes(b *testing.B) {
	testData := []interface{}{
		"string",
		int32(42),
		int64(64),
		true,
		3.14159,
		time.Now(),
		uuid.New(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, data := range testData {
			switch data.(type) {
			case string:
				_ = utils.ToNullableSQL[sql.NullString](data)
			case int32:
				_ = utils.ToNullableSQL[sql.NullInt32](data)
			case int64:
				_ = utils.ToNullableSQL[sql.NullInt64](data)
			case bool:
				_ = utils.ToNullableSQL[sql.NullBool](data)
			case float64:
				_ = utils.ToNullableSQL[sql.NullFloat64](data)
			case time.Time:
				_ = utils.ToNullableSQL[sql.NullTime](data)
			case uuid.UUID:
				_ = utils.ToNullableSQL[uuid.NullUUID](data)
			}
		}
	}
}
