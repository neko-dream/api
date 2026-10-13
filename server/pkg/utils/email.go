package utils

import "regexp"

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsEmail matches the address format supported by password login.
func IsEmail(address string) bool { return emailPattern.MatchString(address) }
