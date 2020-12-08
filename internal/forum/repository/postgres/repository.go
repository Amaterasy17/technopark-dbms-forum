package postgres

import (
	"github.com/jackc/pgx"
	domain "technopark-dbms-forum/internal/forum"
)

type postgresForumRepository struct {
	Conn *pgx.Conn
}

func NewPostgresForumRepository(Conn *pgx.Conn) domain.ForumRepository {
	return &postgresForumRepository{Conn: Conn}
}

func (p *postgresForumRepository) InsertForum() error {
	return nil
}