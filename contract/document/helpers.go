package document

// I64 возвращает указатель на int64. Удобно для заполнения *int64 полей.
func I64(v int64) *int64 { return &v }

// I64Val возвращает значение из указателя или 0 если nil.
func I64Val(p *int64) int64 {
	if p != nil {
		return *p
	}
	return 0
}
