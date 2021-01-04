package forum

import (
	"technopark-dbms-forum/models"
)

type ForumUseCase interface {
	Forum(forum models.Forum) (models.Forum, error)
}

type ForumRepository interface {
	InsertForum(forum models.Forum) error
	CheckForum(forum models.Forum) (models.Forum, bool)
}
