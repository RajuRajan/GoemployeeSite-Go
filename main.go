package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

var db *gorm.DB
var err error

type test_struct struct {
	Test string
}
type Users struct {
	gorm.Model
	Accounttype string
	Name        string
	Empid       string
	Password    string
	Userid      string
}
type Leave struct {
	Empid       string
	Empname     string
	Description string
	Daysofleave string
	Accepted    string
	Date        string
	Dateend     string
	Type        string
}
type Adminemployeereg struct {
	Name              string
	EmployeeId        string
	Sex               string
	Age               string
	Designation       string
	Qualification     string
	YearsOfExperiance string
	Imgsrc            string
}

func InitialMigration() {

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect")
	} else {
		fmt.Println("Connected successfully")
	}
	defer db.Close()
	db.AutoMigrate(&Users{})
	db.AutoMigrate(&Adminemployeereg{})
	db.AutoMigrate(&Leave{})
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Helloworld")
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/user/{name}/{empid}/{accounttype}/{password}/{userid}", Newuser).Methods("POST")
	myRouter.HandleFunc("/admin/", adminEmpReg).Methods("POST")
	myRouter.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("."+"/images/"))))
	myRouter.HandleFunc("/login/{userid}/{password}", login).Methods("POST")
	myRouter.HandleFunc("/emphome/{description}/{empid}/{daysofleave}/{name}/{date}/{dateend}/{type}", mail).Methods("POST")
	myRouter.HandleFunc("/pending", pending).Methods("GET")
	myRouter.HandleFunc("/acceptleave/{empname}/{empid}/{date}", acceptleave).Methods("POST")
	myRouter.HandleFunc("/rejectleave/{empname}/{empid}/{date}", rejectleave).Methods("POST")
	myRouter.HandleFunc("/approval/{empname}/{empid}", approval).Methods("POST")
	myRouter.HandleFunc("/emplist", emplist).Methods("GET")
	myRouter.HandleFunc("/get-token", GetTokenHandler).Methods("GET")
	myRouter.HandleFunc("/upload/{file}", fileupload).Methods("POST")
	myRouter.HandleFunc("/admin/test", test).Methods("POST")
	myRouter.HandleFunc("/leavecheck/{empid}/{empname}/{date}", leavecheck).Methods("POST")
	myRouter.HandleFunc("/pendingapproval/{empid}/{empname}", pendingapproval).Methods("POST")
	log.Fatal(http.ListenAndServe(":8123", cors.Default().Handler(myRouter)))
}

func main() {
	fmt.Println("Go ORM Tutorial")
	InitialMigration()
	handleRequests()
}
