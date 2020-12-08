package delivery

import (
	"github.com/gorilla/mux"
	"net/http"
	domain "technopark-dbms-forum/internal/forum/usecase"
)

type ForumHandler struct {
	ForumUseCase domain.ForumUsecase
}

func NewForumHandler(r *mux.Router, forumUseCase domain.ForumUsecase) {
	handler := &ForumHandler{ForumUseCase: forumUseCase}

	r.HandleFunc("/forum", handler.Forum).Methods(http.MethodPost)
}

func (f *ForumHandler) Forum(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}