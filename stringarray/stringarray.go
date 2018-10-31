package stringarray

func Contains(arr []string, target string) bool {
	m := make(map[string]struct{}, len(arr))
	for _, s := range arr {
		m[s] = struct{}{}
	}
	_, ok := m[target]
	return ok
}
