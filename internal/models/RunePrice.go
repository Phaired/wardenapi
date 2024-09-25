package models

import (
	"time"
)

type RunePrice struct {
	Rune_name string    `json:"runeName"`
	Server    string    `json:"server"`
	Price     float64   `json:"price"`
	Date      time.Time `json:"date"`
}
