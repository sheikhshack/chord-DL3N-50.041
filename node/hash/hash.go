package hash

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
)

// max_slot_capacity represents the maximum num of slots for consistent hashing
const max_slot_capacity = 100000

// Hash converts a string into an int corresponding to a consistent hashing slot
func Hash(key string) int {
	data := []byte(key)
	hashed := sha1.Sum(data)
	return slotHash(hashed[:])
}

// slotHash converts []byte into an int corresponding to a consistent hashing slot
func slotHash(hashed []byte) (slot int) {
	data := binary.BigEndian.Uint64(hashed)
	return int(data % max_slot_capacity)
}

// IsInRange checks if an slot can be located between current slot and another slot
func IsInRange(keySlot, localSlot, remoteSlot int) bool {
	if localSlot == remoteSlot {
		return true
	} else if keySlot == Hash("") {
		return false
	}
	// check if range covers start of circle
	if localSlot > remoteSlot {
		// case of keySlot before start of circle: always in range
		if keySlot > localSlot {
			return true
		}
		// case of keySlot located after start of circle
		return (remoteSlot > keySlot)
	}

	return (localSlot < keySlot && keySlot < remoteSlot)
}

func main() {
	fmt.Println(Hash("remember.rar"))
	fmt.Println(Hash("forgotten.rar"))
	fmt.Println("alpha:", Hash("alpha"))
	fmt.Println("bravo:", Hash("nodeBravo"))
	fmt.Println("charlie:", Hash("nodeCharlie"))
	fmt.Println("delta:", Hash("nodeDelta"))
}
