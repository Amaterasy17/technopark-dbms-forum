package models

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
	"time"
)

type Forum struct {
	ID      int    `json:"-"`
	UserId  int    `json:"-"`
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}

type ThreadOut struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	AuthorId  int    `json:"-"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int       `json:"votes"`
	Slug    string    `json:"-"`
	Created time.Time `json:"created"`
}

type Thread struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	AuthorId  int    `json:"-"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int       `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}



type ThreadIn struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int       `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"-"`
}


type User struct {
	ID       int    `json:"-"`
	Nickname string `json:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type Post struct {
	ID       int              `json:"id"`
	Author   string           `json:"author"`
	AuthorId int              `json:"-"`
	Created  time.Time        `json:"created"`
	Forum    string           `json:"forum"`
	IsEdited bool             `json:"isEdited"`
	Message  string           `json:"message"`
	Parent   JsonNullInt64    `json:"parent"`
	Thread   int              `json:"thread,"`
	Path     pgtype.Int8Array `json:"-"`
}

type Status struct {
	User   int `json:"user"`
	Forum  int `json:"forum"`
	Thread int `json:"thread"'`
	Post   int `json:"post"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	Thread   int    `json:"-"`
}

type PostUpdate struct {
	ID       int       `json:"-"`
	Message string `json:"message"`
}

type PostFull struct {
	Author *User       `json:"author"`
	Forum  *Forum      `json:"forum"`
	Post   *Post       `json:"post"`
	Thread interface{} `json:"thread"`
}

type Parameters struct {
	Limit int    `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

func IsUuid(value string) bool {
	n := len(value)

	if n > 36 || n < 32 {
		return false
	}

	_, err := uuid.Parse(value)

	return err == nil
}

func ThreadToThreadOut(thread Thread) ThreadOut {
	return ThreadOut{
		Id:      thread.Id,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   thread.Forum,
		Message: thread.Message,
		Votes:   thread.Votes,
		Slug:    thread.Slug,
		Created: thread.Created,
	}
}