package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

//Response struct
type Response struct {
	Status   int32
	Response ResponseMsg
	Error    ErrorMessage
}

//ResponseMsgV2 struct
type ResponseMsgV2 struct {
	Msg  string
	User []models.User
}

//ResponseV2 struct
type ResponseV2 struct {
	Status   int32
	Response ResponseMsgV2
	Error    ErrorMessage
}

//ResponseMsgV3 struct
type ResponseMsgV3 struct {
	Msg string
	Bus []models.BusDetail
}

//ResponseV3 struct
type ResponseV3 struct {
	Status   int32
	Response ResponseMsgV3
	Error    ErrorMessage
}

//BusSearchRequest struct
type BusSearchRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	TravelDate  string `json:"travelDate"`
}

//BookingRequest struct
type BookingRequest struct {
	UserID       int
	BusID        int
	NoOfSeats    int
	DateOfTravel string
}

//ShowAllUser function
func ShowAllUser(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	users = models.GetAllUsers()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
}

//SearchBus function
func SearchBus(w http.ResponseWriter, r *http.Request) {
	var reqJSON BusSearchRequest
	err := json.NewDecoder(r.Body).Decode(&reqJSON)
	if err != nil {
		panic(err)
	}

	b, errOccured := models.SearchBus(reqJSON.Source, reqJSON.Destination, reqJSON.TravelDate)

	if errOccured == false {
		data := ResponseV3{
			Status: 200,
			Response: ResponseMsgV3{
				Bus: b,
			},
			Error: ErrorMessage{},
		}
		if len(b) > 0 {
			data.Response.Msg = "Buses avaliable are: "
		} else {
			data.Response.Msg = "Bus not avaliable"
		}
		json.NewEncoder(w).Encode(data)
	} else if errOccured == true {
		data := Response{
			Status: 500,
			Response: ResponseMsg{
				Msg: "Internal server error",
			},
			Error: ErrorMessage{},
		}
		json.NewEncoder(w).Encode(data)
	}
}

//NewBooking function
func NewBooking(w http.ResponseWriter, r *http.Request) {

	// Req Obj
	var reqJSON BookingRequest

	// Res Obj
	resp := Response{}
	w.Header().Set("Content-Type", "application/json")

	// Req Decode
	err := json.NewDecoder(r.Body).Decode(&reqJSON)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Add Booking
	errOccured := models.AddBooking(reqJSON.UserID, reqJSON.BusID, reqJSON.NoOfSeats, reqJSON.DateOfTravel)
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

		// Should send SMS to the customer with booking details.
		type SMSRequest struct {
			Mobile  string
			Message string
		}

		smsRequest := SMSRequest{
			Mobile:  "+918904621381",
			Message: "Your Bus booking was successful",
		}

		smsRequestInBytes, err := json.Marshal(smsRequest)

		req, err := http.NewRequest("POST", "http://localhost:8081/SendSMS", bytes.NewBuffer(smsRequestInBytes))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error in reading response body from smsService", err)
		}

		fmt.Println(string(body))
	}
	json.NewEncoder(w).Encode(resp)
}
