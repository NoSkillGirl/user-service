package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NoSkillGirl/user-service/models"
)

//ErrorMessage struct
type ErrorMessage struct {
	Msg string
}

//ResponseMsg struct
type ResponseMsg struct {
	Msg string
}

//ResponseMsgV2 struct
type ResponseMsgV2 struct {
	Msg  string
	User []models.User
}

//Response struct
type Response struct {
	Status   int32
	Response ResponseMsg
	Error    ErrorMessage
}

//ResponseV2 struct
type ResponseV2 struct {
	Status   int32
	Response ResponseMsgV2
	Error    ErrorMessage
}

//ShowAllUser function
func ShowAllUser(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	users = models.GetAllUsers()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
}

//RegisterUser function
// func RegisterUser(w http.ResponseWriter, r *http.Request) {
// 	name := r.FormValue("name")
// 	phoneNo := r.FormValue("phone_no")
// 	emailID := r.FormValue("email_id")
// 	// fmt.Println(name, phoneNo, emailID)

// 	err := models.AddUser(name, phoneNo, emailID)
// 	w.Header().Set("Content-Type", "application/json")
// 	resp := Response{}
// 	if err == true {
// 		resp.Status = 500
// 		resp.Response = ResponseMsg{}
// 		resp.Error = ErrorMessage{
// 			Msg: "Internal Server Error",
// 		}

// 	} else {
// 		resp.Status = 200
// 		resp.Response = ResponseMsg{
// 			Msg: "user succesfully created",
// 		}
// 		resp.Error = ErrorMessage{}

// 	}
// 	json.NewEncoder(w).Encode(resp)
// }

//SearchBus function
func SearchBus(w http.ResponseWriter, r *http.Request) {
	source := r.FormValue("source")
	destination := r.FormValue("destination")
	fmt.Println(source, destination)

	busDetails := models.SearchBus(source, destination)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(busDetails)
}

//NewBooking function
func NewBooking(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	userIDInt, _ := strconv.Atoi(userID)

	busID := r.FormValue("bus_id")
	busIDInt, _ := strconv.Atoi(busID)
	noOfSeats := r.FormValue("no_of_seats")
	noOfSeatsInt, _ := strconv.Atoi(noOfSeats)
	date := r.FormValue("date")

	errOccured := models.AddBooking(userIDInt, busIDInt, noOfSeatsInt, date)

	w.Header().Set("Content-Type", "application/json")
	resp := Response{}
	if errOccured == true {
		resp.Status = 500
		resp.Response = ResponseMsg{}
		resp.Error = ErrorMessage{
			Msg: "Internal Server Error",
		}

	} else {
		resp.Status = 200
		resp.Response = ResponseMsg{
			Msg: "booking details succesfully added",
		}
		resp.Error = ErrorMessage{}

	}
	json.NewEncoder(w).Encode(resp)
}

//SearchUser function
// func SearchUser(w http.ResponseWriter, r *http.Request) {
// 	name := r.FormValue("name")
// 	password := r.FormValue("password")

// 	u, errOccured := models.UserExist(name, password)

// 	w.Header().Set("Content-Type", "application/json")
// 	resp := ResponseV2{}
// 	if errOccured == true {
// 		resp.Status = 500
// 		resp.Response = ResponseMsgV2{}
// 		resp.Error = ErrorMessage{
// 			Msg: "Internal Server Error",
// 		}

// 	} else {
// 		resp.Status = 200
// 		resp.Response = ResponseMsgV2{
// 			Msg:  "user found",
// 			User: u,
// 		}
// 		resp.Error = ErrorMessage{}
// 	}
// 	// fmt.Println(u)
// 	json.NewEncoder(w).Encode(resp)
// }
