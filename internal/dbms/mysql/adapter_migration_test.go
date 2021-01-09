package mysql

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dotvezz/yoyo/internal/datatype"
	"github.com/dotvezz/yoyo/internal/dbms/base"
	"github.com/dotvezz/yoyo/internal/dbms/dialect"
	"github.com/dotvezz/yoyo/internal/schema"
)

func TestNewadapter(t *testing.T) {
	tests := []struct {
		name string
		want *adapter
	}{
		{
			name: "just an adapter",
			want: &adapter{
				Base: base.Base{
					Dialect: dialect.MySQL,
				},
				validator: validator{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAdapter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAdapter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_adapter_TypeString(t *testing.T) {
	tests := map[string]struct {
		dt      datatype.Datatype
		wantS   string
		wantErr string
	}{
		datatype.Integer.String(): {
			dt:    datatype.Integer,
			wantS: "INT",
		},
		datatype.BigInt.String(): {
			dt:    datatype.BigInt,
			wantS: "BIGINT",
		},
		datatype.SmallInt.String(): {
			dt:    datatype.SmallInt,
			wantS: "SMALLINT",
		},
		"unsupported datatype": {
			dt:      datatype.Boolean,
			wantErr: "unsupported datatype",
		},
		"invalid datatype": {
			dt:      0,
			wantErr: "invalid datatype",
		},
	}

	m := &adapter{
		Base:      base.Base{Dialect: dialect.MySQL},
		validator: validator{},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotS, err := m.TypeString(tt.dt)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected error `nil`, got error `%v`", err)
				} else if gotS != tt.wantS {
					t.Errorf("expected string `%s`, got string `%s`", tt.wantS, gotS)
				}
			}

			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("expected error `%v`, got error `nil`", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected error `%v`, got error `%v`", tt.wantErr, err)
				}
			}
		})
	}
}

func Test_adapter_CreateTable(t *testing.T) {
	tests := map[string]struct {
		tName string
		t     schema.Table
		wantS string
	}{
		"empty table": {
			tName: "table",
			wantS: "CREATE TABLE `table` (\n\n);",
		},
		"single column no primary key": {
			tName: "table",
			t: schema.Table{
				Columns: map[string]schema.Column{
					"column": {
						Datatype: datatype.Integer,
					},
				},
			},
			wantS: "CREATE TABLE `table` (\n    `column` INT SIGNED NOT NULL\n);",
		},
		"two column with primary key": {
			tName: "table",
			t: schema.Table{
				Columns: map[string]schema.Column{
					"column": {
						Datatype:   datatype.Integer,
						PrimaryKey: true,
					},
					"column2": {
						Datatype: datatype.Integer,
					},
				},
			},
			wantS: "CREATE TABLE `table` (\n" +
				"    `column` INT SIGNED NOT NULL,\n" +
				"    `column2` INT SIGNED NOT NULL\n" +
				"    PRIMARY KEY (`column`)\n);",
		},
	}

	m := &adapter{
		Base:      base.Base{Dialect: dialect.MySQL},
		validator: validator{},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotS := m.CreateTable(tt.tName, tt.t)
			if gotS != tt.wantS {
				t.Errorf("expected string `%s`, got string `%s`", tt.wantS, gotS)
			}
		})
	}
}

func Test_adapter_AddColumn(t *testing.T) {
	tests := map[string]struct {
		tName string
		cName string
		c     schema.Column
		wantS string
	}{
		"basic int column": {
			tName: "table",
			cName: "column",
			c: schema.Column{
				Datatype: datatype.Integer,
			},
			wantS: "ALTER TABLE `table` ADD COLUMN `column` INT SIGNED NOT NULL;",
		},
	}

	m := &adapter{
		Base:      base.Base{Dialect: dialect.MySQL},
		validator: validator{},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotS := m.AddColumn(tt.tName, tt.cName, tt.c)
			if gotS != tt.wantS {
				t.Errorf("expected string `%s`, got string `%s`", tt.wantS, gotS)
			}
		})
	}
}

func Test_adapter_AddIndex(t *testing.T) {
	tests := map[string]struct {
		tName string
		iName string
		i     schema.Index
		wantS string
	}{
		"non-unique single-column": {
			tName: "table",
			iName: "foreign",
			i: schema.Index{
				Columns: []string{"col"},
				Unique:  false,
			},
			wantS: "ALTER TABLE `table` ADD INDEX `foreign` (`col`);",
		},
		"non-unique two columns": {
			tName: "table",
			iName: "foreign",
			i: schema.Index{
				Columns: []string{"col", "col2"},
				Unique:  false,
			},
			wantS: "ALTER TABLE `table` ADD INDEX `foreign` (`col`, `col2`);",
		},
		"unique single-column": {
			tName: "table",
			iName: "foreign",
			i: schema.Index{
				Columns: []string{"col"},
				Unique:  true,
			},
			wantS: "ALTER TABLE `table` ADD UNIQUE INDEX `foreign` (`col`);",
		},
		"unique two columns": {
			tName: "table",
			iName: "foreign",
			i: schema.Index{
				Columns: []string{"col", "col2"},
				Unique:  true,
			},
			wantS: "ALTER TABLE `table` ADD UNIQUE INDEX `foreign` (`col`, `col2`);",
		},
	}

	m := &adapter{
		Base:      base.Base{Dialect: dialect.MySQL},
		validator: validator{},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotS := m.AddIndex(tt.tName, tt.iName, tt.i)
			if gotS != tt.wantS {
				t.Errorf("expected string `%s`, got string `%s`", tt.wantS, gotS)
			}
		})
	}
}

func Test_adapter_generateColumn(t *testing.T) {
	point := func(s string) *string {
		return &s
	}
	tests := map[string]struct {
		cName string
		c     schema.Column
		wantS string
	}{
		"int": {
			cName: "col",
			c: schema.Column{
				Datatype: datatype.Integer,
			},
			wantS: "`col` INT SIGNED NOT NULL",
		},
		"nullable int": {
			cName: "col",
			c: schema.Column{
				Datatype: datatype.Integer,
				Nullable: true,
			},
			wantS: "`col` INT SIGNED DEFAULT NULL NULL",
		},
		"int default 1": {
			cName: "col",
			c: schema.Column{
				Datatype: datatype.Integer,
				Default:  point("1"),
			},
			wantS: "`col` INT SIGNED DEFAULT 1 NOT NULL",
		},
		"unsigned int": {
			cName: "col",
			c: schema.Column{
				Datatype: datatype.Integer,
				Unsigned: true,
			},
			wantS: "`col` INT UNSIGNED NOT NULL",
		},
		"int auto_increment": {
			cName: "col",
			c: schema.Column{
				Datatype:      datatype.Integer,
				PrimaryKey:    true,
				AutoIncrement: true,
			},
			wantS: "`col` INT SIGNED NOT NULL AUTO_INCREMENT",
		},
		"decimal": {
			cName: "col",
			c: schema.Column{
				Datatype:  datatype.Decimal,
				Scale:     6,
				Precision: 4,
			},
			wantS: "`col` DECIMAL(6, 4) SIGNED NOT NULL",
		},
		"text default blah": {
			cName: "col",
			c: schema.Column{
				Datatype: datatype.Text,
				Default:  point("blah"),
			},
			wantS: "`col` TEXT DEFAULT \"blah\" NOT NULL",
		},
		"varchar": {
			cName: "col",
			c: schema.Column{
				Datatype: datatype.Varchar,
			},
			wantS: "`col` VARCHAR NOT NULL",
		},
		"sized varchar": {
			cName: "col",
			c: schema.Column{
				Datatype: datatype.Varchar,
				Scale:    64,
			},
			wantS: "`col` VARCHAR(64) NOT NULL",
		},
		"sized hypothetical thing": {
			cName: "col",
			c: schema.Column{
				Datatype:  0,
				Scale:     64,
				Precision: 43,
			},
			wantS: "`col` (64, 43) NOT NULL",
		},
	}

	m := &adapter{
		Base:      base.Base{Dialect: dialect.MySQL},
		validator: validator{},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotS := m.generateColumn(tt.cName, tt.c)
			if gotS != tt.wantS {
				t.Errorf("expected string `%s`, got string `%s`", tt.wantS, gotS)
			}
		})
	}
}

func Test_adapter_AddReference(t *testing.T) {
	tests := map[string]struct {
		tName  string
		ftName string
		fTable schema.Table
		r      schema.Reference
		wantS  string
	}{
		"simple single foreign key": {
			tName:  "local",
			ftName: "foreign",
			fTable: schema.Table{
				Columns: map[string]schema.Column{
					"id":       {PrimaryKey: true, Datatype: datatype.Integer},
					"otherCol": {Datatype: datatype.Integer},
				},
			},
			r: schema.Reference{
				Required: true,
			},
			wantS: "ALTER TABLE `local` ADD COLUMN `fk_foreign_id` INT SIGNED NOT NULL;\n" +
				"ALTER TABLE `local` ADD CONSTRAINT `reference_foreign` FOREIGN KEY (`fk_foreign_id`) REFERENCES foreign(`id`);",
		},
		"optional single foreign key": {
			tName:  "local",
			ftName: "foreign",
			fTable: schema.Table{
				Columns: map[string]schema.Column{
					"id":       {PrimaryKey: true, Datatype: datatype.Integer},
					"otherCol": {Datatype: datatype.Integer},
				},
			},
			wantS: "ALTER TABLE `local` ADD COLUMN `fk_foreign_id` INT SIGNED DEFAULT NULL NULL;\n" +
				"ALTER TABLE `local` ADD CONSTRAINT `reference_foreign` FOREIGN KEY (`fk_foreign_id`) REFERENCES foreign(`id`);",
		},
		"single foreign key with on delete": {
			tName:  "local",
			ftName: "foreign",
			fTable: schema.Table{
				Columns: map[string]schema.Column{
					"id":       {PrimaryKey: true, Datatype: datatype.Integer},
					"otherCol": {Datatype: datatype.Integer},
				},
			},
			r: schema.Reference{
				OnDelete: "CASCADE",
				Required: true,
			},
			wantS: "ALTER TABLE `local` ADD COLUMN `fk_foreign_id` INT SIGNED NOT NULL;\n" +
				"ALTER TABLE `local` ADD CONSTRAINT `reference_foreign` FOREIGN KEY (`fk_foreign_id`) REFERENCES foreign(`id`) ON DELETE CASCADE;",
		},
		"single foreign key with on update": {
			tName:  "local",
			ftName: "foreign",
			fTable: schema.Table{
				Columns: map[string]schema.Column{
					"id":       {PrimaryKey: true, Datatype: datatype.Integer},
					"otherCol": {Datatype: datatype.Integer},
				},
			},
			r: schema.Reference{
				OnUpdate: "CASCADE",
				Required: true,
			},
			wantS: "ALTER TABLE `local` ADD COLUMN `fk_foreign_id` INT SIGNED NOT NULL;\n" +
				"ALTER TABLE `local` ADD CONSTRAINT `reference_foreign` FOREIGN KEY (`fk_foreign_id`) REFERENCES foreign(`id`) ON UPDATE CASCADE;",
		},
		"dual foreign key": {
			tName:  "local",
			ftName: "foreign",
			fTable: schema.Table{
				Columns: map[string]schema.Column{
					"id":       {PrimaryKey: true, Datatype: datatype.Integer},
					"id2":      {PrimaryKey: true, Datatype: datatype.Integer},
					"otherCol": {Datatype: datatype.Integer},
				},
			},
			r: schema.Reference{
				Required: true,
			},
			wantS: "ALTER TABLE `local` ADD COLUMN `fk_foreign_id` INT SIGNED NOT NULL;\n" +
				"ALTER TABLE `local` ADD COLUMN `fk_foreign_id2` INT SIGNED NOT NULL;\n" +
				"ALTER TABLE `local` ADD CONSTRAINT `reference_foreign` FOREIGN KEY (`fk_foreign_id`, `fk_foreign_id2`) REFERENCES foreign(`id`, `id2`);",
		},
		"single foreign key with custom name": {
			tName:  "local",
			ftName: "foreign",
			fTable: schema.Table{
				Columns: map[string]schema.Column{
					"id": {PrimaryKey: true, Datatype: datatype.Integer},
				},
			},
			r: schema.Reference{
				ColumnNames: []string{"fk"},
				Required:    true,
			},
			wantS: "ALTER TABLE `local` ADD COLUMN `fk` INT SIGNED NOT NULL;\n" +
				"ALTER TABLE `local` ADD CONSTRAINT `reference_foreign` FOREIGN KEY (`fk`) REFERENCES foreign(`id`);",
		},
	}

	m := &adapter{
		Base:      base.Base{Dialect: dialect.MySQL},
		validator: validator{},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotS := m.AddReference(tt.tName, tt.ftName, tt.fTable, tt.r)
			if gotS != tt.wantS {
				t.Errorf("expected string `%s`, got string `%s`", tt.wantS, gotS)
			}
		})
	}
}
