package postgres

import (
	"database/sql"
	"errors"

	"github.com/sector-f/eggchan"
	"golang.org/x/crypto/bcrypt"
)

func (s *EggchanService) DeleteThread(board string, thread int) (int64, error) {
	result, err := s.DB.Exec(
		`DELETE FROM threads
		WHERE board_id = (SELECT id FROM boards WHERE name = $1)
		AND post_num = $2`,
		board,
		thread,
	)

	if err != nil {
		return 0, eggchan.DatabaseError{}
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
		return 0, eggchan.DatabaseError{}
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
		return eggchan.DatabaseError{}
	}

	return nil
}

func (s *EggchanService) DeleteUser(user string) error {
	result, err := s.DB.Exec(
		`DELETE FROM users WHERE username = $1`,
		user,
	)

	if err != nil {
		return eggchan.DatabaseError{}
	}

	affected, _ := result.RowsAffected()
	if affected != 1 {
		return errors.New("User not found")
	}

	return nil
}

func (s *EggchanService) ListUsers() ([]eggchan.User, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, eggchan.DatabaseError{}
	}
	defer tx.Commit()

	userList := []eggchan.User{}
	rows, err := tx.Query(`SELECT username FROM users ORDER BY id ASC`)
	if err != nil {
		return userList, eggchan.DatabaseError{}
	}

	for i := 0; rows.Next(); i++ {
		var u string
		if err := rows.Scan(&u); err != nil {
			return userList, eggchan.DatabaseError{}
		}
		userList = append(userList, eggchan.User{u, []string{}})
	}
	rows.Close()

	for i, user := range userList {
		rows, err = tx.Query(
			`SELECT name FROM permissions p
			INNER JOIN user_permissions up ON p.id = up.permission
			INNER JOIN users u ON u.id = up.user_id
			WHERE u.username = $1
			ORDER BY p.id ASC`,
			user.Name,
		)
		if err != nil {
			return userList, eggchan.DatabaseError{}
		}

		permissions := []string{}
		for rows.Next() {
			var p string
			if err := rows.Scan(&p); err != nil {
				return userList, eggchan.DatabaseError{}
			}
			permissions = append(permissions, p)
		}
		rows.Close()

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
			return eggchan.DatabaseError{}
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
			return eggchan.DatabaseError{}
		}
	}

	return nil
}

func (s *EggchanService) ListPermissions() ([]eggchan.Permission, error) {
	perms := []eggchan.Permission{}
	rows, err := s.DB.Query(`SELECT name FROM permissions ORDER BY id ASC`)
	if err != nil {
		return perms, eggchan.DatabaseError{}
	}

	for i := 0; rows.Next(); i++ {
		var n string
		if err := rows.Scan(&n); err != nil {
			return perms, eggchan.DatabaseError{}
		}
		perms = append(perms, eggchan.Permission{n})
	}

	return perms, nil
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
		return eggchan.DatabaseError{}
	}

	return nil
}

func (s *EggchanService) ValidatePassword(user string, password []byte) (bool, error) {
	pw_row := s.DB.QueryRow(`SELECT password FROM users WHERE username = $1`, user)
	var db_pw []byte

	err := pw_row.Scan(&db_pw)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return false, eggchan.UserNotFoundError{}
	default:
		return false, eggchan.DatabaseError{}
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

	err := row.Scan(&count)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return false, eggchan.UserNotFoundError{}
	default:
		return false, eggchan.DatabaseError{}
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
