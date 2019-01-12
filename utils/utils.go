package utils

func Join(ss []string) string {
	var out string
	for _, s := range ss {
		out += s
	}

	return out
}
