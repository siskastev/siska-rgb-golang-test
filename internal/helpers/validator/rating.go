package validator

import "regexp"

func IsValidRating(rating string) bool {
	pattern := `^[0-9](\.([0-9]))?$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(rating)
}
