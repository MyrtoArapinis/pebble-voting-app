package hashmap

func HashBytes(p []byte) int {
	h := 12345
	for _, b := range p {
		h = (h + int(b)) * 16777619
		h ^= h >> 24
	}
	return h
}
