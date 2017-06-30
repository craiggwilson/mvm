package internal

func stringsIndex(values []string, find string) int {
	for i, s := range values {
		if s == find {
			return i
		}
	}

	return -1
}
