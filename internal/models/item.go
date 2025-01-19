package models

import (
	"time"
)

type Item struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	Price      float64   `json:"price"`
	CreateDate time.Time `json:"create_date"`
}
