package mysql

func (a *adapter) PreparedStatementPlaceholders(count int) []string {
	out := make([]string, count)
	for i := range out {
		out[i] = "?"
	}
	return out
}
