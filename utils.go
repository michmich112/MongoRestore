package main

// StringsEqual is used to compare two arrays by value
// Perform an in-order coparison
func StringsEqual(base []string, comp []string) bool {
	if len(base) != len(comp) {
		return false
	}
	for index, element := range base {
		if element != comp[index] {
			return false
		}
	}
	return true
}

// StringsIndexOf is used to retrieve the index of an element from a slice
func StringsIndexOf(slice []string, element string) int {
	for index, str := range slice {
		if str == element {
			return index
		}
	}
	return -1
}

// StringsRemoveElements is used to remove elements from an string slice
func StringsRemoveElements(slice []string, elements ...string) []string {
	var indexes []int
	cSlice := make([]string, len(slice))
	copy(cSlice, slice)
	for index, str := range slice {
		if StringsIndexOf(elements, str) > -1 {
			indexes = append(indexes, index)
		}
	}
	for pad, index := range indexes {
		cSlice = append(cSlice[:index-pad], cSlice[index-pad+1:]...)
	}
	return cSlice
}

//func ExecutionManager(fn func(int)) {}
