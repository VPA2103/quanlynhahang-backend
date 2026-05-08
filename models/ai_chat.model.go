package models

type AIChatResponse struct {
	Type    string      `json:"type"` // "text" | "food_list"
	Message string      `json:"message,omitempty"`
	Foods   []FoodChat  `json:"foods,omitempty"`
}

type FoodChat struct {
	MaMonAn int     `json:"ma_mon_an"`
	Ten     string  `json:"ten_mon_an"`
	Gia     float64 `json:"gia_tien"`
	Anh     string  `json:"anh_mon_an"`
}