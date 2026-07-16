package db

import (
	"database/sql"
	"fmt"
	"strings"
)

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task Task) (int64, error) {
	res, err := DB.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id int64) (*Task, error) {
	var task Task
	row := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task Task) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteTask(id int64) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func ListTasks(search string) ([]Task, error) {
	tasks := make([]Task, 0)
	rows, err := DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if strings.TrimSpace(search) == "" {
		return tasks, nil
	}

	filtered := make([]Task, 0, len(tasks))
	lower := strings.ToLower(search)
	for _, task := range tasks {
		formattedDate := ""
		if len(task.Date) == 8 {
			formattedDate = fmt.Sprintf("%s.%s.%s", task.Date[6:8], task.Date[4:6], task.Date[0:4])
		}
		if strings.Contains(strings.ToLower(task.Title), lower) ||
			strings.Contains(strings.ToLower(task.Comment), lower) ||
			strings.Contains(strings.ToLower(formattedDate), lower) {
			filtered = append(filtered, task)
		}
	}
	return filtered, nil
}
