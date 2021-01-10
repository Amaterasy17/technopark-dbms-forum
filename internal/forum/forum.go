package forum

import (
	"github.com/jackc/pgx"
	"technopark-dbms-forum/models"
)

type ForumUseCase interface {
	Forum(forum models.Forum) (models.Forum, error)
	CreateUser(user models.User) ([]models.User, error)
	GetUser(nickname string) (models.User, error)
	ChangeUserProfile(user models.User) (models.User, error)
	ForumDetails(slug string) (models.Forum, error)
	CreatingThread(thread models.Thread) (models.Thread, error)
	CreatePosts(posts []models.Post, thread models.Thread) ([]models.Post, error)
	ThreadDetails(slug string) (models.Thread, error)
	StatusDB() models.Status
	ClearDB() error
	MakeVote(vote models.Vote, thread models.Thread) (models.Thread, error)
	SumVotesInThread(id int) int
	UpdateMessagePost(update models.PostUpdate) (models.Post, error)
	PostFullDetails(id int, related string) (models.PostFull, error)
	ListThreads(slug string, params models.Parameters) ([]models.Thread, error)
	GetUsersByForum(slug string, params models.Parameters) ([]models.User, error)
	GetPostsOfThread(threadId int, parameters models.Parameters, sort string) ([]models.Post, error)
	UpdateThread(thread models.Thread) (models.Thread, error)
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
	InsertThread(thread models.Thread) (models.Thread,error)
	SelectThreadById(id int) (models.Thread, error)
	CheckParent(post models.Post) bool
	InsertPost(post models.Post) (models.Post, error)
	StatusOfForum() models.Status
	ClearDB() error
	SelectVote(vote models.Vote) (models.Vote, error)
	UpdateVote(vote models.Vote) (models.Vote, error)
	InsertVote(vote models.Vote)  error
	SumVotesInThread(id int) int
	SelectPost(id int) (models.Post, error)
	UpdatePost(post models.Post, postUpdate models.PostUpdate) (models.Post, error)
	SelectThreads(slug string, params models.Parameters) ([]models.Thread, error)
	SelectUsersByForum(slug string, params models.Parameters) ([]models.User, error)
	PostParentTreeSort(threadId int, parameters models.Parameters) ([]models.Post, error)
	PostTreeSort(threadId int, parameters models.Parameters) ([]models.Post, error)
	PostFlatSort(id int, parameters models.Parameters) ([]models.Post, error)
	UpdateThread(thread models.Thread) (models.Thread, error)
	NewTransaction() (*pgx.Tx, error)
	Rollback(tx *pgx.Tx)
}
