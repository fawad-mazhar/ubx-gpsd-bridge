package utils

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Contains(slice []byte, val byte) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
