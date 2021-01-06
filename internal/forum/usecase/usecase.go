package usecase

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	domain "technopark-dbms-forum/internal/forum"
	"technopark-dbms-forum/models"
	"time"
)

type ForumUsecase struct {
	forumRepo domain.ForumRepository
}

func NewForumUsecase(forumRepo domain.ForumRepository) domain.ForumUseCase {
	return &ForumUsecase{forumRepo: forumRepo}
}

func (f *ForumUsecase) Forum(forum models.Forum) (models.Forum, error) {
	_, err := f.forumRepo.SelectUser(forum.User)
	if err != nil {
		return models.Forum{}, err
	}

	forumModel, err := f.forumRepo.SelectForum(forum.Slug)
	if err == nil {
		return forumModel, models.ErrConflict
	}

	err = f.forumRepo.InsertForum(forum)
	if err != nil {
		return models.Forum{}, err
	}

	forum.Posts = 0
	forum.Threads = 0

	return forum, nil
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

func (f *ForumUsecase) ChangeUserProfile(user models.User) (models.User, error) {
	_, err := f.forumRepo.SelectUser(user.Nickname)
	if err != nil {
		return models.User{}, err
	}

	_, err = f.forumRepo.SelectUserByEmail(user)
	if err != nil {
		return models.User{}, err
	}

	userModel, err := f.forumRepo.UpdateUserInfo(user)
	if err != nil {
		return models.User{}, err
	}

	return userModel, nil
}

func (f *ForumUsecase) ForumDetails(slug string) (models.Forum, error) {
	forum, err := f.forumRepo.SelectForum(slug)
	if err != nil {
		return models.Forum{}, err
	}
	return forum, nil
}

func (f *ForumUsecase) CreatingThread(thread models.Thread) (models.Thread, error) {
	_, err := f.forumRepo.SelectForum(thread.Forum)
	if err != nil {
		return models.Thread{}, err
	}

	_, err = f.forumRepo.SelectUser(thread.Author)
	if err != nil {
		return models.Thread{}, err
	}

	if (thread.Slug != "") {
		threadModel, err := f.forumRepo.SelectThreadBySlug(thread.Slug)
		if err == nil {
			return threadModel, models.ErrConflict
		}
	} else {
		slug, err := uuid.NewRandom()
		if err != nil {
			return models.Thread{}, err
		}
		thread.Slug = slug.String()
	}


	err = f.forumRepo.InsertThread(thread)
	if err != nil {
		return models.Thread{}, err
	}
	fmt.Println("popal cuda")

	result, err := f.forumRepo.SelectThreadBySlug(thread.Slug)
	if err != nil {
		fmt.Println(err)
		return models.Thread{}, err
	}

	return result, nil
}

func (f *ForumUsecase) CreatePosts(posts []models.Post, slug string) ([]models.Post, error) {
	id, err := strconv.Atoi(slug)
	var thread models.Thread
	if err != nil {
		thread, err = f.forumRepo.SelectThreadBySlug(slug)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = f.forumRepo.SelectThreadById(id)
		if err != nil {
			return nil, err
		}
	}

	for _, post := range posts {
		_, err := f.forumRepo.SelectUser(post.Author)
		if err != nil {
			return nil, err
		}

		if post.Parent != 0 && !f.forumRepo.CheckParent(post) {
			return nil, models.ErrConflict
		}
	}

	created := time.Now()

	var postsCreated []models.Post
	for _, post := range posts {
		post.Thread = thread.Id
		post.Forum = thread.Forum
		post.Created = created

		post, err = f.forumRepo.InsertPost(post)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		postsCreated = append(postsCreated, post)
	}

	return postsCreated, nil
}