package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Interface interface {
	orm.DB
}

type DB struct {
	Master Interface
}

func (d *DB) Close() {
	if masterConn, ok := d.Master.(*pg.DB); ok {
		_ = masterConn.Close()
	}
}
