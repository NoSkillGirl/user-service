package main

import (
	"fmt"
	"net/http"

	// "github.com/NoSkillGirl/user-service/models"
	"github.com/NoSkillGirl/user-service/routers"
)

func main() {
	// UserRoutes Initilization

	// e, u := models.AddUser("Blah", "8908", "fkjs@gmail.com", "fdskjfks")
	// fmt.Println(e, u)
	routers.UserRoutes()
	http.ListenAndServe(":8082", nil)
}
