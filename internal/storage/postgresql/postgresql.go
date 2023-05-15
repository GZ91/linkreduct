package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/storage/postgresql/postgresqlconfig"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
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
	ps             string
}

func New(config ConfigerStorage, generatorRunes GeneratorRunes) (*DB, error) {
	db := &DB{conf: config, generatorRunes: generatorRunes}
	ConfDB := db.conf.GetConfDB()
	db.ps = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		ConfDB.Address, ConfDB.User, ConfDB.Password, ConfDB.Dbname)
	err := createTable(db)
	return db, err
}

func createTable(db *DB) error {
	ctx := context.Background()
	dbs, err := sql.Open("pgx", db.ps)
	if err != nil {
		return err
	}
	defer dbs.Close()
	con, err := dbs.Conn(ctx)
	if err != nil {
		return err
	}
	defer con.Close()
	_, err = con.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS short_origin_reference 
(
	id serial PRIMARY KEY,
	uuid VARCHAR(45)  NOT NULL,
	ShortURL VARCHAR(250) NOT NULL,
	OriginalURL TEXT
);`)
	return err
}

func (d DB) Ping() error {
	ctx := context.Background()
	db, err := sql.Open("pgx", d.ps)
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
	ctx := context.Background()
	db, err := sql.Open("pgx", d.ps)
	if err != nil {
		logger.Log.Error("failed to open the database", zap.Error(err))
		return ""
	}
	defer db.Close()
	con, err := db.Conn(ctx)
	if err != nil {
		logger.Log.Error("failed to connect to the database", zap.Error(err))
		return ""
	}
	lenShort := d.conf.GetStartLenShortLink()
	index := 0
	var shorturl string
	for {
		shorturl = d.generatorRunes.RandStringRunes(lenShort)
		row := con.QueryRowContext(ctx, "SELECT COUNT(id) FROM short_origin_reference WHERE shorturl = $1", shorturl)
		var count_shorturl int
		row.Scan(&count_shorturl)
		if count_shorturl == 0 {
			break
		}
		index++
		if index == d.conf.GetMaxIterLen() {
			lenShort++
			index = 0
		}
	}

	_, err = con.ExecContext(ctx, "INSERT INTO short_origin_reference(uuid, shorturl, originalurl) VALUES ($1, $2, $3);", uuid.New().String(), shorturl, URL)
	if err != nil {
		logger.Log.Error("error when adding a record to the database", zap.Error(err))
	}
	return shorturl
}

func (d DB) GetURL(shortURL string) (string, bool) {
	ctx := context.Background()
	db, err := sql.Open("pgx", d.ps)
	if err != nil {
		logger.Log.Error("failed to open the database", zap.Error(err))
		return "", false
	}
	defer db.Close()
	con, err := db.Conn(ctx)
	if err != nil {
		logger.Log.Error("failed to connect to the database", zap.Error(err))
		return "", false
	}
	row := con.QueryRowContext(ctx, `SELECT originalurl
	FROM short_origin_reference WHERE shorturl = $1 limit 1`, shortURL)
	var originurl string
	row.Scan(&originurl)
	if originurl != "" {
		return originurl, true
	}
	return "", false
}

func (d DB) Close() error {
	return nil
}
