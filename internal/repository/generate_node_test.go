package repository

import (
	"strings"
	"testing"

	"github.com/yoyo-project/yoyo/internal/repository/template"
)

func TestNewQueryNodeGenerator(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "basic test",
			want:    template.NodeFile,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := strings.Builder{}
			gotErr := NewQueryNodeGenerator()(&sb)
			got := sb.String()
			if tt.wantErr && gotErr == nil {

			} else if tt.want != got {
				t.Errorf("want:%s\n got:%s", tt.want, got)
			}
		})
	}
}
