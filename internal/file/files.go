package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateWithDirs(path string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}
	return os.Create(path)
}

func FindPackagePath(path string) (string, error) {
	gopath := os.Getenv("GOPATH")
	var str, p string
	for p = path; len(p) > 4; p = filepath.Dir(p) {
		str = func(p string) string {
			f, err := os.Open(fmt.Sprintf("%s/go.mod", p))
			defer func() {
				if f != nil {
					_ = f.Close()
				}
			}()
			if !os.IsNotExist(err) {
				scanner := bufio.NewScanner(f)

				var line string
				for scanner.Scan() {
					line = scanner.Text()
					if line[0:6] == "module" {
						return fmt.Sprintf("%s%s", strings.Trim(line[7:], `"'`), path[len(p):])
					}
				}
			}

			// Trim the strings of slashes so we don't need to worry about trailing slashes messing up the comparison
			if strings.Trim(p, "\\/") == strings.Trim(gopath, "\\/") {
				return path[len(p):]
			}

			return ""
		}(p)

		if len(str) > 0 {
			break
		}
	}

	if len(str) > 0 {
		return str, nil
	}

	return "", fmt.Errorf("could not determine package path")
}
