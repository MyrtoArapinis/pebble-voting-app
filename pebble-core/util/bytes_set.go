package util

import "bytes"

type bucket struct {
	value []byte
	ok    bool
}

type BytesSet struct {
	v   []bucket
	len int
}

func hashBytes(p []byte) int {
	h := 12345
	for _, b := range p {
		h = (h + int(b)) * 16777619
		h ^= h >> 24
	}
	return h
}

func (s *BytesSet) Len() int {
	return s.len
}

func (s *BytesSet) Contains(p []byte) bool {
	if s.len == 0 {
		return false
	}
	hash := hashBytes(p)
	mask := len(s.v) - 1
	idx := hash & mask
	for {
		if !s.v[idx].ok || hashBytes(s.v[idx].value) != hash {
			return false
		}
		if bytes.Equal(p, s.v[idx].value) {
			return true
		}
		idx = (idx + 1) & mask
	}
}

func (s *BytesSet) Put(p []byte) {
	if s.len*2 >= len(s.v) {
		s.resize()
	}
	mask := len(s.v) - 1
	idx := hashBytes(p) & mask
	for {
		if !s.v[idx].ok {
			s.v[idx].value = p
			s.v[idx].ok = true
			s.len++
			return
		}
		idx = (idx + 1) & mask
	}
}

func (s *BytesSet) Clear() {
	s.v = nil
	s.len = 0
}

func (s *BytesSet) resize() {
	old := s.v
	s.len = 0
	if len(old) == 0 {
		s.v = make([]bucket, 16)
	} else {
		s.v = make([]bucket, len(old)*2)
		for _, b := range old {
			if b.ok {
				s.Put(b.value)
			}
		}
	}
}
