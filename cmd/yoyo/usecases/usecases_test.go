package usecases

import (
	"testing"
)

func TestInit(t *testing.T) {
	gotUCS := Init()

	switch {
	case gotUCS.LoadRepositoryGenerator == nil:
		t.Errorf("LoadRepositoryGenerator is nil")
	case gotUCS.LoadRepositoryAdapter == nil:
		t.Errorf("LoadRepositoryAdapter is nil")
	case gotUCS.LoadMigrationAdapter == nil:
		t.Errorf("LoadMigrationAdapter is nil")
	case gotUCS.ReadDatabase == nil:
		t.Errorf("ReadDatabase is nil")
	case gotUCS.GetCurrentTime == nil:
		t.Errorf("GetCurrentTime is nil")
	case gotUCS.BuildPostgresAdapter == nil:
		t.Errorf("BuildPostgresAdapter is nil")
	case gotUCS.BuildMySQLAdapter == nil:
		t.Errorf("BuildMySQLAdapter is nil")
	case gotUCS.LoadMigrationGenerator == nil:
		t.Errorf("LoadMigrationGenerator is nil")
	}
}
