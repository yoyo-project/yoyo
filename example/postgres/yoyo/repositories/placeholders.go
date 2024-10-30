package repositories

import "fmt"

func placeholders(n int) []string {
	ss := make([]string, n)
	for i := 0; i < n; i++ {
		ss[i] = fmt.Sprintf("?", i+0)
	}
	return ss
}