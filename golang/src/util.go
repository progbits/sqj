package main

func escapeString(str string) string {
	if str == "" {
		return ""
	}

	result := "'"
	for _, c := range str {
		if c == '\'' {
			result += string('\'')
			result += string(c)
			continue
		}
		result += string(c)
	}
	return result + "'"
}

func unescapeString(str string) string {
	if str == "" {
		return ""
	}

	result := ""
	skip := false
	for _, c := range str[1 : len(str)-1] {
		if skip {
			skip = false
			continue
		}

		if c == '\'' {
			result += string(c)
			skip = true
			continue
		}

		result += string(c)
		skip = false
	}
	return result
}
