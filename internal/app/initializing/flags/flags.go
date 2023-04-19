package flags

import "flag"

func ReadFlags() (string, string, string) {
	addressServer := flag.String("a", "localhost:8080", "Run Address server")
	addressServerURL := flag.String("b", "http://localhost:8080/", "Address server for URL")
	lvlLogs := flag.String("l", "info", "log level")

	flag.Parse()
	return *addressServer, *addressServerURL, *lvlLogs
}
