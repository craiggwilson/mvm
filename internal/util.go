package internal

func stringsContains(values []string, find string) bool {
	return stringsIndex(values, find) >= 0
}

func stringsIndex(values []string, find string) int {
	for i, s := range values {
		if s == find {
			return i
		}
	}

	return -1
}
