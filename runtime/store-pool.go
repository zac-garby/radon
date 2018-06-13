package runtime

const initialStorePoolSize = 32

// A StorePool contains a number of Stores, which can be released quickly. This avoids
// creating stores, which is slower than just reusing existing ones.
type StorePool struct {
	stores []*Store
}

// NewStorePool makes a new store pool with initialStorePoolSize stores.
func NewStorePool() *StorePool {
	pool := &StorePool{
		stores: make([]*Store, initialStorePoolSize),
	}

	for i := 0; i < initialStorePoolSize; i++ {
		pool.stores[i] = NewStore(nil)
	}

	return pool
}

// IsEmpty checks whether or not the pool is empty.
func (s *StorePool) IsEmpty() bool {
	return len(s.stores) == 0
}

// Release releases a store from the pool, allowing for use elsewhere.
func (s *StorePool) Release(enclosing *Store) *Store {
	if s.IsEmpty() {
		return NewStore(enclosing)
	}

	store := s.stores[0]
	s.stores = s.stores[1:]
	store.Enclosing = enclosing

	return store
}

// Add adds a store back into the pool.
func (s *StorePool) Add(sto *Store) {
	sto.Data = make(map[string]*Variable)
	sto.Enclosing = nil
	s.stores = append(s.stores, sto)
}
