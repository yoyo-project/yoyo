package person

import (
	"fmt"
	"testing"
)

func TestID(t *testing.T) {
	sql, args := ID(1).IDLessThan(2).Or(ID(3).IDGreaterThan(4)).SQL()
	fmt.Printf("%v %v", sql, args)
}
