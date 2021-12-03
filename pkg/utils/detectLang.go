package utils

func DetectLanguage(text string) string {
	for _, value := range text {
		if value > 255 {
			return Fa
		}
	}
	return En
}
