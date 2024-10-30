package postgres

import (
	"fmt"
	"github.com/yoyo-project/yoyo/example/postgres/yoyo/repositories/query/person"
	"testing"
)

func TestThing(t *testing.T) {
	q := person.Age(10).Name("ben").FavoriteColor("blue").Or(person.Age(24))
	fmt.Println(q.SQL())
}
