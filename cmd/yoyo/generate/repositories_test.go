package generate

import (
	"os"
	"testing"

	"github.com/yoyo-project/yoyo/cmd/yoyo/usecases"
)

func TestRepos(t *testing.T) {
	oldwd, _ := os.Getwd()
	t.Cleanup(func() {
		os.Chdir(oldwd)
	})

	os.Chdir("../../../example/")

	ucs := usecases.Init()
	err := Repos(ucs.LoadRepositoryGenerator)([]string{})
	if err != nil {
		t.Errorf("error when generating repositories")
	}
}
