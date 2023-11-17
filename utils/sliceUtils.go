package utils

func IsUintInSlice(slice []uint, item uint) bool {
	m := make(map[uint]bool)
	for _, s := range slice {
		m[s] = true
	}
	return m[item]
}
