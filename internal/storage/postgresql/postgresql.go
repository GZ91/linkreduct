package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GZ91/linkreduct/internal/storage/postgresql/postgresqlconfig"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type ConfigerStorage interface {
	GetMaxIterLen() int
	GetConfDB() *postgresqlconfig.ConfigDB
	GetStartLenShortLink() int
}

type GeneratorRunes interface {
	RandStringRunes(int) string
}

type DB struct {
	conf           ConfigerStorage
	generatorRunes GeneratorRunes
}

func New(config ConfigerStorage, generatorRunes GeneratorRunes) *DB {
	return &DB{conf: config, generatorRunes: generatorRunes}
}

func (d DB) Ping() error {
	ConfDB := d.conf.GetConfDB()
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		ConfDB.Address, ConfDB.User, ConfDB.Password, ConfDB.Dbname)
	ctx := context.Background()
	db, err := sql.Open("pgx", ps)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Conn(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (d DB) AddURL(URL string) string {
	return ""
}

func (d DB) GetURL(shortURL string) (string, bool) {

	return "", false
}

func (d DB) Close() error {
	return nil
}
