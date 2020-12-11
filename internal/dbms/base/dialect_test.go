package base

import (
	"github.com/dotvezz/yoyo/internal/datatype"
	"testing"
)

func TestBase_TypeString(t *testing.T) {
	type fields struct {
		Dialect string
	}
	type args struct {
		dt datatype.Datatype
	}
	tests := []struct {
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			fields: fields{
				Dialect: "mysql",
			},
			args: args{
				dt: datatype.Integer,
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			fields: fields{
				Dialect: "mysql",
			},
			args: args{
				dt: 0,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.dt.String(), func(t *testing.T) {
			d := &Base{
				Dialect: tt.fields.Dialect,
			}
			got, err := d.TypeString(tt.args.dt)
			if (err != nil) != tt.wantErr {
				t.Errorf("TypeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TypeString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
