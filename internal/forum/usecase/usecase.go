package usecase

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"strconv"
	"strings"
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
	user, err := f.forumRepo.SelectUser(forum.User)
	if err != nil {
		return models.Forum{}, err
	}

	forum.User = user.Nickname
	err = f.forumRepo.InsertForum(forum)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
			return models.Forum{}, models.ErrNotFound
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			forumModel, _ := f.forumRepo.SelectForum(forum.Slug)
			return forumModel, models.ErrConflict
		}
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

	userModel, err = f.forumRepo.SelectUser(user.Nickname)
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
	forum, err := f.forumRepo.SelectForum(thread.Forum)
	if err != nil {
		return models.Thread{}, err
	}

	user, err := f.forumRepo.SelectUser(thread.Author)
	if err != nil {
		return models.Thread{}, err
	}
	thread.Author = user.Nickname
	thread.Forum = forum.Slug

	if (thread.Slug != "") {
		//threadModel, err := f.forumRepo.SelectThreadBySlug(thread.Slug)
		//if err == nil {
		//	return threadModel, models.ErrConflict
		//}
	} else {
		slug, err := uuid.NewRandom()
		if err != nil {
			return models.Thread{}, err
		}
		thread.Slug = slug.String()
	}

	slug := thread.Slug
	thread, err = f.forumRepo.InsertThread(thread)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			threadModel, _ := f.forumRepo.SelectThreadBySlug(slug)
			return threadModel, models.ErrConflict
		}


		return models.Thread{}, err
	}

	return thread, nil
}

func (f *ForumUsecase) CreatePosts(posts []models.Post, thread models.Thread) ([]models.Post, error) {
	tx, err := f.forumRepo.NewTransaction()
	if err != nil {
		return nil, err
	}
	var postsCreated []models.Post

	postsCreated, err = f.forumRepo.InsertPosts(posts, thread)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
						fmt.Println(err)
						tx.Rollback()
						return nil, models.ErrNotFound
		} else {
			fmt.Println(err)
			tx.Rollback()
			return nil, models.ErrConflict
		}
	}
	//created := time.Now()
	//
	//var postsCreated []models.Post
	//for _, post := range posts {
	//	post.Thread = thread.Id
	//	post.Forum = thread.Forum
	//	post.Created = created
	//
	//	post, err := f.forumRepo.InsertPost(post)
	//	if err != nil {
	//		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
	//			fmt.Println(err)
	//			tx.Rollback()
	//			return nil, models.ErrNotFound
	//		} else {
	//			fmt.Println(err)
	//			tx.Rollback()
	//			return nil, models.ErrConflict
	//		}
	//
	//	}
	//
	//	postsCreated = append(postsCreated, post)
	//}
	//
	tx.Commit()
	return postsCreated, nil
}

func (f *ForumUsecase) ThreadDetails(slug string) (models.Thread, error) {
	id, err := strconv.Atoi(slug)
	var thread models.Thread
	if err != nil {
		thread, err = f.forumRepo.SelectThreadBySlug(slug)
		if err != nil {
			return models.Thread{}, err
		}
	} else {
		thread, err = f.forumRepo.SelectThreadById(id)
		if err != nil {
			return models.Thread{}, err
		}
	}

	return thread, nil
}

func (f* ForumUsecase) ListThreads(slug string, params models.Parameters) ([]models.Thread, error) {

	_, err := f.forumRepo.SelectForum(slug)
	if err != nil {
		return nil, err
	}

	threads, err := f.forumRepo.SelectThreads(slug, params)
	if err != nil {
		return nil, err
	}

	return threads, nil

}

func (f* ForumUsecase) StatusDB() models.Status {
	return f.forumRepo.StatusOfForum()
}

func (f *ForumUsecase) ClearDB() error {
	return f.forumRepo.ClearDB()
}

func (f *ForumUsecase) MakeVote(vote models.Vote, thread models.Thread) (models.Thread, error) {
	//voteResult, err := f.forumRepo.SelectVote(vote)
	fmt.Println("WORK PLEASE MRAZ")
	fmt.Println(vote)
	fmt.Println(thread)
	err := f.forumRepo.InsertVote(vote)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			_, err := f.forumRepo.UpdateVote(vote)
			if err != nil {
				return models.Thread{}, err
			}
			return thread, nil
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
			return models.Thread{}, models.ErrNotFound
		}
		return models.Thread{}, err
	}
	return thread, nil

	//if err != nil {
	//	err = f.forumRepo.InsertVote(vote)
	//	if err != nil {
	//		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
	//			return models.Thread{}, models.ErrNotFound
	//		}
	//		return models.Thread{}, err
	//	}
	//	if vote.Voice == 1 {
	//		thread.Votes++
	//	} else {
	//		thread.Votes--
	//	}
	//	return thread, nil
	//}

	//if vote.Voice == voteResult.Voice {
	//	return thread, nil
	//} else {
	//	voteResult, err = f.forumRepo.UpdateVote(vote)
	//	if err != nil {
	//		return models.Thread{}, err
	//	}
	//	if vote.Voice == 1 {
	//		thread.Votes += 2
	//	} else {
	//		thread.Votes -= 2
	//	}
	//	return thread, nil
	//}

}

func (f *ForumUsecase) SumVotesInThread(id int) int {
	return f.forumRepo.SumVotesInThread(id)
}

func (f *ForumUsecase) UpdateMessagePost(update models.PostUpdate) (models.Post, error){
	var post models.Post
	//post, err := f.forumRepo.SelectPost(update.ID)
	//if err != nil {
	//	return models.Post{}, err
	//}

	post.ID = update.ID
	post, err := f.forumRepo.UpdatePost(post, update)
	if err != nil {
		return models.Post{}, models.ErrNotFound
	}

	return post, nil
}


func (f *ForumUsecase) PostFullDetails(id int, related string) (models.PostFull, error) {
	var postFull models.PostFull
	post, err := f.forumRepo.SelectPost(id)
	if err != nil {
		return models.PostFull{}, err
	}
	postFull.Post = &post

	if strings.Contains(related, "user") {
		author, err := f.forumRepo.SelectUser(post.Author)
		if err != nil {
			return models.PostFull{}, err
		}
		postFull.Author = &author
	}

	if strings.Contains(related, "thread") {
		thread, err := f.forumRepo.SelectThreadById(post.Thread)
		if err != nil {
			return models.PostFull{}, err
		}

		if models.IsUuid(thread.Slug) {
			result := models.ThreadToThreadOut(thread)
			postFull.Thread = result
		} else {
			postFull.Thread = thread
		}
	}

	if strings.Contains(related, "forum") {
		forum, err := f.forumRepo.SelectForum(post.Forum)
		if err != nil {
			return models.PostFull{}, err
		}
		postFull.Forum = &forum
	}

	return postFull, nil
}


func (f *ForumUsecase) GetUsersByForum(slug string, params models.Parameters) ([]models.User, error) {
	_, err := f.forumRepo.SelectForum(slug)
	if err != nil {
		return nil, err
	}

	users, err := f.forumRepo.SelectUsersByForum(slug, params)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (f *ForumUsecase) GetPostsOfThread(threadId int, parameters models.Parameters, sort string) ([]models.Post, error) {
	switch sort {
	case "flat":
		return f.forumRepo.PostFlatSort(threadId, parameters)
	case "tree":
		return f.forumRepo.PostTreeSort(threadId, parameters)
	case "parent_tree":
		return f.forumRepo.PostParentTreeSort(threadId, parameters)
	default:
		return f.forumRepo.PostFlatSort(threadId, parameters)
	}
}

func (f *ForumUsecase) UpdateThread(thread models.Thread) (models.Thread, error) {
	var oldThread models.Thread
	var err error
	if thread.Slug == "" {
		oldThread, err = f.forumRepo.SelectThreadById(thread.Id)
		if err != nil {
			return models.Thread{}, err
		}
	} else {
		oldThread, err = f.forumRepo.SelectThreadBySlug(thread.Slug)
		if err != nil {
			return models.Thread{}, err
		}
	}

	if thread.Title == "" {
		thread.Title = oldThread.Title
	}
	if thread.Message == "" {
		thread.Message = oldThread.Message
	}

	return f.forumRepo.UpdateThread(thread)
}