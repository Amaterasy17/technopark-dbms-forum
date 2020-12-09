package postgres

import (
	"github.com/jackc/pgx"
	domain "technopark-dbms-forum/internal/forum"
)

type postgresForumRepository struct {
	Conn *pgx.ConnPool
}

func NewPostgresForumRepository(Conn *pgx.ConnPool) domain.ForumRepository {
	return &postgresForumRepository{Conn: Conn}
}

func (p *postgresForumRepository) InsertForum() error {
	return nil
}