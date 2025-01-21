package repositories

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"project_sem/internal/logger"
	"project_sem/internal/models"
	"sync"
)

type DataStorable interface {
	AddItems(ctx context.Context, item *[]models.Item) (int, error)
	GetAllItems(ctx context.Context) (*[]models.Item, error)
	GetStatisticItems(ctx context.Context) (*models.TotalPrice, error)
	Close() error
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

func (d DataStore) AddItems(ctx context.Context, items *[]models.Item) (int, error) {
	count := 0
	sqlStm := `INSERT INTO prices (id, name, category, price, create_date) 
							VALUES ($1, $2, $3, $4, $5)`

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return count, err
	}

	stmt, err := tx.PrepareContext(ctx, sqlStm)
	defer func(stmt *sql.Stmt) {
		errStm := stmt.Close()
		if errStm != nil {
			log.WithError(errStm).Error("Close the statement")
		}
	}(stmt)

	if err != nil {
		return count, err
	}

	for _, item := range *items {
		_, err = stmt.ExecContext(ctx, item.Id, item.Name, item.Category, item.Price, item.CreateDate)

		if err != nil {
			errRB := tx.Rollback()
			if errRB != nil {
				log.WithError(errRB).Error("Rollback has filed")
			}
			log.WithError(err).Error("Insert the price")
			count = 0
			return count, err
		}

		count++
	}

	err = tx.Commit()
	if err != nil {
		return count, err
	}

	return count, nil
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

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &items, nil
}

func (d DataStore) GetStatisticItems(ctx context.Context) (*models.TotalPrice, error) {
	sqlStm := `SELECT COUNT(DISTINCT (p.category)) AS total_categories, SUM(p.price) AS total_price FROM prices p `

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

	row := stmt.QueryRowContext(ctx)
	totalPrice := models.TotalPrice{}

	err = row.Scan(&totalPrice.TotalCategories, &totalPrice.TotalPrice)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &totalPrice, nil
		}
		return nil, err
	}

	err = row.Err()
	if err != nil {
		return nil, err
	}

	return &totalPrice, nil
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

func (d DataStore) Close() error {
	return d.Close()
}
