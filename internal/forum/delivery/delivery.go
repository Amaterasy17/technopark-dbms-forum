package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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

	r.HandleFunc("/forum/create", handler.Forum).Methods(http.MethodPost)
	r.HandleFunc("/forum/{slug}/create", handler.CreateThread).Methods(http.MethodPost)
	r.HandleFunc("/user/{nickname}/create", handler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/user/{nickname}/profile", handler.ProfileUser).Methods(http.MethodGet)
}

func (f *ForumHandler) Forum(w http.ResponseWriter, r *http.Request) {
	forum := models.Forum{}
	err := json.NewDecoder(r.Body).Decode(&forum)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}



	w.WriteHeader(http.StatusOK)
}

func (f *ForumHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/forum/")
	slug = strings.TrimSuffix(slug, "/create")
	fmt.Println(slug)
}

func (f *ForumHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create user")
	nickname := strings.TrimPrefix(r.URL.Path, "/user/")
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
	nickname := strings.TrimPrefix(r.URL.Path, "/user/")
	nickname = strings.TrimSuffix(nickname, "/create")
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
	return
}