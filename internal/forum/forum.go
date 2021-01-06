package forum

import (
	"technopark-dbms-forum/models"
)

type ForumUseCase interface {
	Forum(forum models.Forum) (models.Forum, error)
	CreateUser(user models.User) ([]models.User, error)
	GetUser(nickname string) (models.User, error)
	ChangeUserProfile(user models.User) (models.User, error)
	ForumDetails(slug string) (models.Forum, error)
	CreatingThread(thread models.Thread) (models.Thread, error)
	CreatePosts(posts []models.Post, slug string) ([]models.Post, error)
	ThreadDetails(slug string) (models.Thread, error)
	StatusDB() models.Status
	ClearDB() error
}

type ForumRepository interface {
	InsertForum(forum models.Forum) error
	CheckForum(forum models.Forum) (models.Forum, bool)
	SelectUsers(user models.User) ([]models.User, error)
	InsertUser(user models.User) error
	SelectUser(user string) (models.User, error)
	SelectUserByEmail(user models.User) (models.User, error)
	UpdateUserInfo(user models.User) (models.User, error)
	SelectForum(forumName string) (models.Forum, error)
	SelectThreadBySlug(slug string) (models.Thread, error)
	InsertThread(thread models.Thread) error
	SelectThreadById(id int) (models.Thread, error)
	CheckParent(post models.Post) bool
	InsertPost(post models.Post) (models.Post, error)
	StatusOfForum() models.Status
	ClearDB() error
}
