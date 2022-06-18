package tools

func EqualSlice(s, v []byte) bool {
	if len(s) != len(v) {
		return false
	}

	if (s == nil) != (v == nil) {
		return false
	}

	for i, val := range s {
		if val != v[i] {
			return false
		}
	}

	return true
}
