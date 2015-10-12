package prettybenchmarks

import "testing"

var (
	testsFloat = []struct {
		input    float64
		format   string
		expected string
	}{
		{12345.6789, "#,###.##", "12,345.68"},
		{12345.6789, "#,###.", "12,346"},
		{12345.6789, "#,###", "12345,679"},
		{12345.6789, "#\u202F###,##", "12 345,68"},
		{12345.6789, "#.###,######", "12.345,678900"},
	}
	testsInt = []struct {
		input    int
		format   string
		expected string
	}{
		{123456789, "#,###.##", "123,456,789.00"},
		{123456789, "#,###.", "123,456,789"},
		{123456789, "#,###", "123456789,000"},
		{123456789, "#\u202F###,##", "123 456 789,00"},
		{123456789, "#.###,######", "123.456.789,000000"},
	}
)

func Test_RenderFloat(t *testing.T) {
	for _, tt := range testsFloat {
		actual := RenderFloat(tt.format, tt.input)
		if tt.expected != actual {
			t.Errorf("Parsing float %f with format %s: expected: %v, got: %v", tt.input, tt.format, tt.expected, actual)
		}
	}
}

func Test_RenderInteger(t *testing.T) {
	for _, tt := range testsInt {
		actual := RenderInteger(tt.format, tt.input)
		if tt.expected != actual {
			t.Errorf("Parsing integer %d with format %s: expected: %v, got: %v", tt.input, tt.format, tt.expected, actual)
		}
	}
}
