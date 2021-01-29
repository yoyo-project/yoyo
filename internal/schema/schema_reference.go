package schema

import (
	"fmt"
)

const (
	foreignKeyPrefix = "fk_"
)

// ColNames returns a list foreign key column names for the given reference. The method assumes that the `fTable` and
// `ftName` values are correct for the foreign table (for example, the "many" side in a one-to-many reference) and not
// necessarily the target of the Reference itself. So make sure the Yoyo-to-RDB reference translation has already happened
// before calling ColNames.
func (r *Reference) ColNames(ftName string, fTable Table) []string {
	var (
		fknames     []string
		fkname      string
		refColNames = r.ColumnNames
	)

	for _, fcName := range fTable.PKColNames() {
		switch {
		case len(refColNames) > 0:
			fkname, refColNames = refColNames[0], refColNames[1:]
		default:
			fkname = fmt.Sprintf("%s%s_%s", foreignKeyPrefix, ftName, fcName)
		}

		fknames = append(fknames, fkname)
	}

	return fknames
}

// ExportedGoFields returns a list of names for Go Fields which represent the columns for the given reference. The method
// assumes that the `fTable` and `ftName` values are correct for the foreign table (for example, the "many" side in a
// one-to-many reference) and not necessarily the target of the Reference itself. So make sure the Yoyo-to-RDB reference
// translation has already happened before calling ExportedGoFields.
func (r *Reference) ExportedGoFields(ftName string, fTable Table) []string {
	cNames := r.ColNames(ftName, fTable)
	fNames := make([]string, len(cNames))

	for i := range cNames {
		if cNames[i][0:len(foreignKeyPrefix)] == foreignKeyPrefix {
			cNames = cNames[len(foreignKeyPrefix):]
		}
		fNames = append(fNames, pascal(cNames[i]))
	}

	return fNames
}
