package routers

import (
	"net/http"

	"github.com/NoSkillGirl/user-service/controllers"
)

// UserRoutes - All User related Routes
func UserRoutes() {
	http.HandleFunc("/users/register", controllers.RegisterUser)
	http.HandleFunc("/users", controllers.ShowAllUser)
	http.HandleFunc("/search", controllers.SearchBus)
	http.HandleFunc("/booking/new", controllers.NewBooking)
	http.HandleFunc("/searchUser", controllers.SearchUser)
}
