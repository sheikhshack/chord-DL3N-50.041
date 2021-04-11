package chord

import (
	"log"

	"github.com/sheikhshack/distributed-chaos-50.041/node/hash"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"

	"github.com/thoas/go-funk"
)

func compareList(oldList []string, newList []string) (extraList, missingList []string) {
	extraList = funk.UniqString(funk.Filter(newList, func(x string) bool {
		return !isElementInList(x, oldList)
	}).([]string))

	missingList = funk.UniqString(funk.Filter(oldList, func(x string) bool {
		return !isElementInList(x, newList)
	}).([]string))

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

func getAllLocalFiles() (keys, values string) {

	// Get all the replica files in the store
	files, err := store.GetAll("local")
	if err != nil {
		print(err)
		return
	}
	keys = ""
	values = ""

	for _, i := range files {
		log.Printf("Has Filename:%v, HashedFile: %v", i.Name(), hash.Hash(i.Name()))
		keys += i.Name() + ","
		val, _ := store.Get("local", i.Name())
		values += string(val) + ","
	}

	return keys, values

}
