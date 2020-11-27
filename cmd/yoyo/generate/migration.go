package generate

//func MigrationGenerator(t time.Time) lime.Func {
//	return func(args []string) error {
//		config, err := yoyo.LoadConfig()
//		if err != nil {
//			return fmt.Errorf("unable to load config: %w", err)
//		}
//
//		generate, err := initSchemaGenerator(config)
//		if err != nil {
//			return fmt.Errorf("unable to initialize migration generator: %w", err)
//		}
//
//		sb := strings.Builder{}
//
//		err = generate(config.Schema, &sb)
//
//		if err != nil {
//			return fmt.Errorf("unable to generate migration: %w", err)
//		}
//
//		var name string
//		if len(args) > 0 {
//			name = fmt.Sprintf("_%s", strings.ToLower(strings.Join(args, "-")))
//		}
//		f, err := os.Create(fmt.Sprintf("%s/%s%s.sql", config.Paths.Migrations, t.Format("20060102150405"), name))
//		if err != nil {
//			return fmt.Errorf("cannot create migration file '%s': %w", config.Paths.Migrations, err)
//		}
//
//		_, err = f.WriteString(sb.String())
//		if err != nil {
//			return fmt.Errorf("cannot write to migration file '%s': %w", config.Paths.Migrations, err)
//		}
//
//		return nil
//	}
//}

//func initSchemaGenerator(config yoyo.Config) (migration.SchemaGenerator, error) {
//	var migrator migration.Dialect
//
//	reverser, err := reverse.LoadReverser(config.Schema.Dialect, env.DBHost(), env.DBUser(), env.DBName(), env.DBPassword(), env.DBPort())
//	if err != nil {
//		return nil, fmt.Errorf("cannot initialize dialect: %w", err)
//	}
//	migrator, err := migration.LoadDialect(config.Schema.Dialect)
//	if err != nil {
//		return nil, fmt.Errorf("cannot initialize dialect: %w", err)
//	}
//
//
//	return migration.NewSchemaGenerator(
//		migration.NewTableAdder(d),
//		migration.NewColumnAdder(d, migration.AddMissing, migration.NewColumnChecker(d)),
//		migration.NewIndexAdder(d, migration.AddMissing, migration.NewIndexChecker(d)),
//		migration.NewIndexAdder(d, migration.AddAll, nil),
//		migration.NewTableChecker(d),
//		migration.NewRefAdder(d, config.Schema, migration.AddMissing),
//		migration.NewRefAdder(d, config.Schema, migration.AddAll, nil),
//	), nil
//}
