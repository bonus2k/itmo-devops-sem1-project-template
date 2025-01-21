package repositories

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"project_sem/internal/logger"
	"project_sem/internal/models"
	"sync"
)

type DataStorable interface {
	AddItem(ctx context.Context, item *models.Item) error
	GetAllItems(ctx context.Context) (*[]models.Item, error)
}

var (
	log   *logger.Logger
	once  sync.Once
	ds    DataStorable
	dbErr error
)

type DataStore struct {
	db *sql.DB
}

func (d DataStore) AddItem(ctx context.Context, item *models.Item) error {
	sqlStm := `INSERT INTO prices (id, name, category, price, create_date) 
							VALUES ($1, $2, $3, $4, $5)`

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, sqlStm)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, item.Id, item.Name, item.Category, item.Price, item.CreateDate)
	defer func(stmt *sql.Stmt) {
		errStm := stmt.Close()
		if errStm != nil {
			log.WithError(errStm).Error("Close the statement")
		}
	}(stmt)

	if err != nil {
		errRB := tx.Rollback()
		if errRB != nil {
			log.WithError(errRB).Error("Rollback has filed")
		}
		log.WithError(err).Error("Insert the price")
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (d DataStore) GetAllItems(ctx context.Context) (*[]models.Item, error) {
	sqlStm := `SELECT id, name, category, price, create_date FROM prices`

	stmt, err := d.db.PrepareContext(ctx, sqlStm)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		errStm := stmt.Close()
		if errStm != nil {
			log.WithError(errStm).Error("Close the statement")
		}
	}(stmt)

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]models.Item, 0)
	for rows.Next() {
		item := models.Item{}
		err := rows.Scan(&item.Id, &item.Name, &item.Category, &item.Price, &item.CreateDate)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &items, nil
}

func NewDataStore(logger *logger.Logger, connect string) (*DataStorable, error) {
	once.Do(
		func() {
			log = logger
			var dataBase *sql.DB
			log.Info("Init DB connection")
			dataBase, dbErr = sql.Open("pgx/v5", connect)
			if dbErr != nil {
				return
			}
			dbErr = dataBase.Ping()
			ds = &DataStore{db: dataBase}
		},
	)
	return &ds, dbErr
}
