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