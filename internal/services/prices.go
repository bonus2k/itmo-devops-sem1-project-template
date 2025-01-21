package services

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"math"
	"project_sem/internal/logger"
	"project_sem/internal/models"
	"project_sem/internal/repositories"
	"strconv"
	"sync"
	"time"
)

var (
	log  *logger.Logger
	once sync.Once
	ps   PriceService
)

type PriceServiceImpl struct {
	ds repositories.DataStorable
}

type PriceService interface {
	SaveItem(ctx context.Context, items [][]string) (*models.TotalPrice, error)
	AllItem(ctx context.Context) ([][]string, error)
}

func NewPriceService(logger *logger.Logger, ds *repositories.DataStorable) *PriceService {
	once.Do(func() {
		log = logger
		ps = PriceServiceImpl{ds: *ds}
	})
	return &ps
}

func (p PriceServiceImpl) SaveItem(ctx context.Context, items [][]string) (*models.TotalPrice, error) {
	reqID := middleware.GetReqID(ctx)
	newItems := make([]models.Item, 0)
	for i, item := range items {
		if i == 0 {
			continue
		}

		parseId, err := strconv.ParseInt(item[id], 10, 64)
		if err != nil {
			log.WithField("reqID", reqID).
				WithField("item", item).
				WithError(err).Error("Convert Id to int")
			continue
		}

		parsePrice, err := strconv.ParseFloat(item[price], 32)
		if err != nil {
			log.WithField("reqID", reqID).
				WithField("item", item).
				WithError(err).Error("Convert Price to float")
			continue
		}
		parsePrice = math.Round(parsePrice*100) / 100

		parseDate, err := time.Parse(time.DateOnly, item[create_date])
		if err != nil {
			log.WithField("reqID", reqID).
				WithField("item", item).
				WithError(err).Error("Convert CreateDate to date")
			continue
		}

		newItem := models.Item{
			Id:         parseId,
			Name:       item[name],
			Category:   item[category],
			Price:      parsePrice,
			CreateDate: parseDate,
		}

		newItems = append(newItems, newItem)
	}

	totalItems, err := p.ds.AddItems(ctx, &newItems)
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("Save the items to database")
	}

	statisticItems, err := p.ds.GetStatisticItems(ctx)
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("Get statistical data from database")
	}

	statisticItems.TotalItems = totalItems
	return statisticItems, nil
}

func (p PriceServiceImpl) AllItem(ctx context.Context) ([][]string, error) {
	reqID := middleware.GetReqID(ctx)
	items := make([][]string, 0)
	allItems, err := p.ds.GetAllItems(ctx)
	if err != nil {
		log.WithField("reqID", reqID).
			WithError(err).Error("Get all item from db")
		return nil, err
	}

	for i, item := range *allItems {
		if i == 0 {
			newItem := []string{"id", "name", "category", "price", "create_date"}
			items = append(items, newItem)
		}
		newItem := make([]string, 5)
		newItem[id] = strconv.FormatInt(item.Id, 10)
		newItem[name] = item.Name
		newItem[category] = item.Category
		newItem[price] = fmt.Sprintf("%.2f", item.Price)
		newItem[create_date] = item.CreateDate.Format(time.DateOnly)
		items = append(items, newItem)
	}
	return items, nil
}

const (
	id = iota
	name
	category
	price
	create_date
)
