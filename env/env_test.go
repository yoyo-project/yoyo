package env

import (
	"os"
	"testing"
)

func TestEnvs(t *testing.T) {
	os.Setenv("YOYO_DB_URL", "test")
	if DBURL() != "test" {
		t.Errorf("DBName result unexpected")
	}
}
