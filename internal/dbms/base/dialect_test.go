package base

import (
	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/schema"
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

// We expect the "base" dialect to panic on every method except TypeString
func TestBase_Panics(t *testing.T) {
	didPanic := func(f func()) (res bool) {
		defer func() {
			if r := recover(); r != nil {
				res = true
			}
		}()

		f()

		return
	}

	b := Base{}

	t.Run("AddReference", func(t *testing.T) {
		panicked := didPanic(func() {
			_, _ = b.AddReference("", "", schema.Database{}, schema.Reference{})
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})

	t.Run("CreateTable", func(t *testing.T) {
		panicked := didPanic(func() {
			_ = b.CreateTable("", schema.Table{})
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})

	t.Run("AddIndex", func(t *testing.T) {
		panicked := didPanic(func() {
			_ = b.AddIndex("", "", schema.Index{})
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})

	t.Run("AddColumn", func(t *testing.T) {
		panicked := didPanic(func() {
			_ = b.AddColumn("", "", schema.Column{})
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})

	t.Run("ListTables", func(t *testing.T) {
		panicked := didPanic(func() {
			_, _ = b.ListTables()
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})

	t.Run("ListColumns", func(t *testing.T) {
		panicked := didPanic(func() {
			_, _ = b.ListColumns("")
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})

	t.Run("ListIndices", func(t *testing.T) {
		panicked := didPanic(func() {
			_, _ = b.ListIndices("")
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})

	t.Run("ListReferences", func(t *testing.T) {
		panicked := didPanic(func() {
			_, _ = b.ListReferences("")
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})

	t.Run("GetColumn", func(t *testing.T) {
		panicked := didPanic(func() {
			_, _ = b.GetColumn("", "")
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})

	t.Run("GetIndex", func(t *testing.T) {
		panicked := didPanic(func() {
			_, _ = b.GetIndex("", "")
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})

	t.Run("GetReference", func(t *testing.T) {
		panicked := didPanic(func() {
			_, _ = b.GetReference("", "")
		})
		if !panicked {
			t.Errorf("Expected a panic but didn't see one")
		}
	})
}
