package forum

type ForumUseCase interface {
	Forum() error
}

type ForumRepository interface {
	InsertForum() error
}
