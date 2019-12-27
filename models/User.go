package models

import (
	"database/sql"
	"fmt"

	//mysql
	_ "github.com/go-sql-driver/mysql"
)

// User - Struct
type User struct {
	ID       int32
	Name     string
	PhoneNo  string
	EmailID  string
	Password string
}

// Company Struct
type Company struct {
	ID        int32
	Name      string
	OwnerName string
	PhoneNo   string
}

// BusDetail Struct
type BusDetail struct {
	ID          int32
	Number      string
	AC          bool
	Sleeper     bool
	TotalSeat   int32
	Source      string
	Destination string
	CompanyID   int32
}

// BookingDetail Struct
type BookingDetail struct {
	ID        int32
	UserID    int32
	BusID     int32
	NoOfSeats int32
	Date      string
}

//GetAllUsers function
func GetAllUsers() (users []User) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tour_travel")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	selectAllUsersQuery := `SELECT id, name, phone_no, email_no FROM user_details`

	// perform a db.Query select
	rows, err := db.Query(selectAllUsersQuery)

	// if there is an error, handle it
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {

		u := User{}

		err = rows.Scan(&u.ID, &u.Name, &u.PhoneNo, &u.EmailID)

		if err != nil {
			// handle this error
			panic(err)
		}
		users = append(users, u)

	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return users

}

//AddUser function
func AddUser(name, phoneNo, emailID, password string) (errorOccured bool) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tour_travel")

	// if there is an error opening the connection, handle it
	if err != nil {
		//panic(err.Error())
		return true
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	duplicateCheckQuery := `select count(*) from user_details where name = '%s' and (phone_no = '%s' or email_id = '%s')`
	duplicateCheckQueryString := fmt.Sprintf(duplicateCheckQuery, name, phoneNo, emailID)
	check, err := db.Query(duplicateCheckQueryString)
	fmt.Println("check : ", check)

	addUserQuery := `INSER INTO user_details (name, phone_no, email_id, password) VALUES ('%s', '%s', '%s', '%s')`

	addUserQueryString := fmt.Sprintf(addUserQuery, name, phoneNo, emailID, password)
	fmt.Println(addUserQueryString)

	// perform a db.Query insert
	insert, err := db.Query(addUserQueryString)

	// if there is an error inserting, handle it
	if err != nil {
		//panic(err.Error())
		return true
	}

	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	return false
}

//SearchBus function
func SearchBus(source, destination string) (busDetails []BusDetail) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tour_travel")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	selectBusDetailQuery := `select * from bus_details where source = '%s' and destination = '%s'`

	selectBusDetailQueryString := fmt.Sprintf(selectBusDetailQuery, source, destination)
	fmt.Println(selectBusDetailQueryString)

	// perform a db.Query select
	rows, err := db.Query(selectBusDetailQueryString)

	// if there is an error, handle it
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		b := BusDetail{}
		err = rows.Scan(&b.ID, &b.Number, &b.AC, &b.Sleeper, &b.TotalSeat, &b.Source, &b.Destination, &b.CompanyID)

		if err != nil {
			// handle this error
			panic(err)
		}
		busDetails = append(busDetails, b)

	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return busDetails
}

//AddBooking function
func AddBooking(userID int, busID int, noOfSeats int, date string) (errorOccured bool) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tour_travel")

	// if there is an error opening the connection, handle it
	if err != nil {
		//panic(err.Error())
		return true
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	addBookingQuery := `INSERT INTO booking_details (user_id, bus_id, no_of_seats, date) VALUES (%d, %d, %d, '%s')`

	addBookingQueryString := fmt.Sprintf(addBookingQuery, userID, busID, noOfSeats, date)
	fmt.Println(addBookingQueryString)

	// perform a db.Query insert
	insert, err := db.Query(addBookingQueryString)

	// if there is an error inserting, handle it
	if err != nil {
		//panic(err.Error())
		return true
	}

	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	return false
}

//UserExist function
func UserExist(name, password string) (user []User, errorOccured bool) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tour_travel")

	// if there is an error opening the connection, handle it
	if err != nil {
		// panic(err.Error())
		return user, true
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	searchUserQuery := `select * from user_details where name = '%s' and password = '%s' limit 1;`

	searchUserQueryString := fmt.Sprintf(searchUserQuery, name, password)
	fmt.Println(searchUserQueryString)

	// perform a db.Query search
	search, err := db.Query(searchUserQueryString)

	// if there is an error inserting, handle it
	if err != nil {
		// panic(err.Error())
		return user, true
	}

	// be careful deferring Queries if you are using transactions
	defer search.Close()

	for search.Next() {
		u := User{}
		err = search.Scan(&u.ID, &u.Name, &u.PhoneNo, &u.EmailID, &u.Password)

		if err != nil {
			// handle this error
			//panic(err)
			return user, true
		}
		user = append(user, u)
	}

	return user, false
}
