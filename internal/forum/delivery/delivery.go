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

func NewForumHandler(r *mux.Router, forumUseCase domain.ForumUseCase) {
	handler := &ForumHandler{ForumUseCase: forumUseCase}

	r.HandleFunc("/forum/create", handler.Forum).Methods(http.MethodPost)
	r.HandleFunc("/forum/{slug}/create", handler.CreateThread).Methods(http.MethodPost)
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
	slug = strings.TrimRight(slug, "/create")
	fmt.Println(slug)
}