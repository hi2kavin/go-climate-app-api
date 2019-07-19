package main
import (
	"bytes"
	"encoding/json"
	"fmt"
	"errors"
	"time"
	"github.com/jung-kurt/gofpdf"
	"strconv"
	"net/http"
	"net/smtp"
	"github.com/domodwyer/mailyak"
	"github.com/jinzhu/gorm"
	"os"
	"net"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

func printClimate(c Climate)[]string{
	var a []string
	
	a=append(a,c.Date)
	a=append(a,c.City)
	a=append(a,strconv.FormatFloat(c.MaxTempC, 'f', 2, 64))
	a=append(a,strconv.FormatFloat(c.MaxTempF, 'f', 2, 64))
	a=append(a,strconv.FormatFloat(c.MinTempC, 'f', 2, 64))
	a=append(a,strconv.FormatFloat(c.MinTempF, 'f', 2, 64))
	a=append(a,strconv.FormatFloat(c.MaxWindK, 'f', 2, 64))
	a=append(a,strconv.FormatFloat(c.MaxWindM, 'f', 2, 64))
	a=append(a,strconv.FormatFloat(c.Humidity, 'f', 2, 64))
	a=append(a,strconv.FormatFloat(c.TotalPrecipM, 'f', 2, 64))
	a=append(a,strconv.FormatFloat(c.TotalPrecipI, 'f', 2, 64))
	a=append(a,c.Text)
	a=append(a,c.Icon)
	a=append(a,strconv.Itoa(c.Code))
	return a
}
func readJSONFromUrl(url string) ([]Climate, error) {
	var climate []Climate 
	resp, err := http.Get(url)
	if err != nil {
		return climate, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	fmt.Println(resp.Body)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &climate); err != nil {
		return climate, err
	}
	fmt.Println(climate)
	return climate, nil
}

func getDataFormat()[]string{
	 s:=[]string{"Date" ,"Name" ,"Maximum Temperature in C" , "Maximum Temperature in F","Minimum Temperature in C","Minimum Temperature in F",
"Maximum Wind in Kph","Maximum Wind in Mph","Humidity","Total Precipitation in mm","Total Precipitation in in",
"Text" ,"Icon","Code"}
    return s
}

func createPdf(city string,name string)string{
	db,err:=gorm.Open("postgres","user=postgres dbname=climate password=root sslmode=disable")
	if err!=nil{
		fmt.Println(err)
	}
	defer db.Close()
	var climate []Climate
	db.Where("city=?",city).Find(&climate)
	if err != nil {
		panic(err)
	}
	var data,cli []string	
	data=getDataFormat()
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.Image("images/logo.png", 50, 100, 130, 0, false, "", 0, "")
	pdf.SetFont("Arial", "B", 50)
	pdf.CellFormat(200,200,city, "0", 1, "BC", false, 0, "linkStr")
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(200, 50,"created for "+name, "0", 1, "BR", false, 0, "linkStr")
	pdf.CellFormat(170, 10,"on "+time.Now().String()[:19], "0", 1, "BR", false, 0, "linkStr")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	for j:=0;j<len(climate);j++{
	cli=printClimate(climate[j])
	pdf.SetFont("Arial", "B", 30)
	pdf.Image("images/water.png", 50, 500, 130, 0, true, "", 0, "")
	pdf.CellFormat(200, 50,cli[1]+" : "+cli[0], "0", 1, "C", false, 0, "linkStr")
	pdf.SetFont("Arial", "B", 12)
	for i:=0;i<14;i++{
		pdf.CellFormat(80, 10, data[i], "1", 0, "L", false, 0, "linkStr")
		pdf.CellFormat(110, 10,cli[i], "1", 1, "L", false, 0, "linkStr")
	}
	pdf.AddPage()
}
    
	err = pdf.OutputFileAndClose("pdf/"+name+".pdf")
	return name
}

func mailOther(w http.ResponseWriter,r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var m SendOther
	err:=decoder.Decode(&m)
	if err!=nil{
		panic(err)
	}
	db,err:=gorm.Open("postgres","user=postgres dbname=climate password=root sslmode=disable")
	if err!=nil{
		fmt.Println(err)
	}
	defer db.Close()
	var s MailLog
	s.From=m.Recievemail
	s.To=m.Sendmail
	s.City=m.City
	s.Date=time.Now().String()
	db.Create(&s)
	emailPdfOther(m.City,m.Sendmail,m.Sendmail,m.Sendsubject,m.Sendmessage)
	fmt.Fprintf(w,"mailsend")
}

func emailPdfOther(city string,email string,name string,subject string,message string){

	name=createPdf(city,name)
	mail := mailyak.New("smtp.gmail.com:587", smtp.PlainAuth("", "fakemail4rj@gmail.com", "beMyLove", "smtp.gmail.com"))
	fmt.Println("sending mail to ",email)
	mail.To(email)
	mail.From("fakemal4rj@gmail.com")
	mail.Subject(subject)
	mail.HTML().Set(message)
	input,err:=os.Open("pdf/"+name+".pdf")
	if err!=nil{
		panic(err)
	}
	mail.Attach("weather Report -"+city+".pdf", input)
	if err := mail.Send(); err != nil {
		panic(" 💣 ")
	}
}
func mailPassword(email string,password string){


	mail := mailyak.New("smtp.gmail.com:587", smtp.PlainAuth("", "fakemail4rj@gmail.com", "beMyLove", "smtp.gmail.com"))
	fmt.Println("sending mail to ",email)
	mail.To(email)
	mail.From("oops@itsallbroken.com")
	mail.Subject("Weather Report")
	mail.HTML().Set("<link rel='stylesheet' href='https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css'"+
	"integrity='sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T' crossorigin='anonymous'>"+
	"<div class='container' style='background:lightblue;height:200px;width:400px'>"+
	"<div class='row' ><h3>Weather.IO</h3></div><div class='row'>Your weather.io password"+
	" is<h6 style='background:darkgrey;color:grey'>"+password+
	"</h6></div><div class='row'><a class='btn btn-primary' href='http://localhost:3000/login'>Click to login</a></div></div>")
	if err := mail.Send(); err != nil {
		panic(" 💣 ")
	}
}
func emailPdf(city string,email string,name string){

	name=createPdf(city,name)
	mail := mailyak.New("smtp.gmail.com:587", smtp.PlainAuth("", "fakemail4rj@gmail.com", "beMyLove", "smtp.gmail.com"))
	fmt.Println("sending mail to ",email)
	mail.To(email)
	mail.From("oops@itsallbroken.com")
	mail.Subject("Weather Report")
	mail.HTML().Set("Have a Nice Day")
	input,err:=os.Open("pdf/"+name+".pdf")
	if err!=nil{
		panic(err)
	}
	mail.Attach("weather Report -"+city+".pdf", input)
	if err := mail.Send(); err != nil {
		panic(" 💣 ")
	}
}
func externalIP() (string, error) {
		ifaces, err := net.Interfaces()
		if err != nil {
			return "", err
		}
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 {
				continue // interface down
			}
			if iface.Flags&net.FlagLoopback != 0 {
				continue // loopback interface
			}
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if ip == nil || ip.IsLoopback() {
					continue
				}
				ip = ip.To4()
				if ip == nil {
					continue // not an ipv4 address
				}
				return ip.String(), nil
			}
		}
		return "", errors.New("are you connected to the network?")
	}