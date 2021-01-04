package usecase

import (
	domain "technopark-dbms-forum/internal/forum"
	"technopark-dbms-forum/models"
)

type ForumUsecase struct {
	forumRepo domain.ForumRepository
}

func NewForumUsecase(forumRepo domain.ForumRepository) domain.ForumUseCase {
	return &ForumUsecase{forumRepo: forumRepo}
}

func (f *ForumUsecase) Forum(forum models.Forum) (models.Forum, error) {
	checkForum, forumBool := f.forumRepo.CheckForum(forum)
	if !forumBool {
		return checkForum, nil
	}

	//err := f.forumRepo.InsertForum(forum)
	//if err != nil {
	//	return err
	//}

	return models.Forum{}, nil
}