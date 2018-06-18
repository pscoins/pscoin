package server

import "os"

// Supported environments
const (
	Dev  = "dev"
	Test = "test"
	Acc  = "acc"
	Prod = "prod"
)

// Env resolves which environment the application is running on
func Env() string {
	r := Get("ENV")
	if r != Dev && r != Test && r != Acc && r != Prod {
		r = Dev
	}
	return r
}

// Get resolves a given key from the environment variables defined in the OS
func Get(key string) string {
	return os.Getenv(key)
}
