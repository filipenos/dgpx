package postal

import "regexp"

var onlyNumber = regexp.MustCompile("[^0-9]+")

//ValidatePostalCode validate if string contains only numbers and
//size equals 8, return a valid postalCode
func ValidatePostalCode(postalCode string) (string, bool) {
	v := onlyNumber.ReplaceAllString(postalCode, "")
	return v, len(v) == 8
}
