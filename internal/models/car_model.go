package models

import "time"

type Car struct {
	ID        int64     `db:"id"`         // Поле id
	Marka     string    `db:"marka"`      // Поле marka
	Model     string    `db:"model"`      // Поле model
	VIN       string    `db:"vin"`        // Поле vin
	Year      int64     `db:"year"`       // Поле year
	Role      string    `db:"role"`       // Поле role (ENUM как строка)
	CreatedAt time.Time `db:"created_at"` // Поле created_at
	UpdatedAt time.Time `db:"updated_at"` // Поле updated_at
}
