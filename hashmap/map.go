package hashmap

type Key interface {
	Equals(other Key) bool
	Hash() int
}

type KeyValue struct {
	Key   Key
	Value interface{}
}

type Map struct {
	v   []KeyValue
	len int
}

func (m *Map) Len() int {
	return m.len
}

func (m *Map) Get(key Key) (interface{}, bool) {
	if m.len == 0 {
		return nil, false
	}
	hash := key.Hash()
	mask := len(m.v) - 1
	idx := hash & mask
	for {
		if m.v[idx].Key == nil || m.v[idx].Key.Hash() != hash {
			return nil, false
		}
		if key.Equals(m.v[idx].Key) {
			return m.v[idx].Value, true
		}
		idx = (idx + 1) & mask
	}
}

func (m *Map) Clear() {
	m.v = nil
	m.len = 0
}

func (m *Map) Put(key Key, value interface{}) {
	if m.len*2 >= len(m.v) {
		m.resize()
	}
	mask := len(m.v) - 1
	idx := key.Hash() & mask
	for {
		if m.v[idx].Key == nil {
			m.v[idx].Key = key
			m.v[idx].Value = value
			m.len++
			return
		}
		idx = (idx + 1) & mask
	}
}

func (m *Map) resize() {
	old := m.v
	if len(old) == 0 {
		m.v = make([]KeyValue, 16)
	} else {
		m.v = make([]KeyValue, len(old)*2)
	}
	m.len = 0
	for _, kv := range old {
		if kv.Key != nil {
			m.Put(kv.Key, kv.Value)
		}
	}
}
