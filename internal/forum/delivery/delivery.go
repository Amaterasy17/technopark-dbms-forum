package delivery

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	domain "technopark-dbms-forum/internal/forum"
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
	w.WriteHeader(http.StatusOK)
}

func (f *ForumHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	slug, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/forum/"))
	if err != nil {
		return
	}
	print(slug)
}