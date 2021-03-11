package repository

import (
	"reflect"
	"testing"
)

func Test_sortedUnique(t *testing.T) {
	type args struct {
		in []string
	}
	tests := []struct {
		name    string
		args    args
		wantOut []string
	}{
		{
			name:    "already sorted and unique",
			args:    args{in: []string{"a", "b", "c"}},
			wantOut: []string{"a", "b", "c"},
		},
		{
			name:    "three identical values",
			args:    args{in: []string{"a", "a", "a"}},
			wantOut: []string{"a"},
		},
		{
			name:    "unsorted with identical values",
			args:    args{in: []string{"b", "a", "a"}},
			wantOut: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := sortedUnique(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("sortedUnique() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
