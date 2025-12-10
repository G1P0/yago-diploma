package db

import (
	"database/sql"
	"errors"
	"strconv"
)

func GetTask(idStr string) (*Task, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	row := db.QueryRow(
		`SELECT id, date, title, comment, repeat
         FROM scheduler
         WHERE id = ?`,
		id,
	)

	var (
		dbID            int64
		date, title     string
		comment, repeat string
	)

	if err := row.Scan(&dbID, &date, &title, &comment, &repeat); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return &Task{
		ID:      strconv.FormatInt(dbID, 10),
		Date:    date,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}, nil
}
