package main;

import(
	"fmt"
	"net/http"
	"github.com/jinzhu/gorm"
	"encoding/json"
	"time"
	"github.com/gorilla/mux"
  _ "github.com/jinzhu/gorm/dialects/postgres"
	
)
type ClimateResponse struct{
	Climate []Climate
	Pdf string
}
func getClimateBycity(w http.ResponseWriter,r *http.Request){
	db,err:=gorm.Open("postgres","user=postgres dbname=climate password=root sslmode=disable")
	if err!=nil{
		fmt.Println(err)
	}
	vars:=mux.Vars(r)
	city:=vars["city"]
	var cli []Climate
	db.Where("city=?",city).Order("date desc").Find(&cli)
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Content-Type","application/json")
	var res ClimateResponse
	res.Climate=cli
	res.Pdf="string"
	json.NewEncoder(w).Encode(res)
	defer db.Close()
}
func getClimateBycityDB(w http.ResponseWriter,r *http.Request){
	db,err:=gorm.Open("postgres","user=postgres dbname=climate password=root sslmode=disable")
	if err!=nil{
		fmt.Println(err)
	}
	defer db.Close()
	vars:=mux.Vars(r)
	city:=vars["city"]
	var cli Climate
	db.Where("city=?",city).Order("date desc").Find(&cli)
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(cli)
}	
func sendPdfMail(w http.ResponseWriter,r *http.Request){
	db,err:=gorm.Open("postgres","user=postgres dbname=climate password=root sslmode=disable")
	if err!=nil{
		fmt.Println(err)
	}
	defer db.Close()
	vars:=mux.Vars(r)
	city:=vars["city"]
	mail:=vars["mail"]
	var m MailLog
	m.From=mail
	m.To=mail
	m.City=city
	m.Date=time.Now().String()
	db.Create(&m)
	emailPdf(city,mail,mail)
	fmt.Fprintf(w,"Mail Sent")
}

func viewPdf(w http.ResponseWriter,r *http.Request){
	vars:=mux.Vars(r)
	city:=vars["city"]
	mail:=vars["mail"]
	createPdf(city,mail)
	fmt.Fprintf(w,"Ok")
}