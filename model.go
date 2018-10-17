package main

import (
	"database/sql"
)

type category struct {
	ID   int    `json:"-"`
	Name string `json:"name"`
}

func (c *category) getCategory(db *sql.DB) error {
	return db.QueryRow("SELECT id, name FROM categories WHERE id=$1", c.ID).Scan(&c.ID, &c.Name)
}

// func getCategories(db *sql.DB, start, count int) ([]category, error) {
func getCategories(db *sql.DB) ([]category, error) {
	rows, err := db.Query("SELECT id, name FROM categories")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := []category{}
	for rows.Next() {
		var c category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}
