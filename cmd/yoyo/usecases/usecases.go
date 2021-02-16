package usecases

import (
	"database/sql"
	"time"

	"github.com/yoyo-project/yoyo/internal/repository"

	"github.com/yoyo-project/yoyo/internal/dbms/mysql"
	"github.com/yoyo-project/yoyo/internal/dbms/postgres"
	"github.com/yoyo-project/yoyo/internal/migration"
	"github.com/yoyo-project/yoyo/internal/reverse"
)

type UseCases struct {
	GetCurrentTime       func() time.Time
	BuildMySQLAdapter    reverse.AdapterBuilder
	BuildPostgresAdapter reverse.AdapterBuilder
	LoadReverseAdapter   reverse.AdapterLoader
	ReadDatabase         reverse.DatabaseReader

	LoadMigrationAdapter   migration.AdapterLoader
	LoadMigrationGenerator migration.GeneratorLoader

	LoadRepositoryAdapter   repository.AdapterLoader
	LoadRepositoryGenerator repository.GeneratorLoader
}

func Init() (ucs UseCases) {
	ucs = UseCases{
		GetCurrentTime:       time.Now,
		LoadMigrationAdapter: migration.LoadAdapter,
		BuildMySQLAdapter:    mysql.InitReverserBuilder(sql.Open),
		BuildPostgresAdapter: postgres.InitReverserBuilder(sql.Open),
	}

	ucs.LoadReverseAdapter = reverse.InitAdapterSelector(ucs.BuildMySQLAdapter, ucs.BuildPostgresAdapter)
	ucs.ReadDatabase = reverse.InitDatabaseReader(ucs.LoadReverseAdapter)
	ucs.LoadMigrationGenerator = migration.InitGeneratorLoader(ucs.LoadReverseAdapter, ucs.LoadMigrationAdapter, migration.NewGenerator)

	ucs.LoadRepositoryAdapter = repository.LoadAdapter
	ucs.LoadRepositoryGenerator = repository.InitGeneratorLoader(repository.NewGenerator, ucs.LoadRepositoryAdapter)

	return ucs
}
