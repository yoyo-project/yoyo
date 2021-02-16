package repositories

import (
	"github.com/yoyo-project/yoyo/example/yoyo/repositories/query/state"
)

type stateRepo struct {
	*repository
}

func (r *stateRepo) FetchOne(builder state.Query) (State, error) {
	panic("implement me")
}

func (r *stateRepo) Search(builder state.Query) (States, error) {
	panic("implement me")
}

func (r *stateRepo) Save(state State) (State, error) {
	panic("implement me")
}

func (r *stateRepo) Delete(builder state.Query) {
	panic("implement me")
}
