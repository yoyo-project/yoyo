package file

import "testing"

func TestFindPackagePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "basic test",
			args: args{path: "../../example/"},
			want: "github.com/yoyo-project/yoyo/example/",
		},
		{
			name: "root directory",
			args: args{path: "/"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindPackagePath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPackagePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindPackagePath()\nwant = %v\n got %v", tt.want, got)
			}
		})
	}
}
