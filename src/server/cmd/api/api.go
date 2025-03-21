package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/cebuh/simpleHolidayPlaner/service/invite"
	"github.com/cebuh/simpleHolidayPlaner/service/team"
	"github.com/cebuh/simpleHolidayPlaner/service/user"
	"github.com/cebuh/simpleHolidayPlaner/service/vacation"
	"github.com/gorilla/mux"
)

type Server struct {
	address string
	db      *sql.DB
}

func NewServer(address string, db *sql.DB) *Server {
	return &Server{
		address: address,
		db:      db,
	}
}

func (s *Server) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	teamStore := team.NewStore(s.db)
	teamHandler := team.NewHandler(s.db, teamStore, userStore)
	teamHandler.RegisterRoutes(subrouter)

	inviteStore := invite.NewStore(s.db)
	inviteHandler := invite.NewHandler(s.db, inviteStore, userStore, teamStore)
	inviteHandler.RegisterRoutes(subrouter)

	vacationStore := vacation.NewStore(s.db)
	vacationHandler := vacation.NewHandler(s.db, userStore, teamStore, vacationStore)
	vacationHandler.RegisterRoutes(subrouter)

	log.Println("Listen on ", s.address)
	return http.ListenAndServe(s.address, router)
}
