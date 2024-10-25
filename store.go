package gobist

import "sync"

type Store interface {
	Set(key string, value any) error
	Get(key string) (any, error)
	Delete(key string) error
}

type memoryStore struct {
	mx   sync.Mutex
	data map[string]any
}

func newMemoryStore() Store {
	return &memoryStore{data: make(map[string]any)}
}

func (m *memoryStore) Set(key string, value any) error {
	m.mx.Lock()
	m.data[key] = value
	m.mx.Unlock()

	return nil
}

func (m *memoryStore) Get(key string) (any, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	return m.data[key], nil
}

func (m *memoryStore) Delete(key string) error {
	m.mx.Lock()
	delete(m.data, key)
	m.mx.Unlock()

	return nil
}
