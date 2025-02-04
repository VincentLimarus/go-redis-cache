package database

import "github.com/google/uuid"

type Students struct {
	ID 			uuid.UUID `gorm:"type:uuid;primaryKey;not null;default:uuid_generate_v4()"`
	Name 		string `gorm:"type:varchar(255);not null"`
	Email 		string `gorm:"type:varchar(255);not null"`
	Address 	string `gorm:"type:text;not null"` 
}