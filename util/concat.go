package util

func Concat(ps ...[]byte) []byte {
	l := 0
	for _, p := range ps {
		l += len(p)
	}
	res := make([]byte, 0, l)
	for _, p := range ps {
		res = append(res, p...)
	}
	return res
}
