package main;

import (
	"fmt"
	"log"
	"net/http"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/rs/cors"
	"github.com/gorilla/mux"
	 "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct{
 gorm.Model
 FName string
 LName string
 Email string
 Pass string
 Image string
}

type Response struct{
	IsValid bool
	ID uint
	 FName string
	 LName string
	 Email string
	 Image string
}
type Climate struct{
	gorm.Model
	Date string
	City string
	MaxTempC float64
	MaxTempF float64
	MinTempC float64
	MinTempF float64
	MaxWindK float64
	MaxWindM float64
	Humidity float64
	TotalPrecipM float64
	TotalPrecipI float64
	Text string
	Icon string
	Code int

}
type SendOther struct{
	Sendmail string
	Sendsubject string
	Sendmessage string
	Recievemail string
	City string
}
type MailLog struct{
	gorm.Model
	From string
	To string
	City string
	Date string
}
var(
	googleOauthConfig=&oauth2.Config{
		RedirectURL:"http://localhost:8000/authcallback",
		ClientID:"862010389923-5itkgg95ir5ukrfottopvmdt67bqh30u.apps.googleusercontent.com",
		ClientSecret :"rzOSABFVTTjK7lYILMza-2SB",
		Scopes:[]string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	randomState="random"
)
func helloWorld(w http.ResponseWriter,r *http.Request){
	fmt.Fprintf(w,"Hello World")
}
func handleRequests(){
	myRouter:=mux.NewRouter().StrictSlash(true)
	myRouter.PathPrefix("/images/").Handler(http.StripPrefix("/images/",http.FileServer(http.Dir("."+"/images/"))))
	myRouter.PathPrefix("/pdf/").Handler(http.StripPrefix("/pdf/",http.FileServer(http.Dir("."+"/pdf/"))))
	myRouter.HandleFunc("/",helloWorld).Methods("GET")
	myRouter.HandleFunc("/getmailHistory/{email}",getMailHistory).Methods("GET")
	myRouter.HandleFunc("/create",createUser).Methods("POST")
	myRouter.HandleFunc("/auth",authUser).Methods("POST")
	myRouter.HandleFunc("/googlelog",loginwithgoogle)
	myRouter.HandleFunc("/authcallback",fetchGoogledetails)
	myRouter.HandleFunc("/mailother",mailOther).Methods("POST")
	myRouter.HandleFunc("/delete/{email}/",deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/update/{name}/{email}/{pass}",updateUser).Methods("PUT")
	myRouter.HandleFunc("/list/",listUsers).Methods("GET")
	myRouter.HandleFunc("/climate/{city}",getClimateBycity).Methods("GET")
	myRouter.HandleFunc("/climateDB/{city}",getClimateBycityDB).Methods("GET")
	myRouter.HandleFunc("/sendMail/{city}/{mail}",sendPdfMail).Methods("GET")
	myRouter.HandleFunc("/viewPdf/{city}/{mail}",viewPdf).Methods("GET")
	myRouter.HandleFunc("/getIp/",getServerIp).Methods("GET")
	myRouter.HandleFunc("/createFromGoole",createFromGoole).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000",cors.Default().Handler(myRouter)))
}
var db *gorm.DB
var err error


func main(){
	fmt.Println("Go running")
	db, err := gorm.Open("postgres", "user=postgres dbname=climate password=root sslmode=disable")
	if(err!=nil){
		panic(err)
		
	}
	defer db.Close()
	conn:=db.DB()
	err=conn.Ping()
	if(err!=nil){
		fmt.Println("D B   E R R O R")
		
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&MailLog{})
	fmt.Println("success")
	handleRequests()
}