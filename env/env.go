package env

import "os"

// DBUser returns the environment value of YOYO_DB_USER
func DBUser() string {
	return os.Getenv("YOYO_DB_USER")
}

// DBPassword returns the environment value of YOYO_DB_PASSWORD
func DBPassword() string {
	return os.Getenv("YOYO_DB_PASSWORD")
}

// DBPort returns the environment value of YOYO_DB_PORT
func DBPort() string {
	return os.Getenv("YOYO_DB_PORT")
}

// DBHost returns the environment value of YOYO_DB_HOST
func DBHost() string {
	return os.Getenv("YOYO_DB_HOST")
}

// DBName returns the environment value of YOYO_DB_NAME
func DBName() string {
	return os.Getenv("YOYO_DB_NAME")
}

// DB returns the environment value of YOYO_DB
func DB() string {
	return os.Getenv("YOYO_DB")
}
