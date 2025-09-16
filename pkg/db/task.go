package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func isDateSearch(s string) bool {
	_, err := time.Parse("02.01.2006", s)
	return err == nil
}

func convertToDBDateFormat(dateString string) (string, error) {
	t, err := time.Parse("02.01.2006", dateString)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}
	return t.Format("20060102"), nil
}

func Tasks(limit int, search string) ([]*Task, error) {

	var rows *sql.Rows
	var err error

	if search == "" {

		query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
		rows, err = db.Query(query, limit)
	} else if isDateSearch(search) {
		dbDate, err := convertToDBDateFormat(search)
		if err != nil {
			return nil, err
		}

		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = db.Query(query, dbDate, limit)
		if err != nil {
			return nil, err
		}
	} else {
		searchPattern := "%" + search + "%"
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		rows, err = db.Query(query, searchPattern, searchPattern, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, id)
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %v", err)
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	if task.Id == "" {
		return fmt.Errorf("task ID is required")
	}

	if task.Date == "" {
		return fmt.Errorf("date is required")
	}

	if task.Title == "" {
		return fmt.Errorf("title is required")
	}

	if len(task.Date) != 8 {
		return fmt.Errorf("invalid date format")
	}

	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func UpdateDate(id string, newDate string) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	result, err := db.Exec(query, newDate, id)
	if err != nil {
		return fmt.Errorf("failed to update date: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
