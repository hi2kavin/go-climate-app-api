package main;

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"encoding/json"
    "github.com/gorilla/mux"
  _ "github.com/jinzhu/gorm/dialects/postgres"
	
)


func createUser(w http.ResponseWriter,r *http.Request){
	db, err := gorm.Open("postgres", "user=postgres dbname=climate password=root sslmode=disable")
	if(err!=nil){
		panic(err)
	}
	r.ParseMultipartForm(10<<20)
	file,handler,err:=r.FormFile("file")
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(handler.Header)
	defer file.Close()
	fileBytes,err:=ioutil.ReadAll(file)
	if err!=nil{
		fmt.Println(err)
	}
	ty:=strings.Split(http.DetectContentType(fileBytes),"/")
	typ:=ty[len(ty)-1:][0]
	fmt.Println(typ)
	temp,err:=ioutil.TempFile("images","*."+typ)
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer temp.Close()
	temp.Write(fileBytes)
	var user User
	user.FName=r.FormValue("fname")
	user.LName=r.FormValue("lname")
	user.Email=r.FormValue("email")
	user.Pass=r.FormValue("pass")
	user.Image="http://localhost:8000/"+temp.Name()
	db.Create(&user).Find(user)
	fmt.Println(user)
	defer db.Close()
	user.Pass=""
	json.NewEncoder(w).Encode(user)
}
func createFromGoole(w http.ResponseWriter,r *http.Request){
	decoder := json.NewDecoder(r.Body)
	 var user,ck User
	 err=decoder.Decode(&user)
	 db, err := gorm.Open("postgres", "user=postgres dbname=climate password=root sslmode=disable")
	 if(err!=nil){
	 	fmt.Fprintf(w,"updating gdjahfbijdhfisgllhgj")
	 }
	 defer db.Close()
	 
	 db.Where("Email=?",user.Email).Find(&ck)
	 if ck.ID==0{
		db.Create(&user).Find(&ck)
	 mailPassword(ck.Email,ck.Pass)
	 ck.Pass=""
	 w.Header().Set("Content-Type","application/json")
	 json.NewEncoder(w).Encode(&ck)
	 return 
	 }else{
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(&user)
	 }
}

func authUser(w http.ResponseWriter,r *http.Request){
	db, err := gorm.Open("postgres", "user=postgres dbname=climate password=root sslmode=disable")
	if(err!=nil){
		panic(err)
	}
	defer db.Close()
	decoder := json.NewDecoder(r.Body)
	var user User
	err=decoder.Decode(&user)
	pass:=user.Pass
	email:=user.Email
	user.Pass=""
	fmt.Println(email,pass)
	var res Response
	db.Where("email=?",email).Find(&user)
	fmt.Println(user)
	fmt.Println(user.ID)
	if user.Pass==pass && user.ID>0{
		res.IsValid=true
		res.FName=user.FName
		res.LName=user.LName
		res.ID=user.ID
		res.Email=user.Email
		res.Image=user.Image
		json.NewEncoder(w).Encode(res)
	}else{
		res.IsValid=false
		json.NewEncoder(w).Encode(res)
	}
	
	
}
func updateUser(w http.ResponseWriter,r *http.Request){
	db, err := gorm.Open("postgres", "user=postgres dbname=climate password=root sslmode=disable")
	if(err!=nil){
		panic(err)
	}
	vars:=mux.Vars(r)
	fname:=vars["fname"]
	email:=vars["email"]
    pass:=vars["pass"]
	var user User
	db.Where("email=?",email).Find(&user)
	user.FName=fname
	user.Pass=pass
	db.Save(user)
	defer db.Close()
	fmt.Fprintf(w,"updating user")
}
func getMailHistory(w http.ResponseWriter,r *http.Request){
	db, err := gorm.Open("postgres", "user=postgres dbname=climate password=root sslmode=disable")
	if(err!=nil){
		panic(err)
	}
	defer db.Close()
	vars:=mux.Vars(r)
	email:=vars["email"]
	var m []MailLog
	db.Where(&MailLog{From:email}).Find(&m)
	fmt.Println(email,m)
	json.NewEncoder(w).Encode(m)

}
func deleteUser(w http.ResponseWriter,r *http.Request){
	db, err := gorm.Open("postgres", "user=postgres dbname=climate password=root sslmode=disable")
	if(err!=nil){
		panic(err)
	}
	vars:=mux.Vars(r)
	email:=vars["email"]
	var user User
	db.Where("email=?",email).Find(&user)
	db.Delete(user)
	defer db.Close()
	fmt.Fprintf(w,"deleting user")
}
func listUsers(w http.ResponseWriter,r *http.Request){

db, err := gorm.Open("postgres", "user=postgres dbname=climate password=root sslmode=disable")
	if(err!=nil){
		panic(err)
	}
	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)

	defer db.Close()
}

func getServerIp(w http.ResponseWriter,r *http.Request){
 ip,err:=externalIP()
 if err!=nil{
	fmt.Fprintf(w,"NO_IP")
	 fmt.Println()
	 return 
 }

 fmt.Fprintf(w,ip)
}
func loginwithgoogle(w http.ResponseWriter,r *http.Request){
	url:=googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w,r,url,http.StatusSeeOther)
}
func fetchGoogledetails(w http.ResponseWriter,r *http.Request){
	if r.FormValue("state")!=randomState {
	fmt.Println("state is not valid")
	http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
	return
	}
	token,err:=googleOauthConfig.Exchange(oauth2.NoContext,r.FormValue("code"))
	if err!=nil{
		fmt.Println("error token: %s",err.Error())
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return 
	}
	
	resp,err:=http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token="+token.AccessToken)
	if err!=nil{
		fmt.Println("error cannt create request: %s",err.Error())
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return 
	}
	defer resp.Body.Close()
	content,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println("error cannot parse: %s",err.Error())
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return 
	}
	fmt.Fprintf(w,"Response %s ",content)
}