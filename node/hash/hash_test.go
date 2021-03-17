package hash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsInRange(t *testing.T) {

	testCases := []struct {
		keySlot    int
		localSlot  int
		remoteSlot int
		expected   bool
	}{
		{
			keySlot:    1,
			localSlot:  0,
			remoteSlot: 2,
			expected:   true,
		}, {
			keySlot:    0,
			localSlot:  max_slot_capacity - 1,
			remoteSlot: 1,
			expected:   true,
		}, {
			keySlot:    max_slot_capacity - 1,
			localSlot:  max_slot_capacity - 2,
			remoteSlot: 1,
			expected:   true,
		}, {
			keySlot:    3,
			localSlot:  0,
			remoteSlot: 2,
			expected:   false,
		}, {
			keySlot:    3,
			localSlot:  max_slot_capacity - 1,
			remoteSlot: 2,
			expected:   false,
		}, {
			// test for exclusive of each end
			keySlot:    1,
			localSlot:  0,
			remoteSlot: 1,
			expected:   false,
		}, {
			keySlot:    0,
			localSlot:  0,
			remoteSlot: 2,
			expected:   false,
		}, {
			keySlot:    max_slot_capacity - 1,
			localSlot:  max_slot_capacity - 1,
			remoteSlot: 1,
			expected:   false,
		}, {
			keySlot:    1,
			localSlot:  max_slot_capacity - 1,
			remoteSlot: 1,
			expected:   false,
		},
	}

	for _, tc := range testCases {
		assert.Equalf(t, tc.expected, IsInRange(tc.keySlot, tc.localSlot, tc.remoteSlot), "%+v\n", tc)
	}
}
