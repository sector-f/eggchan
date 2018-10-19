package main

import (
	"database/sql"
	"gopkg.in/guregu/null.v3"
	"time"
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

type board struct {
	ID          int         `json:"-"`
	Name        string      `json:"name"`
	Description null.String `json:"description"`
	Category    null.String `json:"category"`
}

func getBoards(db *sql.DB) ([]board, error) {
	rows, err := db.Query("SELECT boards.id, boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	boards := []board{}
	for rows.Next() {
		var b board
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Category); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	return boards, nil
}

func showCategory(db *sql.DB, name string) ([]board, error) {
	rows, err := db.Query("SELECT boards.id, boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id WHERE categories.name = $1", name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	boards := []board{}
	for rows.Next() {
		var b board
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Category); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	return boards, nil
}

type post struct {
	PostNum int       `json:"post_num"`
	ReplyTo null.Int  `json:"reply_to"`
	Time    time.Time `json:"time"`
	Comment string    `json:"comment"`
}

func showBoard(db *sql.DB, name string) ([]post, error) {
	rows, err := db.Query("SELECT posts.post_num, posts.reply_to, posts.time, posts.comment FROM posts LEFT JOIN boards ON boards.id = posts.board_id WHERE boards.name = $1 AND posts.reply_to IS NULL", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []post{}
	for rows.Next() {
		var p post
		if err := rows.Scan(&p.PostNum, &p.ReplyTo, &p.Time, &p.Comment); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func showThread(db *sql.DB, board string, thread int) ([]post, error) {
	rows, err := db.Query(
		"SELECT DISTINCT posts.post_num, posts.reply_to, posts.time, posts.comment FROM posts INNER JOIN boards ON (SELECT id FROM boards WHERE name = $1 LIMIT 1) = posts.board_id WHERE posts.post_num = $2 OR posts.reply_to = $2 ORDER BY posts.time ASC",
		board,
		thread,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []post{}
	for rows.Next() {
		var p post
		if err := rows.Scan(&p.PostNum, &p.ReplyTo, &p.Time, &p.Comment); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}
