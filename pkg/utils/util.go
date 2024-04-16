package utils

import "strings"

func RemoveDuplicateOfStringArray(content []string) []string {
	var result []string
	temp := map[string]struct{}{}
	for _, item := range content {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func RemoveDuplicateOfString(content string) string {
	str := strings.Split(content, ",")
	m := make(map[string]bool)
	for _, item := range str {
		m[item] = true
	}

	result := ""
	for k := range m {
		if result != "" {
			result += ","
		}
		result += k
	}

	return result
}

func RemoveDuplicatePort(old []int) []int {
	var result []int
	temp := map[int]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicateElement(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
