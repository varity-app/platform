package common

// StringToPtr converts a string into a pointer of the same string
func StringToPtr(s string) *string {
	return &s
}
