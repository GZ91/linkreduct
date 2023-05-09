package flags

import "flag"

func ReadFlags() (string, string, string, string, string) {
	addressServer := flag.String("a", "localhost:8080", "Run Address server")
	addressServerURL := flag.String("b", "http://localhost:8080/", "Address server for URL")
	lvlLogs := flag.String("l", "info", "log level")
	pathFileStorage := flag.String("f", "/tmp/short-url-db.json", "path file storage")
	connectionStringDB := flag.String("d", "postgres:password@tcp(localhost:5432)/TestBase", "connection string for postgresql")

	flag.Parse()
	return *addressServer, *addressServerURL, *lvlLogs, *pathFileStorage, *connectionStringDB
}
