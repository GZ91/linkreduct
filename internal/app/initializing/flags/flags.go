package flags

import "flag"

// ReadFlags считывает параметры из флагов командной строки и возвращает их значения.
func ReadFlags() (string, string, string, string, string) {
	// Определение флагов с их значениями по умолчанию и описанием
	addressServer := flag.String("a", "localhost:8080", "Run Address server")
	addressServerURL := flag.String("b", "http://localhost:8080/", "Address server for URL")
	lvlLogs := flag.String("l", "info", "log level")
	pathFileStorage := flag.String("f", "", "path file storage")
	connectionStringDB := flag.String("d", "", "connection string for PostgreSQL")

	// Считывание значений флагов
	flag.Parse()
	return *addressServer, *addressServerURL, *lvlLogs, *pathFileStorage, *connectionStringDB
}
