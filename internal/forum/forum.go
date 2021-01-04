package forum

import (
	"technopark-dbms-forum/models"
)

type ForumUseCase interface {
	Forum() error
}

type ForumRepository interface {
	InsertForum(forum models.Forum) error
}
