package utils

func IsUintInSlice(slice []uint64, item uint64) bool {
	m := make(map[uint64]bool)
	for _, s := range slice {
		m[s] = true
	}
	return m[item]
}
