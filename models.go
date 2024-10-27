package main

import "time"

type URL struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	Tag       string `gorm:"size:3;index;not null;unique"`
	URL       string `gorm:"not null;unique"`
}
