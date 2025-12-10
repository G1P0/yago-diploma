package db

import (
	"strconv"
)

func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query(
		`SELECT id, date, title, comment, repeat 
         FROM scheduler 
         ORDER BY date 
         LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
