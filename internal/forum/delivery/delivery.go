package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	domain "technopark-dbms-forum/internal/forum"
	"technopark-dbms-forum/models"
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
	r.HandleFunc("/api/thread/{slug_or_id}/details", handler.UpdateThread).Methods(http.MethodPost)

	r.HandleFunc("/api/service/status", handler.StatusDB).Methods(http.MethodGet)
	r.HandleFunc("/api/service/clear", handler.ClearDB).Methods(http.MethodPost)

	r.HandleFunc("/api/thread/{slug_or_id}/vote", handler.MakeVote).Methods(http.MethodPost)

	r.HandleFunc("/api/post/{id}/details", handler.PostUpdate).Methods(http.MethodPost)
	r.HandleFunc("/api/post/{id}/details", handler.PostDetails).Methods(http.MethodGet)

	r.HandleFunc("/api/forum/{slug}/threads", handler.ThreadsOfForum).Methods(http.MethodGet)
	r.HandleFunc("/api/forum/{slug}/users", handler.UsersOfForum).Methods(http.MethodGet)
}

func (f *ForumHandler) Forum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	forum := models.Forum{}
	err := json.NewDecoder(r.Body).Decode(&forum)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	forum, err = f.ForumUseCase.Forum(forum)
	status := models.GetStatusCodePost(err)
	if status == 409 {

		w.WriteHeader(models.GetStatusCodePost(err))

		body, err := json.Marshal(forum)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		w.Write(body)
		return
	}
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(forum)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodePost(err))
	w.Write(body)
}

func (f *ForumHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := strings.TrimPrefix(r.URL.Path, "/api/forum/")
	slug = strings.TrimSuffix(slug, "/create")


	var thread models.Thread
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	thread.Forum = slug
	check := thread.Slug



	thread, err = f.ForumUseCase.CreatingThread(thread)
 	status := models.GetStatusCodePost(err)
 	if status == 409 {

		w.WriteHeader(models.GetStatusCodePost(err))

		body, err := json.Marshal(thread)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		w.Write(body)
		return
	}

	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	if check == "" {
		threadOut := models.ThreadOut{
			Id:      thread.Id,
			Title:   thread.Title,
			Author:  thread.Author,
			Forum:   thread.Forum,
			Message: thread.Message,
			Votes:   thread.Votes,
			Slug:    thread.Slug,
			Created: thread.Created,
		}
		body, err := json.Marshal(threadOut)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		w.WriteHeader(models.GetStatusCodePost(err))
		w.Write(body)
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodePost(err))
	w.Write(body)
}

func (f *ForumHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nickname := strings.TrimPrefix(r.URL.Path, "/api/user/")
	nickname = strings.TrimSuffix(nickname, "/create")
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Nickname = nickname

	users, err := f.ForumUseCase.CreateUser(user)
	status := models.GetStatusCodePost(err)

	if status == 409 {
		body, err := json.Marshal(users)
		if err != nil {

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
	w.Header().Set("Content-Type", "application/json")
	nickname := strings.TrimPrefix(r.URL.Path, "/api/user/")
	nickname = strings.TrimSuffix(nickname, "/profile")


	user, err := f.ForumUseCase.GetUser(nickname)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(user)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodeGet(err))
	w.Write(body)
}

func (f *ForumHandler) ChangeProfileInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nickname := strings.TrimPrefix(r.URL.Path, "/api/user/")
	nickname = strings.TrimSuffix(nickname, "/profile")


	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Nickname = nickname

	userModel, err := f.ForumUseCase.ChangeUserProfile(user)
	if err != nil {
		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(userModel)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodeGet(err))
	w.Write(body)
}

func (f *ForumHandler) ForumInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := strings.TrimPrefix(r.URL.Path, "/api/forum/")
	slug = strings.TrimSuffix(slug, "/details")


	forum, err := f.ForumUseCase.ForumDetails(slug)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(forum)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodeGet(err))
	w.Write(body)
}

func (f *ForumHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := strings.TrimPrefix(r.URL.Path, "/api/thread/")
	slug = strings.TrimSuffix(slug, "/create")


	thread, err := f.ForumUseCase.ThreadDetails(slug)
	if err != nil {

		w.WriteHeader(models.GetStatusCodePost(err))
		w.Write(JSONError(err.Error()))
		return
	}


	var posts []models.Post
	err = json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		w.WriteHeader(201)
		w.Write([]byte("[]"))
		return
	}

	//author := (posts)[0].Author
	newPosts, err := f.ForumUseCase.CreatePosts(&posts, thread)
	//if len(*newPosts) == 0 {
	//	err = pgx.ErrNoRows
	//}
	if err != nil {

		//if err == pgx.ErrNoRows {
		//	//_, err = f.ForumUseCase.ThreadDetails(slug)
		//	//if err != nil {
		//	//
		//	//	w.WriteHeader(models.GetStatusCodePost(err))
		//	//	w.Write(JSONError(err.Error()))
		//	//	return
		//	//}
		//	_, err = f.ForumUseCase.GetUser(author)
		//	if err != nil {
		//
		//		w.WriteHeader(models.GetStatusCodePost(err))
		//		w.Write(JSONError(err.Error()))
		//		return
		//	}
		//	Err := models.ErrConflict
		//	w.WriteHeader(models.GetStatusCodePost(Err))
		//	w.Write(JSONError(Err.Error()))
		//	return
		//}

		w.WriteHeader(models.GetStatusCodePost(err))
		w.Write(JSONError(err.Error()))
		return
	}


	body, err := json.Marshal(*newPosts)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodePost(err))
	w.Write(body)
}

func (f *ForumHandler) ThreadDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slug := strings.TrimPrefix(r.URL.Path, "/api/thread/")
	slug = strings.TrimSuffix(slug, "/details")


	thread, err := f.ForumUseCase.ThreadDetails(slug)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	if models.IsUuid(thread.Slug) {
		result := models.ThreadToThreadOut(thread)
		body, err := json.Marshal(result)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(body)
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(models.GetStatusCodeGet(err))
	w.Write(body)
}

func (f *ForumHandler) StatusDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := f.ForumUseCase.StatusDB()
	body, err := json.Marshal(status)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func (f *ForumHandler) ClearDB(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	err := f.ForumUseCase.ClearDB()
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
}

func (f *ForumHandler) MakeVote(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	slug := strings.TrimPrefix(r.URL.Path, "/api/thread/")
	slug = strings.TrimSuffix(slug, "/vote")

	var thread models.Thread
	id, err := strconv.Atoi(slug)
	if err != nil {
		thread, err = f.ForumUseCase.ThreadDetails(slug)
	}

	//thread, err := f.ForumUseCase.ThreadDetails(slug)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	var vote models.Vote
	err = json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if thread.Id != 0 {
		id = thread.Id
	}

	vote.Thread = id

	thread, err = f.ForumUseCase.MakeVote(vote, thread)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}


	thread, err = f.ForumUseCase.ThreadDetails(slug)

	if models.IsUuid(thread.Slug) {
		result := models.ThreadToThreadOut(thread)
		body, err := json.Marshal(result)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(body)
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func (f *ForumHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := strings.TrimPrefix(r.URL.Path, "/api/post/")
	slug = strings.TrimSuffix(slug, "/details")
	id, err := strconv.Atoi(slug)
	if err != nil {

		w.WriteHeader(400)
		w.Write(JSONError(err.Error()))
		return
	}

	var postUpdate models.PostUpdate
	err = json.NewDecoder(r.Body).Decode(&postUpdate)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postUpdate.ID = id

	post, err := f.ForumUseCase.UpdateMessagePost(postUpdate)
	if err != nil {

		w.WriteHeader(models.GetStatusCodePost(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(post)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func (f *ForumHandler) PostDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := strings.TrimPrefix(r.URL.Path, "/api/post/")
	slug = strings.TrimSuffix(slug, "/details")
	id, err := strconv.Atoi(slug)
	if err != nil {

		w.WriteHeader(400)
		w.Write(JSONError(err.Error()))
		return
	}

	related := r.URL.Query().Get("related")

	postFull, err := f.ForumUseCase.PostFullDetails(id, related)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}


	body, err := json.Marshal(postFull)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func (f *ForumHandler) ThreadsOfForum(w http.ResponseWriter, r *http.Request) {


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



	threads, err := f.ForumUseCase.ListThreads(slug, params)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	var result []interface{}
	for _, thr := range threads {
		if models.IsUuid(thr.Slug) {
			tOut := models.ThreadToThreadOut(thr)

			result = append(result, tOut)
		} else {
			result = append(result, thr)
		}
	}

	body, err := json.Marshal(result)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	if len(threads) != 0 {
		w.Write(body)
	} else {
		w.Write([]byte("[]"))
	}
}

func (f *ForumHandler) UsersOfForum(w http.ResponseWriter, r *http.Request) {


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


	users, err := f.ForumUseCase.GetUsersByForum(slug, params)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(users)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}


	w.WriteHeader(http.StatusOK)
	if len(users) != 0 {
		w.Write(body)
	} else {
		w.Write([]byte("[]"))
	}
}


func (f *ForumHandler) PostsOfThread(w http.ResponseWriter, r *http.Request) {


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


	thread, err := f.ForumUseCase.ThreadDetails(slug)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	posts, err := f.ForumUseCase.GetPostsOfThread(thread.Id, params, sort)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	body, err := json.Marshal(posts)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(200)
	if len(posts) != 0 {
		w.Write(body)
	} else {
		w.Write([]byte("[]"))
	}

}

func (f *ForumHandler) UpdateThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slugOrId := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/thread/"), "/details")

	var thread models.Thread
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread.Slug = slugOrId
	} else {
		thread.Id = id
	}

	thread, err = f.ForumUseCase.UpdateThread(thread)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	if models.IsUuid(thread.Slug) {
		result := models.ThreadToThreadOut(thread)
		body, err := json.Marshal(result)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(JSONError(err.Error()))
			return
		}

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(body)

		return
	}

	body, err := json.Marshal(thread)
	if err != nil {

		w.WriteHeader(models.GetStatusCodeGet(err))
		w.Write(JSONError(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}