package repository

import (
	"fmt"
	"github.com/yoyo-project/yoyo/internal/repository/template"
	"io"
	"strings"

	"github.com/yoyo-project/yoyo/internal/schema"
)

func NewRepositoriesGenerator(packageName, reposPath string, packagePath Finder, ts map[string]schema.Table) WriteGenerator {
	return func(db schema.Database, w io.StringWriter) (err error) {
		var repoInterfaces, imports, reposStructFields, repoInits []string

		for _, t := range ts {
			// replacer for interface templates
			ir := strings.NewReplacer(
				template.EntityName,
				t.ExportedGoName(),
				template.QueryPackageName,
				t.QueryPackageName(),
			)
			repoInterfaces = append(repoInterfaces, ir.Replace(template.RepositoryInterfaceTemplate))

			imp, err := packagePath(fmt.Sprintf(`%s/query/%s`, reposPath, t.QueryPackageName()))
			if err != nil {
				return fmt.Errorf("unable to generate repositories file: %w", err)
			}
			imports = append(imports, fmt.Sprintf(`"%s"`, imp))
			reposStructFields = append(reposStructFields, fmt.Sprintf("%sRepository", t.ExportedGoName()))
			repoInits = append(
				repoInits,
				fmt.Sprintf(
					"%sRepository: &%sRepo{baseRepo},",
					t.ExportedGoName(),
					t.QueryPackageName(),
				),
			)
		}

		r := strings.NewReplacer(
			template.PackageName,
			packageName,
			template.RepositoryInterfaces,
			strings.Join(repoInterfaces, "\n"),
			template.Imports,
			strings.Join(sortedUnique(imports), "\n	"),
			template.ReposStructFields,
			strings.Join(reposStructFields, "\n	"),
			template.RepoInits,
			strings.Join(repoInits, "\n		"),
		)

		_, err = w.WriteString(r.Replace(template.RepositoriesFile))

		return err
	}
}
