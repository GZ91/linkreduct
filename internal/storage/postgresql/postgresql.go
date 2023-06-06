package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/GZ91/linkreduct/internal/storage/postgresql/postgresqlconfig"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
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
	chsURLsForDel  chan []models.StructDelURLs
	chURLsForDel   chan models.StructDelURLs
}

func New(ctx context.Context, config ConfigerStorage, generatorRunes GeneratorRunes) (*DB, error) {
	db := &DB{conf: config, generatorRunes: generatorRunes}
	ConfDB := db.conf.GetConfDB()
	db.ps = ConfDB.StringServer
	err := db.openDB()
	if err != nil {
		return nil, err
	}
	err = db.createTable(ctx)
	if err != nil {
		return nil, err
	}
	db.chURLsForDel = make(chan models.StructDelURLs)
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
    userID VARCHAR(45)  NOT NULL,
    deletedFlag boolean DEFAULT FALSE, 
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
	var UserID string
	var userIDCTX models.CtxString = "userID"
	UserIDVal := ctx.Value(userIDCTX)
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	}

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

	_, err = con.ExecContext(ctx, "INSERT INTO short_origin_reference(uuid, shorturl, originalurl, userID) VALUES ($1, $2, $3, $4);",
		uuid.New().String(), shorturl, URL, UserID)
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
	row := con.QueryRowContext(ctx, `SELECT originalurl, deletedFlag 
	FROM short_origin_reference WHERE shorturl = $1 limit 1`, shortURL)
	var originurl string
	var deletedFlag bool
	err = row.Scan(&originurl, &deletedFlag)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error("when scanning the request for the original link", zap.Error(err))
		return "", false, err
	}
	if deletedFlag {
		return "", false, errorsapp.ErrLineURLDeleted
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
	var UserID string
	var userIDCTX models.CtxString = "userID"
	UserIDVal := ctx.Value(userIDCTX)
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	}

	tx, err := d.db.Begin()
	defer tx.Rollback()

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
	execInsertLongURLInBase, err := tx.PrepareContext(ctx, "INSERT INTO short_origin_reference(uuid, shorturl, originalurl, userID) VALUES ($1, $2, $3, $4);")
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
				return nil, err
			}
		}
		if shorturl != "" {
			releasedBatchURL = append(releasedBatchURL, models.ReleasedBatchURL{
				CorrelationID: incomingLink.CorrelationID,
				ShortURL:      shorturl,
			})
			errs = errorsapp.ErrLinkAlreadyExists
			continue
		}
		for {
			shorturl = d.generatorRunes.RandStringRunes(lenShort)
			row := reqShortURL.QueryRowContext(ctx, shorturl)
			var countShorturl int
			err := row.Scan(&countShorturl)
			if err != nil {
				logger.Log.Error("when searching for a scan of the result of finding a repeat of a short link", zap.Error(err))
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

		_, err = execInsertLongURLInBase.ExecContext(ctx, uuid.New().String(), shorturl, incomingLink.OriginalURL, UserID)
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
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return
}

func (d *DB) GetLinksUser(ctx context.Context, userID string) ([]models.ReturnedStructURL, error) {
	con, err := d.db.Conn(ctx)
	if err != nil {
		logger.Log.Error("failed to connect to the database", zap.Error(err))
		return nil, err
	}
	defer con.Close()

	rows, err := con.QueryContext(ctx, `SELECT ShortURL, OriginalURL
	FROM short_origin_reference WHERE userID = $1`, userID)
	if err != nil || rows.Err() != nil {
		if err != sql.ErrNoRows || rows.Err() != sql.ErrNoRows {
			logger.Log.Error("when reading data from the database", zap.Error(err))
			return nil, err
		}
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	returnData := make([]models.ReturnedStructURL, 0)

	for rows.Next() {
		var shortURL, originalURL string
		rows.Scan(&shortURL, &originalURL)
		returnData = append(returnData, models.ReturnedStructURL{OriginalURL: originalURL, ShortURL: shortURL})

	}
	return returnData, nil
}

func (d *DB) InitializingRemovalChannel(chsURLs chan []models.StructDelURLs) error {
	d.chsURLsForDel = chsURLs
	go d.GroupingDataForDeleted()
	go d.FillBufferDelete()
	return nil
}

func (d *DB) GroupingDataForDeleted() {
	for sliceVal := range d.chsURLsForDel {
		sliceVal := sliceVal
		go func() {
			for _, val := range sliceVal {
				d.chURLsForDel <- val
			}
		}()
	}
}

func (d *DB) FillBufferDelete() {
	t := time.NewTicker(time.Second * 10)
	var listForDel []models.StructDelURLs
	for {
		select {
		case val := <-d.chURLsForDel:
			listForDel = append(listForDel, val)
		case <-t.C:
			if len(listForDel) > 0 {
				d.deletedURLs(listForDel)
				listForDel = nil
			}
		}

	}
}

func (d *DB) deletedURLs(listForDel []models.StructDelURLs) {
	ctx := context.Background()
	tx, err := d.db.Begin()
	defer tx.Rollback()
	if err != nil {
		logger.Log.Error("error when trying to create a connection to the database", zap.Error(err))
		return
	}
	pr, err := tx.PrepareContext(ctx, "UPDATE short_origin_reference SET deletedFlag = true WHERE ShortURL = $1 and userID=$2")
	if err != nil {
		logger.Log.Error("error when trying to create a runtime request template", zap.Error(err))
		return
	}
	for _, val := range listForDel {
		pr.Exec(val.URL, val.UserID)
	}
	tx.Commit()
}
