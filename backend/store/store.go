package store
import (
	"math/rand"
	"sync"

	"github.com/rjmarques/something-of-the-day/model"
)

type Store struct {
	somethings []*model.Something
	index      map[int64]*model.Something
	mu         sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		index: map[int64]*model.Something{},
	}
}

func (st *Store) Add(somethings []*model.Something) {
	st.mu.Lock()
	defer st.mu.Unlock()

	for _, s := range somethings {
		if _, found := st.index[s.Id]; found {
			continue // already exists in the store
		}

		st.index[s.Id] = s
		st.somethings = append(st.somethings, s)
	}
}

func (st *Store) GetRand() *model.Something {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if len(st.somethings) == 0 {
		return nil
	}

	index := rand.Int() % len(st.somethings)
	return st.somethings[index]
}

func (st *Store) GetLatest() *model.Something {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if len(st.somethings) == 0 {
		return nil
	}

	return st.somethings[len(st.somethings)-1] // need to be careful of the order the API returns tweets. Pos 0 has the latest but shouldn't...
}

