package utils

func ContainInt64(arr []int64, value int64) bool {
	for _, a := range arr {
		if a == value {
			return true
		}
	}
	return false
}
func ContainString(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func FindInt(arr []int, search int) int {
	for index, value := range arr {
		if value == search {
			return index
		}
	}
	return -1
}

func Unique(arr []float64) []float64 {
	occured := map[float64]bool{}
	result := []float64{}
	for e := range arr {
		if occured[arr[e]] != true {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}
