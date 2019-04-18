package postgres

import (
	"database/sql"

	"github.com/sector-f/eggchan"
)

func (s *EggchanService) ShowBoardReply(name string) (eggchan.BoardReply, error) {
	// TODO: Use transaction for SELECTs here?

	board, err := s.ShowBoardDesc(name)
	if err != nil {
		return eggchan.BoardReply{}, err
	}

	posts, err := s.ShowBoard(name)
	if err != nil {
		return eggchan.BoardReply{}, err
	}

	return eggchan.BoardReply{board, posts}, nil
}

func (s *EggchanService) ShowThreadReply(name string, id int) (eggchan.ThreadReply, error) {
	// TODO: Use transaction for SELECTs here?

	board, err := s.ShowBoardDesc(name)
	if err != nil {
		return eggchan.ThreadReply{}, err
	}

	op, err := s.ShowThreadOP(name, id)
	if err != nil {
		return eggchan.ThreadReply{}, err
	}

	posts, err := s.ShowThread(name, id)
	if err != nil {
		return eggchan.ThreadReply{}, err
	}

	return eggchan.ThreadReply{board, op, posts}, nil
}

func (s *EggchanService) ListCategories() ([]eggchan.Category, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, eggchan.DatabaseError{}
	}
	defer tx.Commit()

	catRows, err := tx.Query("SELECT name FROM categories ORDER BY name ASC")

	if err != nil {
		return nil, eggchan.DatabaseError{}
	}

	categories := []eggchan.Category{}
	for catRows.Next() {
		var c eggchan.Category
		if err := catRows.Scan(&c.Name); err != nil {
			return nil, eggchan.DatabaseError{}
		}

		categories = append(categories, c)
	}
	catRows.Close()

	for i, category := range categories {
		boardRows, err := tx.Query("SELECT b.name, b.description, $1::text FROM boards b INNER JOIN categories c ON c.id = b.category WHERE c.name = $1::text", category.Name)
		if err != nil {
			return nil, eggchan.DatabaseError{}
		}

		boards := []eggchan.Board{}
		for boardRows.Next() {
			var b eggchan.Board
			if err := boardRows.Scan(&b.Name, &b.Description, &b.Category); err != nil {
				return nil, eggchan.DatabaseError{}
			}
			boards = append(boards, b)
		}
		boardRows.Close()

		categories[i].Boards = boards
	}

	return categories, nil
}

func (s *EggchanService) ListBoards() ([]eggchan.Board, error) {
	rows, err := s.DB.Query("SELECT boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id ORDER BY boards.name ASC")

	if err != nil {
		return nil, eggchan.DatabaseError{}
	}

	defer rows.Close()

	boards := []eggchan.Board{}
	for rows.Next() {
		var b eggchan.Board
		if err := rows.Scan(&b.Name, &b.Description, &b.Category); err != nil {
			return nil, eggchan.DatabaseError{}
		}
		boards = append(boards, b)
	}

	return boards, nil
}

func (s *EggchanService) ShowCategory(name string) ([]eggchan.Board, error) {
	var catExists int
	row := s.DB.QueryRow("SELECT count(1) FROM categories WHERE name = $1", name)
	err := row.Scan(&catExists)
	if err != nil {
		return nil, eggchan.DatabaseError{}
	}

	if catExists == 0 {
		return nil, eggchan.CategoryNotFoundError{}
	}

	rows, err := s.DB.Query("SELECT boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id WHERE categories.name = $1 ORDER BY boards.name ASC", name)
	if err != nil {
		return nil, eggchan.DatabaseError{}
	}

	defer rows.Close()

	boards := []eggchan.Board{}
	for rows.Next() {
		var b eggchan.Board
		if err := rows.Scan(&b.Name, &b.Description, &b.Category); err != nil {
			return nil, eggchan.DatabaseError{}
		}
		boards = append(boards, b)
	}

	return boards, nil
}

func (s *EggchanService) ShowBoard(name string) ([]eggchan.Thread, error) {
	var boardExists int
	row := s.DB.QueryRow("SELECT count(1) FROM boards WHERE name = $1", name)
	err := row.Scan(&boardExists)
	if err != nil {
		return nil, eggchan.DatabaseError{}
	}

	if boardExists == 0 {
		return nil, eggchan.BoardNotFoundError{}
	}

	rows, err := s.DB.Query(
		`SELECT
			$1::text,
			threads.post_num,
			threads.subject,
			threads.author,
			threads.time,
			(SELECT COUNT(*) FROM comments WHERE comments.reply_to = threads.id) AS num_replies,
			CASE
				WHEN MAX(comments.time) IS NOT NULL AND COUNT(*) >= (SELECT bump_limit FROM boards WHERE name = $1)  THEN (SELECT comments.time FROM comments OFFSET (SELECT bump_limit FROM boards WHERE name = $1) LIMIT 1)
				WHEN MAX(comments.time) IS NOT NULL THEN MAX(comments.time)
				ELSE MAX(threads.time)
			END AS sort_latest_reply,
			threads.comment
		FROM threads
		LEFT JOIN comments ON threads.id = comments.reply_to
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		GROUP BY threads.id
		ORDER BY sort_latest_reply DESC`,
		name,
	)

	if err != nil {
		return nil, eggchan.DatabaseError{}
	}
	defer rows.Close()

	threads := []eggchan.Thread{}
	for rows.Next() {
		var t eggchan.Thread
		if err := rows.Scan(&t.Board, &t.PostNum, &t.Subject, &t.Author, &t.Time, &t.NumReplies, &t.SortLatestReply, &t.Comment); err != nil {
			return threads, eggchan.DatabaseError{}
		}
		threads = append(threads, t)
	}

	return threads, nil
}

func (s *EggchanService) ShowThread(board string, thread_num int) ([]eggchan.Post, error) {
	var boardExists int
	row := s.DB.QueryRow("SELECT count(1) FROM boards WHERE name = $1", board)
	err := row.Scan(&boardExists)
	if err != nil {
		return nil, eggchan.DatabaseError{}
	}

	if boardExists == 0 {
		return nil, eggchan.BoardNotFoundError{}
	}

	var threadExists int
	row = s.DB.QueryRow(
		`SELECT count(1) FROM threads
		INNER JOIN boards ON boards.id = threads.board_id
		WHERE boards.name = $1
		AND threads.post_num = $2`,
		board,
		thread_num,
	)
	err = row.Scan(&threadExists)
	if err != nil {
		return nil, eggchan.DatabaseError{}
	}

	if threadExists == 0 {
		return nil, eggchan.ThreadNotFoundError{}
	}

	c_rows, err := s.DB.Query(
		`SELECT comments.reply_to, comments.post_num, comments.author, comments.time, comments.comment
		FROM comments
		INNER JOIN threads ON comments.reply_to = threads.id
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		AND comments.reply_to = (SELECT threads.id FROM threads INNER JOIN boards ON threads.board_id = boards.id WHERE boards.name = $1 AND threads.post_num = $2)
		ORDER BY post_num ASC`,
		board,
		thread_num,
	)

	posts := []eggchan.Post{}

	if err != nil {
		return posts, eggchan.DatabaseError{}
	}
	defer c_rows.Close()

	for c_rows.Next() {
		var p eggchan.Post
		if err := c_rows.Scan(&p.ReplyTo, &p.PostNum, &p.Author, &p.Time, &p.Comment); err != nil {
			return posts, eggchan.DatabaseError{}
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (s *EggchanService) MakeThread(board string, comment string, author string, subject string) (int, error) {
	var boardExists int
	row := s.DB.QueryRow("SELECT count(1) FROM boards WHERE name = $1", board)
	err := row.Scan(&boardExists)
	if err != nil {
		return 0, eggchan.DatabaseError{}
	}

	if boardExists == 0 {
		return 0, eggchan.BoardNotFoundError{}
	}

	rows, err := s.DB.Query(
		`INSERT INTO threads (board_id, comment, author, subject)
		VALUES(
			(SELECT id FROM boards WHERE name = $1),
			$2,
			$3,
			CASE WHEN $4 = '' THEN NULL ELSE $4 END
		)
		RETURNING post_num`,
		board,
		comment,
		author,
		subject,
	)

	if err != nil {
		return 0, eggchan.DatabaseError{}
	}

	post_nums := []int{}
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			return 0, eggchan.DatabaseError{}
		}
		post_nums = append(post_nums, i)
	}

	return post_nums[0], nil
}

func (s *EggchanService) MakeComment(board string, thread int, comment string, author string) (int, error) {
	var boardExists int
	row := s.DB.QueryRow("SELECT count(1) FROM boards WHERE name = $1", board)
	err := row.Scan(&boardExists)
	if err != nil {
		return 0, eggchan.DatabaseError{}
	}

	if boardExists == 0 {
		return 0, eggchan.BoardNotFoundError{}
	}

	var threadExists int
	row = s.DB.QueryRow(
		`SELECT count(1) FROM threads
		INNER JOIN boards ON boards.id = threads.board_id
		WHERE boards.name = $1
		AND threads.post_num = $2`,
		board,
		thread,
	)
	err = row.Scan(&threadExists)
	if err != nil {
		return 0, eggchan.DatabaseError{}
	}

	if threadExists == 0 {
		return 0, eggchan.ThreadNotFoundError{}
	}

	row = s.DB.QueryRow(
		`INSERT INTO comments (reply_to, comment, author)
		VALUES(
			(SELECT threads.id FROM threads INNER JOIN boards ON threads.board_id = boards.id WHERE boards.name = $1 AND threads.post_num = $2),
			$3,
			$4
		)
		RETURNING post_num`,
		board,
		thread,
		comment,
		author,
	)

	var post_num int

	err = row.Scan(&post_num)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		// TODO: make sure a nonexistent thread actually results in ErrNoRows
		return 0, eggchan.ThreadNotFoundError{}
	default:
		return 0, eggchan.DatabaseError{}
	}

	return post_num, nil
}

func (s *EggchanService) ShowBoardDesc(board string) (eggchan.Board, error) {
	var b eggchan.Board

	var boardExists int
	row := s.DB.QueryRow("SELECT count(1) FROM boards WHERE name = $1", board)
	err := row.Scan(&boardExists)
	if err != nil {
		return b, eggchan.DatabaseError{}
	}

	if boardExists == 0 {
		return b, eggchan.BoardNotFoundError{}
	}

	b_row := s.DB.QueryRow(
		`SELECT boards.name, boards.description, boards.category
		FROM boards
		WHERE boards.name = $1`,
		board,
	)

	err = b_row.Scan(&b.Name, &b.Description, &b.Category)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return b, eggchan.BoardNotFoundError{}
	default:
		return b, eggchan.DatabaseError{}
	}

	return b, nil
}

func (s *EggchanService) ShowThreadOP(board string, id int) (eggchan.Thread, error) {
	var t eggchan.Thread

	var boardExists int
	row := s.DB.QueryRow("SELECT count(1) FROM boards WHERE name = $1", board)
	err := row.Scan(&boardExists)
	if err != nil {
		return t, eggchan.DatabaseError{}
	}

	if boardExists == 0 {
		return t, eggchan.BoardNotFoundError{}
	}

	var threadExists int
	row = s.DB.QueryRow(
		`SELECT count(1) FROM threads
		INNER JOIN boards ON boards.id = threads.board_id
		WHERE boards.name = $1
		AND threads.post_num = $2`,
		board,
		id,
	)
	err = row.Scan(&threadExists)
	if err != nil {
		return t, eggchan.DatabaseError{}
	}

	if threadExists == 0 {
		return t, eggchan.ThreadNotFoundError{}
	}

	t_row := s.DB.QueryRow(
		`SELECT
			$1::text,
			threads.post_num,
			threads.subject,
			threads.author,
			threads.time,
			(SELECT COUNT(*) FROM comments WHERE comments.reply_to = threads.id) AS num_replies,
			CASE
				WHEN MAX(comments.time) IS NOT NULL AND COUNT(*) >= (SELECT bump_limit FROM boards WHERE name = $1)  THEN (SELECT comments.time FROM comments OFFSET (SELECT bump_limit FROM boards WHERE name = $1) LIMIT 1)
				WHEN MAX(comments.time) IS NOT NULL THEN MAX(comments.time)
				ELSE MAX(threads.time)
			END AS sort_latest_reply,
			threads.comment
		FROM threads
		LEFT JOIN comments ON threads.id = comments.reply_to
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		AND threads.post_num = $2
		GROUP BY threads.id
		ORDER BY sort_latest_reply DESC`,
		board,
		id,
	)

	err = t_row.Scan(&t.Board, &t.PostNum, &t.Subject, &t.Author, &t.Time, &t.NumReplies, &t.SortLatestReply, &t.Comment)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return t, eggchan.ThreadNotFoundError{}
	default:
		return t, eggchan.DatabaseError{}
	}

	return t, nil
}

func (s *EggchanService) AddCategory(category string) error {
	_, err := s.DB.Exec(
		`INSERT INTO categories (name) VALUES ($1)`,
		category,
	)

	if err != nil {
		return eggchan.DatabaseError{}
	}

	return nil
}

func (s *EggchanService) AddBoard(board, description, category string) error {
	_, err := s.DB.Exec(
		`INSERT INTO boards (name, description, category) VALUES ($1, $2, (SELECT id FROM categories WHERE name = $3))`,
		board,
		description,
		category,
	)

	if err != nil {
		return eggchan.DatabaseError{}
	}

	return nil
}
