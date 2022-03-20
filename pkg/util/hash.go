package util

import (
	"hash/fnv"
	"strconv"
)

func NewHash(input string) string {
	h := fnv.New32a()
	h.Write([]byte(input))
	return strconv.FormatUint(uint64(h.Sum32()), 10)
}
