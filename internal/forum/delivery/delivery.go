package delivery

import (
	"github.com/gorilla/mux"
	"net/http"
	domain "technopark-dbms-forum/internal/forum"
)

type ForumHandler struct {
	ForumUseCase domain.ForumUseCase
}

func NewForumHandler(r *mux.Router, forumUseCase domain.ForumUseCase) {
	handler := &ForumHandler{ForumUseCase: forumUseCase}

	r.HandleFunc("/forum", handler.Forum).Methods(http.MethodPost)
}

func (f *ForumHandler) Forum(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}