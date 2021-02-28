package schema

import (
	"fmt"
)

const (
	foreignKeyPrefix = "fk_"
)

// ColNames returns a list foreign key column names for the given reference. The method assumes that the `fTable` value
//is correct (for example, the "many" side in a one-to-many reference) and not necessarily the target table of the
// Reference itself. So make sure the Yoyo-to-RDB reference translation has already happened before calling ColNames.
func (r *Reference) ColNames(ft Table) []string {
	var (
		fknames     []string
		fkname      string
		refColNames = r.ColumnNames
	)

	for _, fcName := range ft.PKColNames() {
		switch {
		case len(refColNames) > 0:
			fkname, refColNames = refColNames[0], refColNames[1:]
		default:
			fkname = fmt.Sprintf("%s%s_%s", foreignKeyPrefix, ft.Name, fcName)
		}

		fknames = append(fknames, fkname)
	}

	return fknames
}
