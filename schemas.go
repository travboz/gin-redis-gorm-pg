package main

type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Price int
}
