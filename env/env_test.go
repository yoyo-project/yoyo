package env

import (
	"os"
	"testing"
)

func TestEnvs(t *testing.T) {
	os.Setenv("YOYO_DB_USER", "test")
	os.Setenv("YOYO_DB_PASSWORD", "test")
	os.Setenv("YOYO_DB_PORT", "test")
	os.Setenv("YOYO_DB_HOST", "test")
	os.Setenv("YOYO_DB_NAME", "test")
	os.Setenv("YOYO_DB", "test")

	if DBUser() != "test" {
		t.Errorf("DBUser result unexpected")
	}
	if DBPassword() != "test" {
		t.Errorf("DBPassword result unexpected")
	}
	if DBPort() != "test" {
		t.Errorf("DBPort result unexpected")
	}
	if DBHost() != "test" {
		t.Errorf("DBHost result unexpected")
	}
	if DBName() != "test" {
		t.Errorf("DBName result unexpected")
	}
	if DB() != "test" {
		t.Errorf("DB result unexpected")
	}
}
