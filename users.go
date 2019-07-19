package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func Newuser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := mux.Vars(r)
	name := vars["name"]
	empid := vars["empid"]
	accounttype := vars["accounttype"]
	password := vars["password"]
	userid := vars["userid"]

	db.Create(&Users{Name: name, Empid: empid, Password: password, Accounttype: accounttype, Userid: userid})

	fmt.Println(w, "New user created successfully")
	fmt.Fprintf(w, "New user created successfully")

}
func adminEmpReg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {

		fmt.Println("Connection Successfull")
	}
	defer db.Close()
	r.ParseMultipartForm(10 << 20)

	var a Adminemployeereg
	a.Name = r.FormValue("name")
	a.EmployeeId = r.FormValue("empid")
	a.Age = r.FormValue("age")
	a.Designation = r.FormValue("des")
	a.Qualification = r.FormValue("qual")
	a.Sex = r.FormValue("sex")
	a.YearsOfExperiance = r.FormValue("yoe")
	file, handler, err := r.FormFile("image")
	if err != nil {
		panic(err)
	}
	fmt.Println(handler.Header)
	filebytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	ty := strings.Split(http.DetectContentType(filebytes), "/")
	typ := ty[len(ty)-1:][0]
	fmt.Println(typ)
	tempFile, err := ioutil.TempFile("images", "*"+"."+typ)
	if err != nil {
		panic(err)
	}
	tempFile.Write(filebytes)
	defer tempFile.Close()
	a.Imgsrc = tempFile.Name()
	fmt.Println(tempFile.Name())
	db.Create(&a)
	fmt.Println(w, "New user created successfully")
	fmt.Fprintf(w, "New user created successfully")
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "txt/plain")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := mux.Vars(r)
	var user Users
	userid := vars["userid"]
	password := vars["password"]

	db.Where("userid=? AND password=?", userid, password).Find(&user)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["name"] = user.Name
	claims["empid"] = user.Empid
	claims["accounttype"] = user.Accounttype
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, _ := token.SignedString(mySigningKey)

	w.Write([]byte(tokenString))
	fmt.Println("token")

	fmt.Println(tokenString)
}	
func mail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := mux.Vars(r)
	name := vars["name"]
	empid := vars["empid"]
	description := vars["description"]
	daysofleave := vars["daysofleave"]
	date := vars["date"]
	dateend := vars["dateend"]
	types := vars["type"]
	acceptance := "no"

	db.Create(&Leave{Empname: name, Empid: empid, Daysofleave: daysofleave, Description: description, Accepted: acceptance, Date: date, Dateend: dateend, Type: types})

	fmt.Println(w, "reason ")
	fmt.Fprintf(w, "reason")
	send(description, name, empid)
}

var mySigningKey = []byte("secret")

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

})

func send(body string, name string, empid string) {
	from := "rajuart678@gmail.com"
	pass := "veluprabha21669"
	to := "rajuart678@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Permission for leave\n\n" +
		"Respected sir/Mam\n\n" +
		body + "\n\n" + name + "\n" + empid

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, visit http://foobarbazz.mailinator.com")
}
func pending(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {
		fmt.Println("Connection Successfull")
	}
	var leave []Leave
	db.Find(&leave)
	fmt.Println("pending")
	json.NewEncoder(w).Encode(leave)

}
func acceptleave(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := mux.Vars(r)
	var leave []Leave
	empid := vars["empid"]
	empname := vars["empname"]
	date := vars["date"]

	db.Model(&leave).Where("Empid = ? AND Empname =? AND Date=?", empid, empname, date).Update("accepted", "yes")

}
func rejectleave(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {
		fmt.Println("Connection Successfull")
	}

	vars := mux.Vars(r)
	var leave []Leave
	empid := vars["empid"]
	empname := vars["empname"]
	date := vars["date"]

	db.Model(&leave).Where("Empid = ? AND Empname =? AND Date=?", empid, empname, date).Update("accepted", "rejected")

}
func approval(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {
		fmt.Println("Connection Successfull")
	}
	vars := mux.Vars(r)
	var leave Leave
	empid := vars["empid"]
	empname := vars["empname"]
	db.Model(&leave).Where("Empid = ? AND Empname =?", empid, empname).Find(&leave)
	fmt.Println("approval")
	json.NewEncoder(w).Encode(leave)

}

func emplist(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {
		fmt.Println("Connection Successfull")
	}
	var emplist []Adminemployeereg
	db.Find(&emplist)
	fmt.Println("pending")
	json.NewEncoder(w).Encode(emplist)
}

func fileupload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("file uploaded")

	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("name")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	fmt.Println(file)
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
func test(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	log.Println(req.Form)
	//LOG: map[{"test": "that"}:[]]
	var t test_struct
	for key, _ := range req.Form {
		log.Println(key)
		//LOG: {"test": "that"}
		err := json.Unmarshal([]byte(key), &t)
		if err != nil {
			log.Println(err.Error())
		}
	}
	log.Println(t.Test)
	//LOG: that
}
func leavecheck(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {
		fmt.Println("Connection Successfull")
	}
	vars := mux.Vars(r)
	var leave Leave
	empid := vars["empid"]
	empname := vars["empname"]
	date := vars["date"]
	db.Model(&leave).Where("Empid = ? AND Empname =? AND Date=?", empid, empname, date).Find(&leave)
	fmt.Println("checkleave")
	json.NewEncoder(w).Encode(leave)

}
func pendingapproval(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err = gorm.Open("postgres", "port=5432 user=postgres dbname=cm password=root sslmode=disable")
	if err != nil {
		panic("Failed to connect")
	} else {
		fmt.Println("Connection Successfull")
	}
	vars := mux.Vars(r)
	var leave []Leave
	empid := vars["empid"]
	empname := vars["empname"]

	db.Model(&leave).Where("Empid = ? AND Empname =? ", empid, empname).Find(&leave)
	fmt.Println("pendingapproval")
	json.NewEncoder(w).Encode(leave)

}
