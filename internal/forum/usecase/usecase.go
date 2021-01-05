package usecase

import (
	"fmt"
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

func (f *ForumUsecase) CreateUser(user models.User) ([]models.User, error) {
	var users []models.User
	users, err := f.forumRepo.SelectUsers(user)
	if err != nil {
		fmt.Println(err)
	}
	if len(users) != 0 {
		return users, models.ErrConflict
	}

	err = f.forumRepo.InsertUser(user)
	if err != nil {
		return nil, err
	}
	users = append(users, user)
	return users, nil
}

func (f *ForumUsecase) GetUser(nickname string) (models.User, error) {
	return f.forumRepo.SelectUser(nickname)
}