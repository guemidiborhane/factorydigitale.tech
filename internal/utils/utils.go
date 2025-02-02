package utils

func GetTimeZone() string {
	return "Africa/Algiers"
}

func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := []T{}

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
