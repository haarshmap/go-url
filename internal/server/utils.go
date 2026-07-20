package server

import "hash/fnv"

func FastHash(data string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(data))
	return h.Sum32()
}
