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
	//_, err := p.Conn.Exec(	`Insert INTO forum()`)
	return nil
}