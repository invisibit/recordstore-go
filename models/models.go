package models

// import (
// 	// "database/sql"
// )

// // Models is the wrapper for database
// type Models struct {
// 	// DB DBModel
// }

// NewModels returns models with db pool
// func NewModels(db *sql.DB) Models {
// 	return Models{
// 		DB: DBModel{DB: db},
// 	}
// }

// type User struct {
// 	ID 			int		`json:"id"`
// 	FirstName	string 	`json:"first_name"`
// 	LastName	string 	`json:"last_name"`
// }

// type Cursors struct {
// 	After	string 	`json:"after"`
// 	Before	string 	`json:"before"`
// }