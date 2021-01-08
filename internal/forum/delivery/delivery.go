package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	domain "technopark-dbms-forum/internal/forum"
	"technopark-dbms-forum/models"
	"time"
)

type ForumHandler struct {
	ForumUseCase domain.ForumUseCase
}

func JSONError(message string) []byte {
	jsonError, err := json.Marshal(models.Error{Message: message})
	if err != nil {
		return []byte("")
	}
	return jsonError
}


func NewForumHandler(r *mux.Router, forumUseCase domain.ForumUseCase) {
	handler := &ForumHandler{ForumUseCase: forumUseCase}

	r.HandleFunc("/api/forum/create", handler.Forum).Methods(http.MethodPost)
	r.HandleFunc("/api/forum/{slug}/create", handler.CreateThread).Methods(http.MethodPost)
	r.HandleFunc("/api/forum/{slug}/details", handler.ForumInfo).Methods(http.MethodGet)

	r.HandleFunc("/api/user/{nickname}/create", handler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/api/user/{nickname}/profile", handler.ProfileUser).Methods(http.MethodGet)
	r.HandleFunc("/api/user/{nickname}/profile", handler.ChangeProfileInformation).Methods(http.MethodPost)


	r.HandleFunc("/api/thread/{slug_or_id}/create", handler.CreatePost).Methods(http.MethodPost)
	r.HandleFunc("/api/thread/{slug_or_id}/details", handler.ThreadDetails).Methods(http.MethodGet)
	r.HandleFunc("/api/thread/{slug_or_id}/posts", handler.PostsOfThread).Methods(http.MethodGet)

	r.HandleFunc("/api/service/status", handler.StatusDB).Methods(http.MethodGet)
	r.HandleFunc("/api/service/clear", handler.ClearDB).Methods(http.MethodPost)

	r.HandleFunc("/api/thread/{slug_or_id}/vote", handler.MakeVote).Methods(http.MethodPost)

	r.HandleFunc("/api/post/{id}/details", handler.PostUpdate).Methods(http.MethodPost)
	r.HandleFunc("/api/post/{id}/details", handler.PostDetails).Methods(http.MethodGet)

	r.HandleFunc("/api/forum/{slug}/threads", handler.ThreadsOfForum).Methods(http.MethodGet)
	r.HandleFunc("/api/forum/{slug}/users", handler.UsersOfForum).Methods(http.MethodGet)
}

func (f *ForumHandler) Forum(w http.ResponseWriter, r *http.Request) {
	forum := models.Forum{}
	err := json.NewDecoder(r.Body).Decode(&forum)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	forum, err = f.ForumUseCase.Forum(forum)
	status := models.GetStatusCodePost(err)
	if status == 409 {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodePost(err))

		body, err := json.Marshal(forum)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		w.Write(body)
		return
	}
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(forum)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodePost(err))
	w.Write(body)
}

func (f *ForumHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/api/forum/")
	slug = strings.TrimSuffix(slug, "/create")
	fmt.Println(slug)

	threadIn := models.ThreadIn{}
	err := json.NewDecoder(r.Body).Decode(&threadIn)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	thread := models.Thread{
		Id:      threadIn.Id,
		Title:    threadIn.Title,
		Author:   threadIn.Author,
		Forum:    threadIn.Forum,
		Message:  threadIn.Message,
		Votes:   0,
		Slug:     threadIn.Slug,
		Created: time.Time{},
	}
	thread.Forum = slug
	fmt.Println("has gone")

	thread, err = f.ForumUseCase.CreatingThread(thread)
 	status := models.GetStatusCodePost(err)
 	if status == 409 {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodePost(err))

		body, err := json.Marshal(thread)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		w.Write(body)
		return
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodePost(err))
	w.Write(body)
}

func (f *ForumHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create user")
	nickname := strings.TrimPrefix(r.URL.Path, "/api/user/")
	nickname = strings.TrimSuffix(nickname, "/create")
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Nickname = nickname

	users, err := f.ForumUseCase.CreateUser(user)
	status := models.GetStatusCodePost(err)
	if status == 409 {
		body, err := json.Marshal(users)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}
		w.WriteHeader(status)
		w.Write(body)
		return
	}
	if status == 201 {
		body, err := json.Marshal(users[0])
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}
		w.WriteHeader(status)
		w.Write(body)
		return
	}
}

func (f *ForumHandler) ProfileUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get profile user")
	nickname := strings.TrimPrefix(r.URL.Path, "/api/user/")
	nickname = strings.TrimSuffix(nickname, "/profile")
	fmt.Println(nickname)

	user, err := f.ForumUseCase.GetUser(nickname)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodeGet(err))
	w.Write(body)
}

func (f *ForumHandler) ChangeProfileInformation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("change profile user")
	nickname := strings.TrimPrefix(r.URL.Path, "/api/user/")
	nickname = strings.TrimSuffix(nickname, "/profile")
	fmt.Println(nickname)

	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Nickname = nickname

	userModel, err := f.ForumUseCase.ChangeUserProfile(user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(userModel)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodeGet(err))
	w.Write(body)
}

func (f *ForumHandler) ForumInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("forumInfo")
	slug := strings.TrimPrefix(r.URL.Path, "/api/forum/")
	slug = strings.TrimSuffix(slug, "/details")
	fmt.Println(slug)

	forum, err := f.ForumUseCase.ForumDetails(slug)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(forum)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodeGet(err))
	w.Write(body)
}

func (f *ForumHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Post")
	slug := strings.TrimPrefix(r.URL.Path, "/api/thread/")
	slug = strings.TrimSuffix(slug, "/create")
	fmt.Println(slug)

	var posts []models.Post
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	posts, err = f.ForumUseCase.CreatePosts(posts, slug)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodePost(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(posts)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodePost(err))
	w.Write(body)
}

func (f *ForumHandler) ThreadDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Println("thread details")
	slug := strings.TrimPrefix(r.URL.Path, "/api/thread/")
	slug = strings.TrimSuffix(slug, "/details")
	fmt.Println(slug)

	thread, err := f.ForumUseCase.ThreadDetails(slug)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodeGet(err))
	w.Write(body)
}

func (f *ForumHandler) StatusDB(w http.ResponseWriter, r *http.Request) {
	status := f.ForumUseCase.StatusDB()
	body, err := json.Marshal(status)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func (f *ForumHandler) ClearDB(w http.ResponseWriter, r *http.Request)  {
	err := f.ForumUseCase.ClearDB()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
}

func (f *ForumHandler) MakeVote(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("Voting")
	slug := strings.TrimPrefix(r.URL.Path, "/api/thread/")
	slug = strings.TrimSuffix(slug, "/vote")
	fmt.Println(slug)

	thread, err := f.ForumUseCase.ThreadDetails(slug)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	var vote models.Vote
	err = json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vote.Thread = thread.Id

	vote, err = f.ForumUseCase.MakeVote(vote)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	thread.Votes = f.ForumUseCase.SumVotesInThread(thread.Id)

	body, err := json.Marshal(thread)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func (f *ForumHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("post update")
	slug := strings.TrimPrefix(r.URL.Path, "/api/post/")
	slug = strings.TrimSuffix(slug, "/details")
	fmt.Println(slug)
	id, err := strconv.Atoi(slug)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		w.Write(JSONError(err.Error()))
		return
	}

	var postUpdate models.PostUpdate
	err = json.NewDecoder(r.Body).Decode(&postUpdate)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postUpdate.ID = id

	post, err := f.ForumUseCase.UpdateMessagePost(postUpdate)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(post)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func (f *ForumHandler) PostDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Println("post details")
	slug := strings.TrimPrefix(r.URL.Path, "/api/post/")
	slug = strings.TrimSuffix(slug, "/details")
	fmt.Println(slug)
	id, err := strconv.Atoi(slug)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		w.Write(JSONError(err.Error()))
		return
	}

	related := r.URL.Query().Get("related")

	postFull, err := f.ForumUseCase.PostFullDetails(id, related)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}


	body, err := json.Marshal(postFull)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func (f *ForumHandler) ThreadsOfForum(w http.ResponseWriter, r *http.Request) {
	fmt.Println("threads of forum")

	w.Header().Set("Content-Type", "application/json")

	var params models.Parameters
	var err error
	params.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		params.Limit = 100
	}

	params.Since = r.URL.Query().Get("since")

	params.Desc, err = strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		params.Desc = false
	}

	slug := strings.TrimPrefix(r.URL.Path, "/api/forum/")
	slug = strings.TrimSuffix(slug, "/threads")
	fmt.Println(slug)


	threads, err := f.ForumUseCase.ListThreads(slug, params)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(threads)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func (f *ForumHandler) UsersOfForum(w http.ResponseWriter, r *http.Request) {
	fmt.Println("users of forum")

	w.Header().Set("Content-Type", "application/json")

	var params models.Parameters
	var err error
	params.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		params.Limit = 100
	}

	params.Since = r.URL.Query().Get("since")

	params.Desc, err = strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		params.Desc = false
	}

	slug := strings.TrimPrefix(r.URL.Path, "/api/forum/")
	slug = strings.TrimSuffix(slug, "/users")
	fmt.Println(slug)

	users, err := f.ForumUseCase.GetUsersByForum(slug, params)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}


func (f *ForumHandler) PostsOfThread(w http.ResponseWriter, r *http.Request) {
	fmt.Println("posts of forum")

	w.Header().Set("Content-Type", "application/json")

	var params models.Parameters
	var err error
	params.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		params.Limit = 100
	}

	params.Since = r.URL.Query().Get("since")
	sort := r.URL.Query().Get("sort")

	params.Desc, err = strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		params.Desc = false
	}

	slug := strings.TrimPrefix(r.URL.Path, "/api/thread/")
	slug = strings.TrimSuffix(slug, "/posts")
	fmt.Println(slug)

	thread, err := f.ForumUseCase.ThreadDetails(slug)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	posts, err := f.ForumUseCase.GetPostsOfThread(thread.Id, params, sort)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(posts)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}