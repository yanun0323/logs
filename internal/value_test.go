package internal

import "testing"

func BenchmarkValueToString(b *testing.B) {
	values := []any{
		"string_value",
		42,
		int64(42),
		uint64(42),
		3.14159,
		true,
		false,
		nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = ValueToString(v)
		}
	}
}
