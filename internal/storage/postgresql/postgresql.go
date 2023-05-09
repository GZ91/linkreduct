package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	address  string
	user     string
	password string
	dbname   string
}

func New(address, user, password, dbname string) *DB {
	return &DB{address: address, user: user, password: password, dbname: dbname}
}

func (d DB) Ping() error {
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		d.address, d.user, d.password, d.dbname)
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
