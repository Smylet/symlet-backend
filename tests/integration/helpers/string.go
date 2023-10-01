package helpers

import (
	"fmt"
	"strings"
)

// StrReplace makes replacing of multiple placeholders by theirs values in a string.
func StrReplace(str string, original []string, replacement []interface{}) (string, error) {
	if len(original) != len(replacement) {
		return "", fmt.Errorf("length of original and replacement slices do not match")
	}
	replacerArgs := make([]string, 0, len(original)*2)
	for i, replace := range original {
		replacerArgs = append(replacerArgs, fmt.Sprintf("%v", replace), fmt.Sprintf("%v", replacement[i]))
	}
	return strings.NewReplacer(replacerArgs...).Replace(str), nil
}
