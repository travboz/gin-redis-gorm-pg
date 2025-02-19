package main

import "gorm.io/gorm"

type application struct {
	db *gorm.DB
}
