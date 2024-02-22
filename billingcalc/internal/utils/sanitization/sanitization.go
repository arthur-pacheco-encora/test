package sanitization

import (
	"regexp"
)

func SanitizeName(name string) string {
	symbolPattern := regexp.MustCompile("[^a-zA-Z0-9]+")
	sanitizedName := symbolPattern.ReplaceAllString(name, "")
	return sanitizedName
}

func SanitizeFloatString(floatStr string) string {
	nonNumericPattern := regexp.MustCompile("[^0-9.Ee+]+")
	sanitizedFloatStr := nonNumericPattern.ReplaceAllString(floatStr, "")
	if sanitizedFloatStr == "" {
		sanitizedFloatStr = "0"
	}
	return sanitizedFloatStr
}

func SanitizeIntegerString(intStr string) string {
	nonNumericPattern := regexp.MustCompile("[^0-9]+")
	sanitizedIntStr := nonNumericPattern.ReplaceAllString(intStr, "")
	if sanitizedIntStr == "" {
		sanitizedIntStr = "0"
	}
	return sanitizedIntStr
}

func SanitizeGoogleSheet(content [][]string, headerRowIndex int) [][]string {
	var sanitizedContent [][]string
	var headerColumnCount int
	for i, row := range content {
		if i < headerRowIndex {
			continue
		} else if i == headerRowIndex {
			headerColumnCount = len(row)
		}
		if len(row) >= headerColumnCount {
			sanitizedContent = append(sanitizedContent, row[:headerColumnCount])
		}
	}
	return sanitizedContent
}
