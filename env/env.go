package env

import "os"

func DBUser() string {
	return os.Getenv("YOYO_DB_USER")
}

func DBPassword() string {
	return os.Getenv("YOYO_DB_PASSWORD")
}

func DBPort() string {
	return os.Getenv("YOYO_DB_PORT")
}

func DBHost() string {
	return os.Getenv("YOYO_DB_HOST")
}

func DBName() string {
	return os.Getenv("YOYO_DB_NAME")
}

func DB() string {
	return os.Getenv("YOYO_DB")
}
