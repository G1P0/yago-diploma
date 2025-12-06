package db

import (
	"errors"
	"strconv"
)

func UpdateTask(task *Task) error {
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return errors.New("invalid id")
	}

	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	res, err := db.Exec(query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		id,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("task not found")
	}

	return nil
}
