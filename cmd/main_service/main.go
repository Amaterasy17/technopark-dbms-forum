package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"technopark-dbms-forum/configs"

	forumHandlers "technopark-dbms-forum/internal/forum/delivery"
	forumRepo "technopark-dbms-forum/internal/forum/repository/postgres"
	forumUseCase "technopark-dbms-forum/internal/forum/usecase"
)


func main() {
	router := mux.NewRouter()
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=%s",
		configs.PostgresPreferences.User,
		configs.PostgresPreferences.Password,
		configs.PostgresPreferences.DBName,
		configs.PostgresPreferences.Port)

	pgxConn, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pgxConn.PreferSimpleProtocol = true

	config := pgx.ConnPoolConfig{
		ConnConfig:     pgxConn,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	connPool, err := pgx.NewConnPool(config)
	if err != nil {
		fmt.Println(err.Error())
	}

	forumRepository := forumRepo.NewPostgresForumRepository(connPool)
	forumUsecase := forumUseCase.NewForumUsecase(forumRepository)
	forumHandlers.NewForumHandler(router, forumUsecase)

	addr := ":8080"
	err = http.ListenAndServe(addr, router)
	if err != nil {
		fmt.Println("error of starting server")
	}
}