package yoyo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/yoyo-project/yoyo/internal/schema"
	"gopkg.in/yaml.v3"
)

const (
	filename              = "yoyo.yml"
	defaultMigrationsPath = "yoyo/migrations/"
	defaultRepositoryPath = "yoyo/repositories/"
)

// LoadConfig searches for a yoyo.yml file, unmarshals it, and returns the unmarshaled struct.
// It travels toward the filesystem root, searching for yoyo.yml in each directory until
// it finds the file, or arrives at either / or <DriveLetter>:\
func LoadConfig() (yml Config, err error) {
	var (
		dir string
		f   []byte
	)

	dir, err = os.Getwd()
	if err != nil {
		return
	}

	for len(dir) >= 3 {
		f, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, filename))
		if err == nil {
			break
		}
		dir, _ = filepath.Split(strings.TrimRight(dir, "/"))
		dir = strings.TrimRight(dir, "/")
	}

	// landed at filesystem root, or something close enough to root to assume it's wrong
	if len(dir) <= 3 {
		err = errors.New("unable to find config file")
		return
	}

	err = yaml.Unmarshal(f, &yml)

	if yml.Paths.Migrations == "" {
		yml.Paths.Migrations = fmt.Sprintf("%s/%s", dir, defaultMigrationsPath)
	} else {
		yml.Paths.Migrations = fmt.Sprintf("%s/%s", dir, yml.Paths.Migrations)
	}

	if yml.Paths.Repositories == "" {
		yml.Paths.Repositories = fmt.Sprintf("%s/%s", dir, defaultRepositoryPath)
	} else {
		yml.Paths.Repositories = fmt.Sprintf("%s/%s", dir, yml.Paths.Repositories)
	}

	return
}

// Config is a struct which represents the yoyo.yml file
type Config struct {
	Paths  Paths
	Schema schema.Database
}

// Paths defines the locations that Migrations and generated Repositories code will be created in
type Paths struct {
	Migrations   string
	Repositories string // Soon...
	Models       string // Soon...
}
