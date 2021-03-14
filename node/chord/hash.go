package chord

import (
	"crypto/sha1"
	"encoding/binary"
)

// max_slot_capacity represents the maximum num of slots for consistent hashing
const max_slot_capacity = 360

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
