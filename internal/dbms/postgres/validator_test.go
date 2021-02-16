package postgres

import (
	"testing"

	"github.com/yoyo-project/yoyo/internal/datatype"
)

func Test_validator_SupportsDatatype(t *testing.T) {
	type args struct {
		dt datatype.Datatype
	}
	tests := []struct {
		args args
		want bool
	}{
		{
			args: args{dt: datatype.Boolean},
			want: true,
		},
		{
			args: args{dt: datatype.Integer},
			want: true,
		},
		{
			args: args{dt: 0},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.dt.String(), func(t *testing.T) {
			va := &validator{}
			if got := va.SupportsDatatype(tt.args.dt); got != tt.want {
				t.Errorf("SupportsDatatype() = %v, want %v", got, tt.want)
			}
		})
	}
}
