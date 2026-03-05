package document

import (
	"strconv"
	"strings"
)

// I64 возвращает указатель на int64. Удобно для заполнения *int64 полей.
func I64(v int64) *int64 { return &v }

// I64Val возвращает значение из указателя или 0 если nil.
func I64Val(p *int64) int64 {
	if p != nil {
		return *p
	}
	return 0
}

// NormalizeDecimalString приводит decimal string к согласованному виду:
// dot separator, без пробелов, без завершающих нулей после точки.
// Функция не гарантирует удаление лидирующих нулей в целой части.
// Examples: "001.500" -> "001.5", "1,250" -> "1.25", "2.000" -> "2".
func NormalizeDecimalString(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, ",", ".")
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	if s == "-0" || s == "+0" {
		return "0"
	}
	return s
}

// FormatDecimalFromFloat64 форматирует float64 в согласованный decimal string.
// Producer'ы должны по возможности избегать промежуточной float arithmetic и предпочитать
// исходные строковые или целочисленные значения из источника.
func FormatDecimalFromFloat64(v float64) string {
	return NormalizeDecimalString(strconv.FormatFloat(v, 'f', -1, 64))
}
