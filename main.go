package main

import (
	"net/http"

	"github.com/NoSkillGirl/user-service/routers"
)

func main() {
	// UserRoutes Initilization
	routers.UserRoutes()
	http.ListenAndServe(":8082", nil)
}
