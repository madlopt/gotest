package entities

import (
	"fmt"
)

// ConvertIPToUint32 converts an IPv4 string to its 32-bit integer representation
func ConvertIPToUint32(ipStr string) (uint32, bool) {
	var b1, b2, b3, b4 uint32
	n, err := fmt.Sscanf(ipStr, "%d.%d.%d.%d", &b1, &b2, &b3, &b4)
	if err != nil || n != 4 {
		return 0, false
	}
	return (b1 << 24) | (b2 << 16) | (b3 << 8) | b4, true
}
