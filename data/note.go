package data

import (
	"fmt"
	"time"
)

type Note struct {
	Id          int
	Uuid        string
	UserId      int
	Date        string
	Yesterday   string
	Today       string
	Impediments string
	GoBacks     string
}

// Get all the notes in the database
func (user *User) Notes() (notes []Note, err error) {
	rows, err := Db.Query("SELECT id, uuid, user_id, date, yesterday, today, impediments, go_backs FROM notes WHERE user_id = $1 ORDER BY date DESC", user.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		note := Note{}
		if err = rows.Scan(&note.Id, &note.Uuid, &note.UserId, &note.Date, &note.Yesterday, &note.Today, &note.Impediments, &note.GoBacks); err != nil {
			return
		}
		notes = append(notes, note)
	}
	rows.Close()
	return
}

// Get a note by the UUID
func NoteByUUID(uuid string) (note Note, err error) {
	note = Note{}
	err = Db.QueryRow("SELECT id, uuid, user_id, date, yesterday, today, impediments, go_backs FROM notes WHERE uuid = $1", uuid).
		Scan(&note.Id, &note.Uuid, &note.UserId, &note.Date, &note.Yesterday, &note.Today, &note.Impediments, &note.GoBacks)
	return
}

// Get a note by the date
func NoteByDate(date string) (note Note, err error) {
	note = Note{}
	err = Db.QueryRow("SELECT id, uuid, user_id, date, yesterday, today, impediments, go_backs FROM notes WHERE date = $1", date).
		Scan(&note.Id, &note.Uuid, &note.UserId, &note.Date, &note.Yesterday, &note.Today, &note.Impediments, &note.GoBacks)
	return
}

// Format the date for display
func (note *Note) DisplayDate() string {
	t, err := time.Parse(time.RFC3339, note.Date)
	if err != nil {
		return note.Date
	}
	return fmt.Sprintf("%d.%02d.%02d", t.Year(), t.Month(), t.Day())
}

// Get the user who made the note
func (note *Note) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", note.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

func (user *User) CreateNote(date string, yesterday string, today string, gobacks string, impediments string) (note Note, err error) {
	statement := "INSERT INTO notes (uuid, user_id, date, yesterday, today, impediments, go_backs) values ($1, $2, $3, $4, $5, $6, $7) returning uuid, user_id, date, yesterday, today, impediments, go_backs"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(createUUID(), user.Id, date, yesterday, today, impediments, gobacks).Scan(&note.Uuid, &note.Id, &note.Date, &note.Yesterday, &note.Today, &note.Impediments, &note.GoBacks)
	return
}

func (user *User) UpdateNote(uuid string, date string, yesterday string, today string, gobacks string, impediments string) (note Note, err error) {
	statement := "UPDATE notes SET date = $1, yesterday = $2, today = $3, impediments = $4, go_backs = $5 WHERE uuid = $6 returning uuid, date, yesterday, today, impediments, go_backs"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(date, yesterday, today, impediments, gobacks, uuid).Scan(&note.Date, &note.Yesterday, &note.Today, &note.Impediments, &note.GoBacks, &note.Uuid)
	return
}

func (user *User) DeleteNote(uuid string) (err error) {
	statement := "DELETE FROM notes WHERE uuid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(statement, uuid)
	if err != nil {
		return
	}
	fmt.Printf("%v", res)
	return
}
