// Package models описывает схему JSON ответ после загрузки zip-архив с данными
package models

type TotalPrice struct {
	TotalItems      int     `json:"total_items"`      //общее количество добавленных элементов
	TotalCategories int     `json:"total_categories"` //общее количество категорий
	TotalPrice      float64 `json:"total_price"`      //суммарная стоимость всех объектов в базе данных
}
