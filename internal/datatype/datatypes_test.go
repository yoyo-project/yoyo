package datatype

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func TestDatatype_UnmarshalYAML(t *testing.T) {
	type TestStruct struct {
		Datatype Datatype `yaml:"datatype"`
	}
	tests := []struct {
		name         string
		input        string
		wantDatatype Datatype
		wantErr      bool
	}{
		{
			name:         boolean,
			input:        "datatype: " + boolean,
			wantDatatype: Boolean,
		},
		{
			name:         sbool,
			input:        "datatype: " + sbool,
			wantDatatype: Boolean,
		},
		{
			name:         integer,
			input:        "datatype: " + integer,
			wantDatatype: Integer,
		},
		{
			name:         bigint,
			input:        "datatype: " + bigint,
			wantDatatype: BigInt,
		},
		{
			name:         mediumint,
			input:        "datatype: " + mediumint,
			wantDatatype: MediumInt,
		},
		{
			name:         smallint,
			input:        "datatype: " + smallint,
			wantDatatype: SmallInt,
		},
		{
			name:         tinyint,
			input:        "datatype: " + tinyint,
			wantDatatype: TinyInt,
		},
		{
			name:         decimal,
			input:        "datatype: " + decimal + "(-5, 8)",
			wantDatatype: Decimal,
		},
		{
			name:         varchar,
			input:        "datatype: " + varchar,
			wantDatatype: Varchar,
		},
		{
			name:         "sized" + varchar,
			input:        "datatype: " + varchar + "(5)",
			wantDatatype: Varchar,
		},
		{
			name:         text,
			input:        "datatype: " + text,
			wantDatatype: Text,
		},
		{
			name:         tinytext,
			input:        "datatype: " + tinytext,
			wantDatatype: TinyText,
		},
		{
			name:         mediumtext,
			input:        "datatype: " + mediumtext,
			wantDatatype: MediumText,
		},
		{
			name:         longtext,
			input:        "datatype: " + longtext,
			wantDatatype: LongText,
		},
		{
			name:         blob,
			input:        "datatype: " + blob,
			wantDatatype: Blob,
		},
		{
			name:         char,
			input:        "datatype: " + char,
			wantDatatype: Char,
		},
		{
			name:         enum,
			input:        "datatype: " + enum + "('Hello', 'world!')",
			wantDatatype: Enum,
		},
		{
			name:    "invalid",
			input:   "datatype: " + "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{}

			if err := yaml.Unmarshal([]byte(tt.input), &ts); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantDatatype != ts.Datatype {
				t.Errorf("UnmarshalYAML() Datatype = %s, wantDatatype = %s", ts.Datatype, tt.wantDatatype)
			}
		})
	}
}
