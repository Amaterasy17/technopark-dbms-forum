package postgres

import (
	"fmt"
	"github.com/jackc/pgx"
	domain "technopark-dbms-forum/internal/forum"
	models "technopark-dbms-forum/models"
)

type postgresForumRepository struct {
	Conn *pgx.ConnPool
}

func NewPostgresForumRepository(Conn *pgx.ConnPool) domain.ForumRepository {
	return &postgresForumRepository{Conn: Conn}
}

func (p *postgresForumRepository) InsertForum(forum models.Forum) error {
	_, err := p.Conn.Exec(	`Insert INTO forum(Slug, "user", Title) VALUES ($1, $2, $3);`,
		forum.Slug, forum.User, forum.Title)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgresForumRepository) SelectForum(forumName string) (models.Forum, error) {
	var forum models.Forum
	row := p.Conn.QueryRow(`Select slug, "user", title, posts, threads From forum
				Where slug=$1`, forumName)
	err := row.Scan(&forum.Slug, &forum.User, &forum.Title, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, models.ErrNotFound
	}
	return forum, nil
}

func (p *postgresForumRepository) CheckForum(forum models.Forum) (models.Forum, bool) {
	resultForum := models.Forum{
		Posts: -1,
	}
	row := p.Conn.QueryRow(`Select slug, user, title, posts, threads From forum
				Where slug=$1`, forum.Slug)
	_ = row.Scan(&resultForum.Slug, &resultForum.User, &resultForum.Title, &resultForum.Posts, &resultForum.Threads)
	if resultForum.Posts == -1 {
		return models.Forum{},false
	}
	return resultForum, true
}

func (p *postgresForumRepository) SelectUsers(user models.User) ([]models.User, error) {
	var users []models.User
	rows, err := p.Conn.Query(`Select Nickname, FullName, About, Email From users Where Nickname=$1 or Email=$2;`,
														user.Nickname, user.Email)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return users, err
	}
	for rows.Next() {
		var userModel models.User
		err := rows.Scan(&userModel.Nickname, &userModel.FullName, &userModel.About, &userModel.Email)
		if err != nil {
			return users, err
		}
		users = append(users, userModel)
	}

	return users, nil
}

func (p *postgresForumRepository) InsertUser(user models.User) error {
	_, err := p.Conn.Exec(	`Insert INTO users(Nickname, FullName, About, Email) VALUES ($1, $2, $3, $4);`,
		user.Nickname, user.FullName, user.About, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgresForumRepository) SelectUser(user string) (models.User, error) {
	var userModel models.User
	row := p.Conn.QueryRow(`Select nickname, fullname, about, email From users Where nickname=$1;`, user)
	err := row.Scan(&userModel.Nickname, &userModel.FullName, &userModel.About, &userModel.Email)
	if err != nil {
		return models.User{}, models.ErrNotFound
	}
	return userModel, nil
}

func (p *postgresForumRepository) SelectUserByEmail(user models.User) (models.User, error) {
	var userModel models.User
	row := p.Conn.QueryRow(`Select nickname, email from users Where email=$1;`, user.Email)
	err := row.Scan(&userModel.Nickname, &userModel.Email)
	if err != nil {
		return models.User{}, nil
	}
	if userModel.Nickname == user.Nickname {
		return models.User{}, nil
	}
	return userModel, models.ErrConflict
}

func (p *postgresForumRepository) UpdateUserInfo(user models.User) (models.User, error) {
	_, err := p.Conn.Exec(`UPDATE users SET fullname=$1 WHERE nickname=$2;`, user.FullName, user.Nickname)
	if err != nil {
		return models.User{}, err
	}

	_, err = p.Conn.Exec(`UPDATE users SET about=$1 WHERE nickname=$2;`, user.About, user.Nickname)
	if err != nil {
		return models.User{}, err
	}

	_, err = p.Conn.Exec(`UPDATE users SET email=$1 WHERE nickname=$2;`, user.Email, user.Nickname)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (p *postgresForumRepository) SelectThreadBySlug(slug string) (models.Thread, error) {
	var thread models.Thread
	row := p.Conn.QueryRow(`Select id, title, author, forum, message, votes, slug, created from thread
							Where slug=$1;`, slug)
	err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
					&thread.Slug, &thread.Created)
	if err != nil {
		return models.Thread{}, models.ErrNotFound
	}
	return thread, nil
}

func (p *postgresForumRepository) InsertThread(thread models.Thread) error {
	_, err := p.Conn.Exec(	`Insert INTO thread(Title, Author, Forum, Message, slug, Votes)
							VALUES ($1, $2, $3, $4, $5, $6);`, thread.Title, thread.Author, thread.Forum,
							thread.Message, thread.Slug, thread.Votes)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgresForumRepository) SelectThreadById(id int) (models.Thread, error) {
	var thread models.Thread
	row := p.Conn.QueryRow(`Select id, title, author, forum, message, votes, slug, created from thread
							Where id=$1;`, id)
	err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
		&thread.Slug, &thread.Created)
	if err != nil {
		return models.Thread{}, models.ErrNotFound
	}
	return thread, nil
}

func (p *postgresForumRepository) CheckParent(post models.Post) bool {
	var id int
	row := p.Conn.QueryRow(`Select id from post where id=$1;`, post.Parent)
	err := row.Scan(&id)
	if err != nil {
		return false
	}
	return true
}
func (p *postgresForumRepository) InsertPost(post models.Post) (models.Post, error) {
	row := p.Conn.QueryRow(`INSERT INTO post(author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`,
		post.Author, post.Created, post.Forum, post.Message, post.Parent, post.Thread)
	var postModel models.Post
	err := row.Scan(&postModel.ID, &postModel.Author, &postModel.Created, &postModel.Forum,  &postModel.IsEdited,
		&postModel.Message, &postModel.Parent, &postModel.Thread)
	return postModel, err
}