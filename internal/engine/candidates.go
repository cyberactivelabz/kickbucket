package engine

import "strings"

// Candidates returns unique bucket-name candidates in deterministic order.
func Candidates(keyword string, words []string) []string {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return nil
	}

	seen := make(map[string]struct{})
	result := make([]string, 0, len(words)*2+1)
	add := func(value string) {
		value = strings.Trim(strings.ToLower(strings.TrimSpace(value)), "-")
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}

	add(keyword)
	for _, word := range words {
		word = strings.Trim(strings.ToLower(strings.TrimSpace(word)), "-")
		if word == "" {
			continue
		}
		add(keyword + "-" + word)
		add(word + "-" + keyword)
	}
	return result
}
