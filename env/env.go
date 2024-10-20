package env

import "os"

func DBURL() string { return os.Getenv("YOYO_DB_URL") }
