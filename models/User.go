package models

import (
	"context"
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
	ID             int32
	Number         string
	AC             bool
	Sleeper        bool
	TotalSeat      int32
	Source         string
	Destination    string
	CompanyID      int32
	DepartureTime  string
	ArrivalTime    string
	SeatsAvailable int32
}

// BookingDetail Struct
type BookingDetail struct {
	ID        int32
	UserID    int32
	BusID     int32
	NoOfSeats int32
	Date      string
}

var (
	ctx context.Context
	db  *sql.DB
)

const mySQLHost = "34.93.137.151"

var mySQLConnection = fmt.Sprintf("root:password@tcp(%s)/tour_travel", mySQLHost)

//GetAllUsers function
func GetAllUsers() (users []User) {
	//db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tour_travel")
	db, err := sql.Open("mysql", mySQLConnection)
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
func AddUser(name, phoneNo, emailID, password string) (errorOccured bool, duplicateUser bool) {
	ctx := context.Background()
	//db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tour_travel")
	db, err := sql.Open("mysql", mySQLConnection)

	// if there is an error opening the connection, handle it
	if err != nil {
		//panic(err.Error())
		return true, false
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	var count int
	error := db.QueryRowContext(ctx, "select count(*) from user_details where name=? and (phone_no=? or email_id=?)", name, phoneNo, emailID).Scan(&count)

	if error != nil {
		return true, false
	}
	if count == 0 {
		addUserQuery := `INSERT INTO user_details (name, phone_no, email_id, password) VALUES ('%s', '%s', '%s', '%s')`

		addUserQueryString := fmt.Sprintf(addUserQuery, name, phoneNo, emailID, password)
		fmt.Println(addUserQueryString)

		// perform a db.Query insert
		insert, err := db.Query(addUserQueryString)

		// if there is an error inserting, handle it
		if err != nil {
			//panic(err.Error())
			return true, false
		}

		// be careful deferring Queries if you are using transactions
		defer insert.Close()
		return false, false
	}
	//error = false, duplicate = true
	return false, true

}

//SearchBus function
func SearchBus(source, destination string, travelDate string) (busDetails []BusDetail, errorOccured bool) {
	// opening mysql connection
	//db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tour_travel")
	db, err := sql.Open("mysql", mySQLConnection)

	if err != nil {
		return []BusDetail{}, true
	}
	defer db.Close()

	// bus_details query preparation
	selectBusDetailQuery := `select * from bus_details where source = '%s' and destination = '%s';`
	selectBusDetailQueryString := fmt.Sprintf(selectBusDetailQuery, source, destination)
	// bus_details running in db
	rows, err := db.Query(selectBusDetailQueryString)
	if err != nil {
		return []BusDetail{}, true
	}
	defer rows.Close()

	// seats_available query preparation
	seatsAvailableQuery := `select bus_id, sum(no_of_seats) from booking_details where bus_id in (select id from bus_details where source = '%s' and destination = '%s') and date_of_travel = '%s' group by bus_id;`
	seatsAvailableQueryString := fmt.Sprintf(seatsAvailableQuery, source, destination, travelDate)
	fmt.Println(seatsAvailableQueryString)

	// seats_available running in db
	rs, err := db.Query(seatsAvailableQueryString)
	if err != nil {
		return []BusDetail{}, true
	}
	defer rs.Close()
	busSeatAvailability := make(map[int32]int32)

	// extracting the results
	for rs.Next() {
		var busID, seatsAvailable int32
		err = rs.Scan(&busID, &seatsAvailable)
		if err != nil {
			return []BusDetail{}, true
		}
		busSeatAvailability[busID] = seatsAvailable
	}

	err = rs.Err()
	if err != nil {
		return []BusDetail{}, true
	}

	// extracting the results
	for rows.Next() {
		b := BusDetail{}
		err = rows.Scan(&b.ID, &b.Number, &b.AC, &b.Sleeper, &b.TotalSeat, &b.Source, &b.Destination, &b.CompanyID, &b.DepartureTime, &b.ArrivalTime)
		if err != nil {
			return []BusDetail{}, true
		}
		if busSeatAvailability[b.ID] != 0 {
			b.SeatsAvailable = b.TotalSeat - busSeatAvailability[b.ID]
		} else {
			b.SeatsAvailable = b.TotalSeat
		}

		busDetails = append(busDetails, b)
	}

	err = rows.Err()
	if err != nil {
		return []BusDetail{}, true
	}

	return busDetails, false
}

//AddBooking function
func AddBooking(userID int, busID int, noOfSeats int, travelDate string) (errorOccured bool) {

	// initilize db connection
	//db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tour_travel")
	db, err := sql.Open("mysql", mySQLConnection)

	if err != nil {
		fmt.Println("Error in initializing db connection", err)
		return true
	}
	defer db.Close()

	// t := time.Now().Format("2006-01-02")
	addBookingQuery := `
	INSERT INTO booking_details (user_id, bus_id, no_of_seats, date_of_travel, date_of_booking) 
	VALUES (%d, %d, %d, '%s', now())
	`
	addBookingQueryString := fmt.Sprintf(addBookingQuery, userID, busID, noOfSeats, travelDate)
	fmt.Println(addBookingQueryString)

	// perform a db.Query insert
	insert, err := db.Query(addBookingQueryString)
	if err != nil {
		fmt.Println(err)
		return true
	}
	defer insert.Close()

	return false
}

//UserExist function
func UserExist(name, password string) (user []User, errorOccured bool) {
	//db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/tour_travel")
	db, err := sql.Open("mysql", mySQLConnection)
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
