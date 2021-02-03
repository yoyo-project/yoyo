package schema

import (
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

func TestDatabase_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yml     string
		wantDB  Database
		wantErr bool
	}{
		{
			yml: "dialect: mysql",
			wantDB: Database{
				Dialect: "mysql",
				Tables:  map[string]Table{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := Database{}
			if err := yaml.Unmarshal([]byte(tt.yml), &db); err != nil {
				t.Errorf("Got error %s", err)
			}

			if !reflect.DeepEqual(db, tt.wantDB) {
				t.Errorf("Want %#v, got %#v", tt.wantDB, db)
			}
		})
	}
}

func TestReference_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yml     string
		wantRef Reference
		wantErr bool
	}{
		{
			yml: "has_one: true",
			wantRef: Reference{
				HasOne: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Reference{}
			if err := yaml.Unmarshal([]byte(tt.yml), &r); err != nil {
				t.Errorf("Got error %s", err)
			}

			if !reflect.DeepEqual(r, tt.wantRef) {
				t.Errorf("Want %#v, got %#v", tt.wantRef, r)
			}
		})
	}
}
