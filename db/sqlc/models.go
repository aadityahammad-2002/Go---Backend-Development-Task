package sqlc

import "time"

type User struct {
	ID   int64     `json:"id"`
	Name string    `json:"name"`
	Dob  time.Time `json:"dob"`
}
