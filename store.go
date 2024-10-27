package gobist

import "sync"

type Store interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type memoryStore struct {
	mx   sync.Mutex
	data map[string]string
}

func newMemoryStore() Store {
	return &memoryStore{data: make(map[string]string)}
}

func (m *memoryStore) Set(key string, value string) error {
	m.mx.Lock()
	m.data[key] = value
	m.mx.Unlock()

	return nil
}

func (m *memoryStore) Get(key string) (string, error) {
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
