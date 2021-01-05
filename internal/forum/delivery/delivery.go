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

	r.HandleFunc("/api/forum/create", handler.Forum).Methods(http.MethodPost)
	r.HandleFunc("/api/forum/{slug}/create", handler.CreateThread).Methods(http.MethodPost)
	r.HandleFunc("/api/user/{nickname}/create", handler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/api/user/{nickname}/profile", handler.ProfileUser).Methods(http.MethodGet)
	r.HandleFunc("/api/user/{nickname}/profile", handler.ChangeProfileInformation).Methods(http.MethodPost)
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
	slug := strings.TrimPrefix(r.URL.Path, "/api/forum/")
	slug = strings.TrimSuffix(slug, "/create")
	fmt.Println(slug)
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