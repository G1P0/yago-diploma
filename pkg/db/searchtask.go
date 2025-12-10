package db

import (
	"database/sql"
	"strconv"
)

func SearchTasks(limit int, search string) ([]*Task, error) {
	like := "%" + search + "%"

	rows, err := db.Query(
		`SELECT id, date, title, comment, repeat
         FROM scheduler
         WHERE title LIKE ? OR comment LIKE ?
         ORDER BY date
         LIMIT ?`,
		like, like, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func TasksByDate(limit int, date string) ([]*Task, error) {
	rows, err := db.Query(
		`SELECT id, date, title, comment, repeat
         FROM scheduler
         WHERE date = ?
         ORDER BY date
         LIMIT ?`,
		date, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task

	for rows.Next() {
		var (
			id              int64
			date, title     string
			comment, repeat string
		)

		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, &Task{
			ID:      strconv.FormatInt(id, 10),
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}
