package chord

import (
	"github.com/thoas/go-funk"
)

func compareList(oldList []string, newList []string) (extraList, missingList []string) {
	extraList = funk.Filter(newList, func(x string) bool {
		return !isElementInList(x, oldList)
	}).([]string)

	missingList = funk.Filter(oldList, func(x string) bool {
		return !isElementInList(x, newList)
	}).([]string)

	return extraList, missingList
}

func isElementInList(element string, list []string) bool {
	for _, val := range list {
		if val == element {
			return true
		}
	}
	return false
}
