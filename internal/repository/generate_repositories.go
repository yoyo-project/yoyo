package repository

import (
	"fmt"
	"io"
	"strings"

	"github.com/yoyo-project/yoyo/internal/repository/template"

	"github.com/yoyo-project/yoyo/internal/schema"
)

func NewRepositoriesGenerator(packageName, reposPath string, packagePath Finder, db schema.Database) WriteGenerator {
	return func(db schema.Database, w io.StringWriter) (err error) {
		var reposStructFields, repoInits []string

		for _, t := range db.Tables {
			reposStructFields = append(reposStructFields, fmt.Sprintf("%sRepository", t.ExportedGoName()))
			repoInits = append(
				repoInits,
				fmt.Sprintf(
					"%sRepository: &%sRepository{baseRepo},",
					t.ExportedGoName(),
					t.ExportedGoName(),
				),
			)
		}

		r := strings.NewReplacer(
			template.PackageName,
			packageName,
			template.ReposStructFields,
			fmt.Sprintf("*%s", strings.Join(reposStructFields, "\n	*")),
			template.RepoInits,
			strings.Join(sortedUnique(repoInits), "\n		"),
		)

		_, err = w.WriteString(r.Replace(template.RepositoriesFile))

		return err
	}
}
