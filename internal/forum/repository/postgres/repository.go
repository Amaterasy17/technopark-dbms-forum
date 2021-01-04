package postgres

import (
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