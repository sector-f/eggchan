package postgres

import (
	"database/sql"
	"errors"

	"github.com/sector-f/eggchan"
	"golang.org/x/crypto/bcrypt"
)

type EggchanService struct {
	DB *sql.DB
}

func (s *EggchanService) ListCategories() ([]eggchan.Category, error) {
	rows, err := s.DB.Query("SELECT name FROM categories ORDER BY name ASC")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := []eggchan.Category{}
	for rows.Next() {
		var c eggchan.Category
		if err := rows.Scan(&c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (s *EggchanService) ListBoards() ([]eggchan.Board, error) {
	rows, err := s.DB.Query("SELECT boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id ORDER BY boards.name ASC")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	boards := []eggchan.Board{}
	for rows.Next() {
		var b eggchan.Board
		if err := rows.Scan(&b.Name, &b.Description, &b.Category); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	return boards, nil
}

func (s *EggchanService) ShowCategory(name string) ([]eggchan.Board, error) {
	rows, err := s.DB.Query("SELECT boards.name, boards.description, categories.name FROM boards LEFT JOIN categories ON boards.category = categories.id WHERE categories.name = $1 ORDER BY boards.name ASC", name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	boards := []eggchan.Board{}
	for rows.Next() {
		var b eggchan.Board
		if err := rows.Scan(&b.Name, &b.Description, &b.Category); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	return boards, nil
}

func (s *EggchanService) ShowBoard(name string) ([]eggchan.Thread, error) {
	rows, err := s.DB.Query(
		`SELECT
			boards.name,
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
		INNER JOIN boards ON threads.board_id = boards.id
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		GROUP BY boards.name, threads.id
		ORDER BY sort_latest_reply DESC`,
		name,
	)

	threads := []eggchan.Thread{}

	if err != nil {
		return threads, err
	}
	defer rows.Close()

	for rows.Next() {
		var t eggchan.Thread
		if err := rows.Scan(&t.Board, &t.PostNum, &t.Subject, &t.Author, &t.Time, &t.NumReplies, &t.SortLatestReply, &t.Comment); err != nil {
			return threads, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}

func (s *EggchanService) ShowThread(board string, thread_num int) ([]eggchan.Post, error) {
	c_rows, err := s.DB.Query(
		`SELECT comments.post_num, comments.author, comments.time, comments.comment
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
		return posts, err
	}
	defer c_rows.Close()

	for c_rows.Next() {
		var p eggchan.Post
		if err := c_rows.Scan(&p.PostNum, &p.Author, &p.Time, &p.Comment); err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (s *EggchanService) MakeThread(board string, comment string, author string, subject string) (int, error) {
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
		return 0, err
	}

	post_nums := []int{}
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			return 0, err
		}
		post_nums = append(post_nums, i)
	}

	return post_nums[0], nil
}

func (s *EggchanService) MakeComment(board string, thread int, comment string, author string) (int, error) {
	// TODO: use QueryRow here
	rows, err := s.DB.Query(
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

	if err != nil {
		return 0, err
	}

	post_nums := []int{}
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			return 0, err
		}
		post_nums = append(post_nums, i)
	}

	return post_nums[0], nil
}

func (s *EggchanService) checkIsOp(board string, thread int) (bool, error) {
	rows, err := s.DB.Query(
		`SELECT post_num
		FROM threads
		INNER JOIN boards ON threads.board_id = boards.id
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		AND threads.post_num = $2`,
		board,
		thread,
	)

	if err != nil {
		return false, err
	}

	posts := []eggchan.Post{}
	for rows.Next() {
		var p eggchan.Post
		if err := rows.Scan(&p.PostNum); err != nil {
			return false, err
		}
		posts = append(posts, p)
	}

	if len(posts) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (s *EggchanService) getUserAuthentication(name string, pw []byte) (bool, error) {
	pw_row := s.DB.QueryRow(`SELECT password FROM users WHERE username = $1`, name)
	var db_pw []byte
	if err := pw_row.Scan(&db_pw); err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword(db_pw, pw); err != nil {
		return false, nil
	} else {
		return true, nil
	}
}

func (s *EggchanService) getUserAuthorization(name string, perm string) (bool, error) {
	row := s.DB.QueryRow(
		`SELECT COUNT(*) FROM user_permissions
		WHERE user_id = (SELECT id FROM users WHERE username = $1)
		AND permission = (SELECT id FROM permissions WHERE name = $2)`,
		name,
		perm,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *EggchanService) DeleteThread(board string, thread int) (int64, error) {
	result, err := s.DB.Exec(
		`DELETE FROM threads
		WHERE board_id = (SELECT id FROM boards WHERE name = $1)
		AND post_num = $2`,
		board,
		thread,
	)

	if err != nil {
		return 0, err
	}

	count, _ := result.RowsAffected()
	return count, nil
}

func (s *EggchanService) DeleteComment(board string, comment int) (int64, error) {
	result, err := s.DB.Exec(
		`DELETE FROM comments
		USING threads
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		AND comments.post_num = $2`,
		board,
		comment,
	)

	if err != nil {
		return 0, err
	}

	count, _ := result.RowsAffected()
	return count, nil
}

func (s *EggchanService) AddUser(user, password string) error {
	_, err := s.DB.Exec(
		`INSERT INTO users (username, password) VALUES ($1, $2)`,
		user,
		password,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *EggchanService) DeleteUser(user string) error {
	result, err := s.DB.Exec(
		`DELETE FROM users WHERE username = $1`,
		user,
	)

	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected != 1 {
		return errors.New("User not found")
	}

	return nil
}

func (s *EggchanService) ListUsers() ([]eggchan.User, error) {
	userList := []eggchan.User{}

	rows, err := s.DB.Query(`SELECT username FROM users ORDER BY id ASC`)
	if err != nil {
		return userList, err
	}

	for i := 0; rows.Next(); i++ {
		var u string
		if err := rows.Scan(&u); err != nil {
			return userList, err
		}
		userList = append(userList, eggchan.User{u, []string{}})
	}

	for i, user := range userList {
		rows, err = s.DB.Query(
			`SELECT name FROM permissions p
			INNER JOIN user_permissions up ON p.id = up.permission
			INNER JOIN users u ON u.id = up.user_id
			WHERE u.username = $1
			ORDER BY p.id ASC`,
			user.Name,
		)
		if err != nil {
			return userList, err
		}

		permissions := []string{}
		for rows.Next() {
			var p string
			if err := rows.Scan(&p); err != nil {
				return userList, err
			}
			permissions = append(permissions, p)
		}

		userList[i].Perms = permissions
	}

	return userList, nil
}

func (s *EggchanService) GrantPermissions(user string, perms []eggchan.Permission) error {
	for _, perm := range perms {
		_, err := s.DB.Exec(
			`INSERT INTO user_permissions (user_id, permission) VALUES
			((SELECT id FROM users WHERE username = $1), (SELECT id FROM permissions WHERE name = $2))`,
			user,
			perm,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *EggchanService) RevokePermissions(user string, perms []eggchan.Permission) error {
	for _, perm := range perms {
		_, err := s.DB.Exec(
			`DELETE FROM user_permissions
			WHERE user_id = (SELECT id FROM users WHERE username = $1)
			AND permission = (SELECT id FROM permissions WHERE name = $2)`,
			user,
			perm,
		)

		// TODO: figure out a better way to do this
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *EggchanService) ListPermissions() ([]eggchan.Permission, error) {
	perms := []eggchan.Permission{}
	rows, err := s.DB.Query(`SELECT name FROM permissions ORDER BY id ASC`)
	if err != nil {
		return perms, err
	}

	for i := 0; rows.Next(); i++ {
		var n string
		if err := rows.Scan(&n); err != nil {
			return perms, err
		}
		perms = append(perms, eggchan.Permission{n})
	}

	return perms, nil
}

func (s *EggchanService) AddBoard(board, description, category string) error {
	var d sql.NullString
	if description == "" {
		d = sql.NullString{"", false}
	} else {
		d = sql.NullString{description, true}
	}

	var c sql.NullString
	if category == "" {
		c = sql.NullString{"", false}
	} else {
		c = sql.NullString{category, true}
	}

	_, err := s.DB.Exec(
		`INSERT INTO boards (name, description, category) VALUES ($1, $2, (SELECT id FROM categories WHERE name = $3))`,
		board,
		d,
		c,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *EggchanService) ValidatePassword(user string, password []byte) (bool, error) {
	pw_row := s.DB.QueryRow(`SELECT password FROM users WHERE username = $1`, user)
	var db_pw []byte
	if err := pw_row.Scan(&db_pw); err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword(db_pw, password); err != nil {
		return false, nil
	} else {
		return true, nil
	}
}

func (s *EggchanService) CheckPermission(user, permission string) (bool, error) {
	row := s.DB.QueryRow(
		`SELECT COUNT(*) from user_permissions
		WHERE user_id = (SELECT id FROM users WHERE username = $1 LIMIT 1)
		AND permission = (SELECT id FROM permissions WHERE name = $2 LIMIT 1)`,
		user,
		permission,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *EggchanService) ShowBoardDesc(board string) (eggchan.Board, error) {
	b_row := s.DB.QueryRow(
		`SELECT boards.name, boards.description, boards.category
		FROM boards
		WHERE boards.name = $1`,
		board,
	)

	var b eggchan.Board
	if err := b_row.Scan(&b.Name, &b.Description, &b.Category); err != nil {
		return b, err
	}

	return b, nil
}

func (s *EggchanService) ShowThreadOP(board string, id int) (eggchan.Thread, error) {
	t_row := s.DB.QueryRow(
		`SELECT
			boards.name,
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
		INNER JOIN boards ON threads.board_id = boards.id
		WHERE threads.board_id = (SELECT id FROM boards WHERE name = $1)
		AND threads.post_num = $2
		GROUP BY boards.name, threads.id
		ORDER BY sort_latest_reply DESC`,
		board,
		id,
	)

	var t eggchan.Thread
	if err := t_row.Scan(&t.Board, &t.PostNum, &t.Subject, &t.Author, &t.Time, &t.NumReplies, &t.SortLatestReply, &t.Comment); err != nil {
		return t, err
	}

	return t, nil
}
