package datatype

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
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
			input:        "datatype: " + decimal + "(6, 3)",
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

func TestDatatype_String(t *testing.T) {
	tests := []struct {
		dt   Datatype
		want string
	}{
		{
			dt:   Integer,
			want: integer,
		},
		{
			dt:   TinyInt,
			want: tinyint,
		},
		{
			dt:   SmallInt,
			want: smallint,
		},
		{
			dt:   MediumInt,
			want: mediumint,
		},
		{
			dt:   BigInt,
			want: bigint,
		},
		{
			dt:   Decimal,
			want: decimal,
		},
		{
			dt:   Varchar,
			want: varchar,
		},
		{
			dt:   Text,
			want: text,
		},
		{
			dt:   TinyText,
			want: tinytext,
		},
		{
			dt:   MediumText,
			want: mediumtext,
		},
		{
			dt:   LongText,
			want: longtext,
		},
		{
			dt:   Char,
			want: char,
		},
		{
			dt:   Blob,
			want: blob,
		},
		{
			dt:   Enum,
			want: enum,
		},
		{
			dt:   Boolean,
			want: boolean,
		},
		{
			dt:   123123,
			want: "NONE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.dt.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatatype_IsInt(t *testing.T) {
	tests := []struct {
		dt   Datatype
		want bool
	}{
		{
			dt:   Integer,
			want: true,
		},
		{
			dt:   Text,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.dt.String(), func(t *testing.T) {
			if got := tt.dt.IsInt(); got != tt.want {
				t.Errorf("IsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatatype_IsNumeric(t *testing.T) {
	tests := []struct {
		dt   Datatype
		want bool
	}{
		{
			dt:   Integer,
			want: true,
		},
		{
			dt:   Decimal,
			want: true,
		},
		{
			dt:   Text,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.dt.String(), func(t *testing.T) {
			if got := tt.dt.IsNumeric(); got != tt.want {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatatype_IsBinary(t *testing.T) {
	tests := []struct {
		dt   Datatype
		want bool
	}{
		{
			dt:   Integer,
			want: false,
		},
		{
			dt:   Decimal,
			want: false,
		},
		{
			dt:   Text,
			want: false,
		},
		{
			dt:   Blob,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.dt.String(), func(t *testing.T) {
			if got := tt.dt.IsBinary(); got != tt.want {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatatype_RequiresScale(t *testing.T) {
	tests := []struct {
		dt   Datatype
		want bool
	}{
		{
			dt:   Integer,
			want: false,
		},
		{
			dt:   Decimal,
			want: true,
		},
		{
			dt:   Text,
			want: false,
		},
		{
			dt:   Enum,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.dt.String(), func(t *testing.T) {
			if got := tt.dt.RequiresParams(); got != tt.want {
				t.Errorf("RequiresParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatatype_IsString(t *testing.T) {
	tests := []struct {
		dt   Datatype
		want bool
	}{
		{
			dt:   Integer,
			want: false,
		},
		{
			dt:   Decimal,
			want: false,
		},
		{
			dt:   Text,
			want: true,
		},
		{
			dt:   Enum,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.dt.String(), func(t *testing.T) {
			if got := tt.dt.IsString(); got != tt.want {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatatype_IsSignable(t *testing.T) {
	tests := []struct {
		dt   Datatype
		want bool
	}{
		{
			dt:   Integer,
			want: true,
		},
		{
			dt:   Decimal,
			want: true,
		},
		{
			dt:   Text,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.dt.String(), func(t *testing.T) {
			if got := tt.dt.IsSignable(); got != tt.want {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatatype_MarshalYAML(t *testing.T) {
	tests := []struct {
		dt   Datatype
		want string
	}{
		{
			dt:   Integer,
			want: "integer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.dt.String(), func(t *testing.T) {
			got, _ := tt.dt.MarshalYAML()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalYAML() got = %v, want %v", got, tt.want)
			}
		})
	}
}
