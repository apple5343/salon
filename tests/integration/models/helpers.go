package models

import (
	"regexp"
)

func MatchesPattern(pattern, word string) bool {
	re := regexp.MustCompile("(?i).*" + pattern + ".*")
	return re.MatchString(word)
}
