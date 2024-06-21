package entity

type Reserve struct {
	UserID int    `json:"user_id"`
	HallID int    `json:"hall_id"`
	ID     int    `json:"id"`
	Date   string `json:"date"`
	Time   int    `json:"time"`
	Status int    `json:"status"`
}

type ReserveReq struct {
	HallID int    `json:"hall_id"`
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Time   int    `json:"time"`
}
