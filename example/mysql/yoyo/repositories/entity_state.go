package repositories

import (
	"database/sql"
	"fmt"
)

type State struct {
	ID   uint
	Name string

	persisted *State
}

type States struct {
	rs sql.Rows
}

func (s *States) Next() bool {
	return s.rs.Next()
}

func (s *States) Scan(state *State) (err error) {
	if state == nil {
		return fmt.Errorf("in States.Scan: passed a nil entity")
	}

	err = s.rs.Scan(state.ID, state.Name)
	persisted := *state
	state.persisted = &persisted
	return err
}
