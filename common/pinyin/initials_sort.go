package pinyin

import (
	"sort"
)

// SortInitials sorts the string slice by pinyin initial letters (in-place).
func SortInitials(strs []string) {
	a := NewArgs()
	l := len(strs)
	initials := make([]string, l)
	newStrs := map[string]string{}

	for i := 0; i < l; i++ {
		for ii, py := range Pinyin(strs[i], a) {
			if len(py) == 0 {
				initials[i] += string([]rune(strs[i])[ii])
			} else {
				initials[i] += py[0]
			}
		}
		newStrs[initials[i]] = strs[i]
	}

	sort.Strings(initials)

	for i := 0; i < l; i++ {
		strs[i] = newStrs[initials[i]]
	}
}
