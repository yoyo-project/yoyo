package template

import (
	"fmt"
	"testing"

	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/schema"
)

func TestGenerateLogic(t *testing.T) {
	type args struct {
		field  string
		col    string
		column schema.Column
	}
	tests := []struct {
		name          string
		args          args
		wantMethods   []string
		wantFunctions []string
		wantImports   []string
	}{
		{
			name: "ID",
			args: args{
				column: schema.Column{
					Datatype: datatype.BigInt,
					Unsigned: true,
					GoName:   "ID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethods, gotFunctions, gotImports := GenerateQueryLogic(tt.args.col, tt.args.column)
			fmt.Println(gotMethods, gotFunctions, gotImports)
		})
	}
}
