package string

func replaceSpaces(S string, length int) string {
	var ans []rune
	rs := []rune(S)
	i := 0
	isReplaced := false
	for rs[i] == ' ' && i < len(rs) {
		i++
	}
	for i < len(rs) {
		if rs[i] == ' ' {
			if !isReplaced {
				ans = append(ans, []rune("%20")...)
				isReplaced = true
			}
		} else {
			ans = append(ans, rs[i])
			isReplaced = false
		}
		i++
	}
	return string(ans)
}
