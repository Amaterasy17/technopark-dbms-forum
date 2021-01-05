package forum

import (
	"technopark-dbms-forum/models"
)

type ForumUseCase interface {
	Forum(forum models.Forum) (models.Forum, error)
	CreateUser(user models.User) ([]models.User, error)
	GetUser(nickname string) (models.User, error)
}

type ForumRepository interface {
	InsertForum(forum models.Forum) error
	CheckForum(forum models.Forum) (models.Forum, bool)
	SelectUsers(user models.User) ([]models.User, error)
	InsertUser(user models.User) error
	SelectUser(user string) (models.User, error)
}
