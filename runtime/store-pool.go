package runtime

const initialStorePoolSize = 32

// A StorePool contains a number of Stores, which can be released quickly. This avoids
// creating stores, which is slower than just reusing existing ones.
type StorePool struct {
	stores []*Store
}

func NewStorePool() *StorePool {
	return &StorePool{
		stores: make([]*Store, initialStorePoolSize),
	}
}

func (s *StorePool) IsEmpty() bool {
	return len(s.stores) == 0
}

func (s *StorePool) Release(enclosing *Store) *Store {
	if s.IsEmpty() {
		return NewStore(enclosing)
	} else {
		store := s.stores[0]
		s.stores = s.stores[1:]
		return store
	}
}

func (s *StorePool) Add(sto *Store) {
	sto.Data = make(map[string]*Variable)
	sto.Enclosing = nil
	s.stores = append(s.stores, sto)
}
