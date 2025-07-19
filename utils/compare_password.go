package utils

func ComparePassword(password1 string, password2 string) bool {
	if password1 == password2 {
		return true
	} else {
		return false
	}
}
