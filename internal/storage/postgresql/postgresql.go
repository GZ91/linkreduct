package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
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
	db             *sql.DB
}

func New(ctx context.Context, config ConfigerStorage, generatorRunes GeneratorRunes) (*DB, error) {
	db := &DB{conf: config, generatorRunes: generatorRunes}
	ConfDB := db.conf.GetConfDB()
	db.ps = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		ConfDB.Address, ConfDB.User, ConfDB.Password, ConfDB.Dbname)
	err := db.openDB()
	if err != nil {
		return nil, err
	}
	err = db.createTable(ctx)
	if err != nil {
		return nil, err
	}
	return db, err
}

func (d *DB) openDB() error {
	db, err := sql.Open("pgx", d.ps)
	if err != nil {
		logger.Log.Error("failed to open the database", zap.Error(err))
		return err
	}
	d.db = db
	return nil
}

func (d *DB) createTable(ctx context.Context) error {
	con, err := d.db.Conn(ctx)
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

func (d *DB) Ping(ctx context.Context) error {
	con, err := d.db.Conn(ctx)
	if err != nil {
		logger.Log.Error("failed to connect to the database", zap.Error(err))
		return err
	}
	defer con.Close()
	return nil
}

func (d *DB) AddURL(ctx context.Context, URL string) (string, error) {
	con, err := d.db.Conn(ctx)
	if err != nil {
		logger.Log.Error("failed to connect to the database", zap.Error(err))
		return "", err
	}
	defer con.Close()
	lenShort := d.conf.GetStartLenShortLink()
	index := 0

	var shorturl string
	for {
		shorturl = d.generatorRunes.RandStringRunes(lenShort)
		row := con.QueryRowContext(ctx, "SELECT COUNT(id) FROM short_origin_reference WHERE shorturl = $1", shorturl)
		var countShorturl int
		err := row.Scan(&countShorturl)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Log.Error("when scanning a request for a shortcut", zap.Error(err))
			return "", err
		}
		if countShorturl == 0 {
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
	return shorturl, nil
}

func (d *DB) GetURL(ctx context.Context, shortURL string) (string, bool, error) {
	con, err := d.db.Conn(ctx)
	if err != nil {
		logger.Log.Error("failed to connect to the database", zap.Error(err))
		return "", false, err
	}
	defer con.Close()
	row := con.QueryRowContext(ctx, `SELECT originalurl
	FROM short_origin_reference WHERE shorturl = $1 limit 1`, shortURL)
	var originurl string
	err = row.Scan(&originurl)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error("when scanning the request for the original link", zap.Error(err))
		return "", false, err
	}
	if originurl != "" {
		return originurl, true, nil
	}
	return "", false, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) FindLongURL(ctx context.Context, OriginalURL string) (string, bool, error) {
	con, err := d.db.Conn(ctx)
	if err != nil {
		logger.Log.Error("failed to connect to the database", zap.Error(err))
		return "", false, err
	}
	defer con.Close()
	row := con.QueryRowContext(ctx, `SELECT ShortURL
	FROM short_origin_reference WHERE OriginalURL = $1 limit 1`, OriginalURL)
	var shortURL string
	err = row.Scan(&shortURL)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error("when scanning the result of the original link search query", zap.Error(err))
		return "", false, err
	}
	if shortURL != "" {
		return shortURL, true, nil
	}
	return "", false, nil
}

func (d *DB) AddBatchLink(ctx context.Context, batchLinks []models.IncomingBatchURL) (releasedBatchURL []models.ReleasedBatchURL, errs error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	reqShortURL, err := tx.PrepareContext(ctx, "SELECT COUNT(id) FROM short_origin_reference WHERE shorturl = $1")
	if err != nil {
		logger.Log.Error("When initializing a shortcut search pattern", zap.Error(err))
		return nil, err
	}
	reqLongLinkInBase, err := tx.PrepareContext(ctx, "SELECT shorturl FROM short_origin_reference WHERE originalurl = $1 LIMIT 1")
	if err != nil {
		logger.Log.Error("when initializing a long link search pattern", zap.Error(err))
		return nil, err
	}
	execInsertLongUrlInBase, err := tx.PrepareContext(ctx, "INSERT INTO short_origin_reference(uuid, shorturl, originalurl) VALUES ($1, $2, $3);")
	if err != nil {
		logger.Log.Error("when initializing the add string pattern", zap.Error(err))
		return nil, err
	}

	lenShort := d.conf.GetStartLenShortLink()
	index := 0
	for _, incomingLink := range batchLinks {
		shorturl := ""
		row := reqLongLinkInBase.QueryRowContext(ctx, incomingLink.OriginalURL)
		err := row.Scan(&shorturl)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logger.Log.Error("when searching for a scan of the result of finding a repeat of a long link", zap.Error(err))
				tx.Rollback()
				return nil, err
			}
		}
		if shorturl != "" {
			releasedBatchURL = append(releasedBatchURL, models.ReleasedBatchURL{
				CorrelationID: incomingLink.CorrelationID,
				ShortURL:      shorturl,
			})
			errs = fmt.Errorf("%w; %w", errs, errorsapp.ErrLinkAlreadyExists)
			continue
		}
		for {
			shorturl = d.generatorRunes.RandStringRunes(lenShort)
			row := reqShortURL.QueryRowContext(ctx, shorturl)
			var countShorturl int
			err := row.Scan(&countShorturl)
			if err != nil {
				logger.Log.Error("when searching for a scan of the result of finding a repeat of a short link", zap.Error(err))
				tx.Rollback()
				return nil, err
			}
			if countShorturl == 0 {
				break
			}
			index++
			if index == d.conf.GetMaxIterLen() {
				lenShort++
				index = 0
			}
		}

		_, err = execInsertLongUrlInBase.ExecContext(ctx, uuid.New().String(), shorturl, incomingLink.OriginalURL)
		if err != nil {
			logger.Log.Error("When creating a string with a long link in the database", zap.Error(err))
			tx.Rollback()
			return nil, err
		}
		releasedBatchURL = append(releasedBatchURL, models.ReleasedBatchURL{
			CorrelationID: incomingLink.CorrelationID,
			ShortURL:      shorturl,
		})
	}
	tx.Commit()
	return
}
