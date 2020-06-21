package string

type String string

// Exists checks if a string is present in a string slice
func (s String) Exists(elements []string) bool {
	for _, ele := range elements {
		if ele == string(s) {
			return true
		}
	}
	return false
}
