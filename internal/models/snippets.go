package models

import (
	"database/sql"
	"errors"
	"time"
)

/* defining a snippet type to hold the data for individual snippet */
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

/* Defining a SnippetModel type which wraps a sql.DB connection pool. */
type SnippetModel struct {
	DB *sql.DB
}

// it will insert a new snippet
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title,content,created,expires)
			 VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}
	//using lastInsertID() method to get result for our latest id
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// This will return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	//a pointer to a zeroed snippet struct
	s := &Snippet{}
	stmt := `SELECT id,title,content,created,expires FROM snippets
		   WHERE expires > UTC_TIMESTAMP() AND id = ?`

	/*//this returns a pointer to sql.Row objext which holds the result from the database
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	*/

	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// this will return 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id,title,content,created,expires FROM snippets
			WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	//to close the connection pool so we dont get short in connection pool before it exits latest() for safety
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	/* when the rows.next() loop has finished we call rows.Err() to retrieve any error that was 
	encounterd during the iteration. so its imp to call this not all succesfull iterations means 
	a succesfull result
	 */

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
