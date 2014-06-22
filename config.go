package main

/*
Config is used for storing data retrieved from environmental variables.
*/
type Config struct {
	Port      int
	LogLevel  string
	LogFormat string
	APIToken  string
	SkipPush  bool
}
