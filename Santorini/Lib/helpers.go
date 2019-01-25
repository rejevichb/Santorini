package lib

import (
	"strings"
	"unicode"
)

//Returns true if the given int slice includes the given target int, and false if otherwise
func IntPresent(arr []int, target int) bool {
	for _, candidate := range arr {
		if candidate == target {
			return true
		}
	}
	return false
}

//Lowercase characters
const ALPHA = "abcdefghijklmnopqrstuvwxyz"

//Returns true if the string is composed of entirely alphanumeric characters,
//or false if it has any others
func IsLowercase(s string) bool {
	for _, char := range s {
		if !strings.Contains(ALPHA, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

//Returns true if the given slice includes the given target, and false if otherwise
func StringPresent(arr []string, target string) bool {
	for _, candidate := range arr {
		if candidate == target {
			return true
		}
	}
	return false
}

func StripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}

func StripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
