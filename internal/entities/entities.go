package entities

import "time"

type Dish struct {
	IDDish      uint8   `json:"id_dish"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float32 `json:"price"`
	IDCategory  string  `json:"id_category"`
}

type Order struct {
	IDOrder  uint8  `json:"id_order"`
	IDDesk   uint8  `json:"id_desk"`
	IDStatus string `json:"status"`
}

type Invoice struct {
	IDOrder uint8     `json:"id_order"`
	IDDish  uint8     `json:"id_dish"`
	Time    time.Time `json:"time_field"`
}

type Category struct {
	IDCategory string `json:"category"`
}

type Status struct {
	IDStatus string `json:"status"`
}
