package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"net/http"
	"github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

type Weather struct {
	
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Forecast Forecast `json:"forecast"`
}
type Location struct {
	
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}
type Condition struct {
	
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}
type Current struct {
	
	LastUpdatedEpoch int       `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            float64   `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
	WindMph          float64   `json:"wind_mph"`
	WindKph          float64   `json:"wind_kph"`
	WindDegree       int       `json:"wind_degree"`
	WindDir          string    `json:"wind_dir"`
	PressureMb       float64   `json:"pressure_mb"`
	PressureIn       float64   `json:"pressure_in"`
	PrecipMm         float64   `json:"precip_mm"`
	PrecipIn         float64   `json:"precip_in"`
	Humidity         int       `json:"humidity"`
	Cloud            int       `json:"cloud"`
	FeelslikeC       float64   `json:"feelslike_c"`
	FeelslikeF       float64   `json:"feelslike_f"`
	VisKm            float64   `json:"vis_km"`
	VisMiles         float64   `json:"vis_miles"`
	Uv               float64   `json:"uv"`
	GustMph          float64   `json:"gust_mph"`
	GustKph          float64   `json:"gust_kph"`
}
type Day struct {
	
	MaxtempC      float64   `json:"maxtemp_c"`
	MaxtempF      float64   `json:"maxtemp_f"`
	MintempC      float64   `json:"mintemp_c"`
	MintempF      float64   `json:"mintemp_f"`
	AvgtempC      float64   `json:"avgtemp_c"`
	AvgtempF      float64   `json:"avgtemp_f"`
	MaxwindMph    float64   `json:"maxwind_mph"`
	MaxwindKph    float64   `json:"maxwind_kph"`
	TotalprecipMm float64   `json:"totalprecip_mm"`
	TotalprecipIn float64   `json:"totalprecip_in"`
	AvgvisKm      float64   `json:"avgvis_km"`
	AvgvisMiles   float64   `json:"avgvis_miles"`
	Avghumidity   float64   `json:"avghumidity"`
	Condition     Condition `json:"condition"`
	Uv            float64   `json:"uv"`
}
type Astro struct {
	
	Sunrise  string `json:"sunrise"`
	Sunset   string `json:"sunset"`
	Moonrise string `json:"moonrise"`
	Moonset  string `json:"moonset"`
}
type Forecastday struct {
	
	Date      string `json:"date"`
	DateEpoch int    `json:"date_epoch"`
	Day       Day    `json:"day"`
	Astro     Astro  `json:"astro"`
}
type Forecast struct {
	
	Forecastday []Forecastday `json:"forecastday"`
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
	a=append(a,string(c.Code))
	return a
}
func readJSONFromUrl(url string) (Weather, error) {
	var weather Weather 
	resp, err := http.Get(url)
	if err != nil {
		return weather, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &weather); err != nil {
		return weather, err
	}
	return weather, nil
}
func deleteOldDatabase(){
	db, err := gorm.Open("postgres", "user=postgres dbname=climate password=root sslmode=disable")
	if(err!=nil){
		panic(err)
		
	}
	defer db.Close()
	var today string 
	today= time.Now().AddDate(0,0,-5).String()[:10]
	fmt.Println(today)
	db.Where("date < ?", today).Delete(Climate{})
	fmt.Println("daa deleted")
}
func weatherToClimate(weather Weather){
	db, err := gorm.Open("postgres", "user=postgres dbname=climate password=root sslmode=disable")
	if(err!=nil){
		panic(err)
		
	}
	defer db.Close()

	var climate Climate
	climate.Date=weather.Forecast.Forecastday[0].Date
	climate.City=weather.Location.Name
	climate.MaxTempC =weather.Forecast.Forecastday[0].Day.MaxtempC
	climate.MaxTempF =weather.Forecast.Forecastday[0].Day.MaxtempF
	climate.MinTempC =weather.Forecast.Forecastday[0].Day.MintempC
	climate.MinTempF =weather.Forecast.Forecastday[0].Day.MintempF
	climate.MaxWindK =weather.Forecast.Forecastday[0].Day.MaxwindKph
	climate.MaxWindM =weather.Forecast.Forecastday[0].Day.MaxwindMph
	climate.Humidity =weather.Forecast.Forecastday[0].Day.Avghumidity
	climate.TotalPrecipM =weather.Forecast.Forecastday[0].Day.TotalprecipMm
	climate.TotalPrecipI =weather.Forecast.Forecastday[0].Day.TotalprecipIn
	climate.Text =weather.Forecast.Forecastday[0].Day.Condition.Text
	climate.Icon =weather.Forecast.Forecastday[0].Day.Condition.Icon
	climate.Code =weather.Forecast.Forecastday[0].Day.Condition.Code
	db.AutoMigrate(&Climate{})
	db.Create(&climate)
	
}
func getDates()(string){
	t:=time.Now()
	//  var dates []string
	//  for i:=4;i<6;i++{
	//  	dates=append(dates,t.AddDate(0,0,i+1).String()[:10])
	//  }

	return  t.AddDate(0,0,8).String()[:10]
	// return dates
	
}
func fetcher(url string){
	weather, err := readJSONFromUrl(url)
	fmt.Println("fetching: ")
	if err != nil {
		fmt.Println(err)
	}

	if weather.Location.Lat!=0 {
	weatherToClimate(weather)
}
}

func main() {
	city:=[...]string{
	"Kolhapur",
	"Adilabad",
	"Adoni",
	"Amalapuram",
	"Anakapalle",
	"Anantapur",
	"Banganapalle",
	"Bapatla",
	"Bhadrachalam",
	"Bhainsa",
	"Bheemunipatnam",
	"Bhimavaram",
	"Bhongir",
	"Bobbili",
	"Bodhan",
	"Chirala",
	"Chittoor",
	"Cuddapah",
	"Dharmavaram",
	"Eluru",
	"Gadwal",
	"Gooty",
	"Gudivada",
	"Gudur",
	"Guntur",
	"Hindupur",
	"Hyderabad",
	"Ichchapuram",
	"Jagtial",
	"Jammalamadugu",
	"Jangaon",
	"Kadapa",
	"Kadiri",
	"Kagaznagar",
	"Kakinada",
	"Kandukur",
	"Karimnagar",
	"Kavali",
	"Khammam",
	"Koratla",
	"Kothagudem",
	"Kovvur",
	"Kurnool",
	"Macherla",
	"Machilipatnam",
	"Madanapalle",
	"Mahbubnagar",
	"Mancherial",
	"Mandapeta",
	"Markapur",
	"Medak",
	"Nagari",
	"Nandyal",
	"Narasapur",
	"Narasaraopet",
	"Narayanpet",
	"Narsipatnam",
	"Nellore",
	"Nidadavole",
	"Nirmal",
	"Nizamabad",
	"Nuzvid",
	"Ongole",
	"Palacole",
	"Pedana",
	"Peddapuram",
	"Pithapuram",
	"Pondur",
	"Proddatur",
	"Punganur",
	"Puttur",
	"Rajahmundry",
	"Rajam",
	"Ramachandrapuram",
	"Ramagundam",
	"Rayachoti",
	"Renigunta",
	"Repalle",
	"Sadasivpet",
	"Salur",
	"Samalkot",
	"Sattenapalle",
	"Siddipet",
	"Singapur",
	"Srikakulam",
	"Suryapet",
	"Tadepalligudem",
	"Tadpatri",
	"Tandur",
	"Tanuku",
	"Tenali",
	"Tirupati",
	"Tuni",
	"Uravakonda",
	"Venkatagiri",
	"Vicarabad",
	"Vijayawada",
	"Vinukonda",
	"Visakhapatnam",
	"Vizianagaram",
	"Wanaparthy",
	"Warangal",
	"Yellandu",
	"Yerraguntla",
	"Zahirabad",
	"Rajampet",
	"Along",
	"Bomdila",
	"Itanagar",
	"Pasighat",
	"Abhayapuri",
	"Amguri",
	"Barpeta",
	"Bilasipara",
	"Bongaigaon",
	"Dhekiajuli",
	"Dhubri",
	"Dibrugarh",
	"Digboi",
	"Diphu",
	"Dispur",
	"Gauripur",
	"Goalpara",
	"Golaghat",
	"Guwahati",
	"Haflong",
	"Hailakandi",
	"Hojai",
	"Jorhat",
	"Karimganj",
	"Kokrajhar",
	"Lanka",
	"Lumding",
	"Mankachar",
	"Margherita",
	"Mariani",
	"Nagaon",
	"Nalbari",
	"Rangia",
	"Sibsagar",
	"Silchar",
	"Tezpur",
	"Tinsukia",
	"Amarpur",
	"Araria",
	"Arrah",
	"Asarganj",
	"Aurangabad",
	"Bagaha",
	"Bahadurganj",
	"Bakhtiarpur",
	"Banka",
	"Barauli",
	"Barh",
	"Begusarai",
	"Bettiah",
	"Bhabua",
	"Bhagalpur",
	"Bikramganj",
	"BodhGaya",
	"Buxar",
	"Chhapra",
	"Colgong",
	"Dalsinghsarai",
	"Darbhanga",
	"Daudnagar",
	"Dehri-on-Sone",
	"Dhaka",
	"Dighwara",
	"Dumraon",
	"Fatwah",
	"Forbesganj",
	"Gaya",
	"Gopalganj",
	"Hajipur",
	"Hilsa",
	"Hisua",
	"Islampur",
	"Jagdispur",
	"Jamalpur",
	"Jamui",
	"Jehanabad",
	"Jhanjharpur",
	"Kanti",
	"Katihar",
	"Khagaria",
	"Kharagpur",
	"Kishanganj",
	"Lakhisarai",
	"Lalganj",
	"Madhubani",
	"Maharajganj",
	"Makhdumpur",
	"Maner",
	"Manihari",
	"Marhaura",
	"Masaurhi",
	"Mirganj",
	"Mokameh",
	"Motihari",
	"Motipur",
	"Munger",
	"Murliganj",
	"Muzaffarpur",
	"Naugachhia",
	"Nawada",
	"Nokha",
	"Patna",
	"Piro",
	"Purnia",
	"Rafiganj",
	"Rajgir",
	"Ramnagar",
	"Revelganj",
	"Saharsa",
	"Samastipur",
	"Sasaram",
	"Sheikhpura",
	"Sheohar",
	"Sherghati",
	"Silao",
	"Sitamarhi",
	"Siwan",
	"Sonepur",
	"Sugauli",
	"Sultanganj",
	"Supaul",
	"Akaltara",
	"Ambikapur",
	"Arang",
	"Balod",
	"Bhatapara",
	"Bilaspur",
	"Birgaon",
	"Champa",
	"Dalli-Rajhara",
	"Dhamtari",
	"Dipka",
	"Dongargarh",
	"Durg-BhilaiNagar",
	"Jagdalpur",
	"Janjgir",
	"Jashpurnagar",
	"Kanker",
	"Kawardha",
	"Kondagaon",
	"Korba",
	"Mahasamund",
	"Mahendragarh",
	"Mungeli",
	"Raigarh",
	"Raipur",
	"Sakti",
	"Amli",
	"Silvassa",
	"Asola",
	"Delhi",
	"Aldona",
	"Madgaon",
	"Mapusa",
	"Margao",
	"Marmagao",
	"Panaji",
	"Ahmedabad",
	"Amreli",
	"Anand",
	"Ankleshwar",
	"Bharuch",
	"Bhavnagar",
	"Bhuj",
	"Cambay",
	"Dahod",
	"Deesa",
	"Dharampur",
	"Dholka",
	"Gandhinagar",
	"Godhra",
	"Himatnagar",
	"Idar",
	"Jamnagar",
	"Junagadh",
	"Kadi",
	"Kalavad",
	"Kalol",
	"Kapadvanj",
	"Karjan",
	"Keshod",
	"Khambhalia",
	"Khambhat",
	"Kheda",
	"Kheralu",
	"Kodinar",
	"Lathi",
	"Limbdi",
	"Mahesana",
	"Mahuva",
	"Manavadar",
	"Mandvi",
	"Mangrol",
	"Mansa",
	"Modasa",
	"Morvi",
	"Nadiad",
	"Navsari",
	"Padra",
	"Palanpur",
	"Palitana",
	"Pardi",
	"Patan",
	"Petlad",
	"Porbandar",
	"Radhanpur",
	"Rajkot",
	"Rajpipla",
	"Rajula",
	"Ranavav",
	"Rapar",
	"Salaya",
	"Sanand",
	"Sidhpur",
	"Sihor",
	"Songadh",
	"Surat",
	"Talaja",
	"Thangadh",
	"Tharad",
	"Umbergaon",
	"Umreth",
	"Una",
	"Unjha",
	"Upleta",
	"Vadnagar",
	"Vadodara",
	"Valsad",
	"Vapi",
	"Veraval",
	"Vijapur",
	"Viramgam",
	"Visnagar",
	"Vyara",
	"Wadhwan",
	"Wankaner",
	"Adalaj",
	"Adityana",
	"Alang",
	"Andada",
	"Anjar",
	"Atul",
	"Ambala",
	"Assandh",
	"Ateli",
	"Bahadurgarh",
	"Barwala",
	"Bhiwani",
	"Ellenabad2",
	"Faridabad",
	"Fatehabad",
	"Gharaunda",
	"Gohana",
	"Gurgaon",
	"Haibat(YamunaNagar)",
	"Hansi",
	"Hisar",
	"Hodal",
	"Jhajjar",
	"Jind",
	"Kaithal",
	"KalanWali",
	"Kalka",
	"Karnal",
	"Ladwa",
	"Narnaul",
	"Narwana",
	"Palwal",
	"Panipat",
	"Pehowa",
	"Pinjore",
	"Rania",
	"Ratia",
	"Rewari",
	"Rohtak",
	"Safidon",
	"Samalkha",
	"Shahbad",
	"Sirsa",
	"Sohna",
	"Sonipat",
	"Thanesar",
	"Tohana",
	"Yamunanagar",
	"Arki",
	"Baddi",
	"Bilaspur",
	"Chamba",
	"Dalhousie",
	"Dharamsala",
	"Hamirpur",
	"Mandi",
	"Nahan",
	"Shimla",
	"Solan",
	"Sundarnagar",
	"Jammu",
	"Anantnag",
	"Arnia",
	"Awantipora",
	"Baramula",
	"Kathua",
	"Leh",
	"Punch",
	"Rajauri",
	"Sopore",
	"Srinagar",
	"Udhampur",
	"Ara",
	"Chaibasa",
	"Chakradharpur",
	"Chatra",
	"Churi",
	"Daltonganj",
	"Deoghar",
	"Dhanbad",
	"Dumka",
	"Garhwa",
	"Ghatshila",
	"Giridih",
	"Godda",
	"Gomoh",
	"Gumia",
	"Gumla",
	"Hazaribag",
	"Hussainabad",
	"Jamshedpur",
	"Jamtara",
	"Khunti",
	"Lohardaga",
	"Madhupur",
	"Mihijam",
	"Musabani",
	"Pakaur",
	"Patratu",
	"Ranchi",
	"Sahibganj",
	"Saunda",
	"Simdega",
	"TenuDam-cum-Kathhara",
	"Bangalore",
	"Belgaum",
	"Bellary",
	"Chamrajnagar",
	"Chintamani",
	"Chitradurga",
	"Gulbarga",
	"Gundlupet",
	"Hassan",
	"Hospet",
	"Hubli",
	"Karkala",
	"Karwar",
	"Kolar",
	"Kota",
	"Lakshmeshwar",
	"Lingsugur",
	"Maddur",
	"Madhugiri",
	"Madikeri",
	"Magadi",
	"Mahalingpur",
	"Malavalli",
	"Malur",
	"Mandya",
	"Mangalore",
	"Manvi",
	"Mudbidri",
	"Muddebihal",
	"Mudhol",
	"Mulbagal",
	"Mundargi",
	"Mysore",
	"Nanjangud",
	"Pavagada",
	"Puttur",
	"Raichur",
	"Ramanagaram",
	"Ranibennur",
	"RobertsonPet",
	"Ron",
	"Sadalgi",
	"Sagar",
	"Sakleshpur",
	"Sandur",
	"Sankeshwar",
	"Saundatti-Yellamma",
	"Savanur",
	"Sedam",
	"Shahabad",
	"Shahpur",
	"Shiggaon",
	"Shimoga",
	"Shorapur",
	"Shrirangapattana",
	"Sidlaghatta",
	"Sindgi",
	"Sindhnur",
	"Sira",
	"Sirsi",
	"Siruguppa",
	"Srinivaspur",
	"Talikota",
	"Tarikere",
	"Tekkalakota",
	"Terdal",
	"Tiptur",
	"Tumkur",
	"Udupi",
	"Vijayapura",
	"Wadi",
	"Yadgir",
	"Alappuzha",
	"Aroor",
	"Attingal",
	"Chengannur",
	"Chittur-Thathamangalam",
	"Erattupetta",
	"Irinjalakuda",
	"Kalpetta",
	"Kanhangad",
	"Kannur",
	"Kasaragod",
	"Kayamkulam",
	"Kochi",
	"Kodungallur",
	"Kollam",
	"Kottayam",
	"Kozhikode",
	"Kunnamkulam",
	"Malappuram",
	"Nedumangad",
	"Neyyattinkara",
	"Ottappalam",
	"Palai",
	"Palakkad",
	"Pappinisseri",
	"Paravoor",
	"Pathanamthitta",
	"Payyannur",
	"Ponnani",
	"Punalur",
	"Shoranur",
	"Taliparamba",
	"Thiruvananthapuram",
	"Thrissur",
	"Tirur",
	"Vaikom",
	"Varkala",
	"Kavaratti",
	"AshokNagar",
	"Balaghat",
	"Betul",
	"Bhopal",
	"Burhanpur",
	"Chhatarpur",
	"Dabra",
	"Datia",
	"Dewas",
	"Dhar",
	"Fatehabad",
	"Gwalior",
	"Indore",
	"Itarsi",
	"Jabalpur",
	"Katni",
	"Kotma",
	"Lahar",
	"Lundi",
	"Maharajpur",
	"Mahidpur",
	"Maihar",
	"Manasa",
	"Manawar",
	"Mandla",
	"Mandsaur",
	"Mauganj",
	"Morena",
	"Multai",
	"Murwara",
	"Nagda",
	"Nainpur",
	"Narsinghgarh",
	"Neemuch",
	"Niwari",
	"Nowgong",
	"Pali",
	"Panagar",
	"Pandhurna",
	"Panna",
	"Pasan",
	"Porsa",
	"Prithvipur",
	"Raghogarh-Vijaypur",
	"Rahatgarh",
	"Raisen",
	"Rajgarh",
	"Ratlam",
	"Rau",
	"Rehli",
	"Rewa",
	"Sabalgarh",
	"Sagar",
	"Sanawad",
	"Sarangpur",
	"Sarni",
	"Satna",
	"Sausar",
	"Sehore",
	"Sendhwa",
	"Seoni",
	"Seoni-Malwa",
	"Shahdol",
	"Shajapur",
	"Shamgarh",
	"Sheopur",
	"Shivpuri",
	"Shujalpur",
	"Sidhi",
	"Sihora",
	"Singrauli",
	"Sironj",
	"Sohagpur",
	"Tarana",
	"Tikamgarh",
	"Ujhani",
	"Ujjain",
	"Umaria",
	"Vidisha",
	"WaraSeoni",
	"Ahmednagar",
	"Akola",
	"Amravati",
	"Aurangabad",
	"Baramati",
	"Chalisgaon",
	"Chinchani",
	"Devgarh",
	"Dhule",
	"Durgapur",
	"Ichalkaranji",
	"Jalna",
	"Kalyan",
	"Latur",
	"Loha",
	"Lonar",
	"Lonavla",
	"Mahad",
	"Mahuli",
	"Malegaon",
	"Malkapur",
	"Manchar",
	"Manjlegaon",
	"Manmad",
	"Manwath",
	"Mehkar",
	"Miraj",
	"Mul",
	"Mumbai",
	"Nagpur",
	"Nanded-Waghala",
	"Nandgaon",
	"Nandura",
	"Nandurbar",
	"Narkhed",
	"Nawapur",
	"Nilanga",
	"Osmanabad",
	"Ozar",
	"Pachora",
	"Paithan",
	"Palghar",
	"Pandharpur",
	"Panvel",
	"Parbhani",
	"Parli",
	"Parola",
	"Partur",
	"Pathardi",
	"Pathri",
	"Patur",
	"Pauni",
	"Pen",
	"Phaltan",
	"Pulgaon",
	"Pune",
	"Purna",
	"Pusad",
	"Rahuri",
	"Rajura",
	"Ramtek",
	"Ratnagiri",
	"Raver",
	"Risod",
	"Sailu",
	"Sangamner",
	"Sangli",
	"Sangole",
	"Sasvad",
	"Satana",
	"Satara",
	"Savner",
	"Shegaon",
	"Shirdi",
	"Shirpur-Warwade",
	"Shrigonda",
	"Shrirampur",
	"Sillod",
	"Sinnar",
	"Solapur",
	"Talode",
	"Tasgaon",
	"Tuljapur",
	"Tumsar",
	"Uran",
	"Wai",
	"Wani",
	"Wardha",
	"Warora",
	"Warud",
	"Washim",
	"Yevla",
	"Udgir",
	"Umarga",
	"Umarkhed",
	"Umred",
	"Vaijapur",
	"Vasai",
	"Virar",
	"Vita",
	"Yavatmal",
	"Yawal",
	"Imphal",
	"Kakching",
	"Lilong",
	"Thoubal",
	"Jowai",
	"Nongstoin",
	"Shillong",
	"Tura",
	"Aizawl",
	"Lunglei",
	"Saiha",
	"Dimapur",
	"Kohima",
	"Mokokchung",
	"Tuensang",
	"Wokha",
	"Zunheboto",
	"Anandapur",
	"Anugul",
	"Asika",
	"Balangir",
	"Balasore",
	"Baleshwar",
	"Bamra",
	"Bargarh",
	"Baripada",
	"Basudebpur",
	"Bhadrak",
	"Bhawanipatna",
	"Bhuban",
	"Bhubaneswar",
	"Biramitrapur",
	"Brahmapur",
	"Brajrajnagar",
	"Cuttack",
	"Dhenkanal",
	"Gunupur",
	"Hinjilicut",
	"Jagatsinghapur",
	"Jajapur",
	"Jaleswar",
	"Jatani",
	"Jharsuguda",
	"Joda",
	"Kantabanji",
	"Karanjia",
	"Kendrapara",
	"Khordha",
	"Koraput",
	"Malkangiri",
	"Paradip",
	"Parlakhemundi",
	"Pattamundai",
	"Phulabani",
	"Puri",
	"Rairangpur",
	"Rajagangapur",
	"Raurkela",
	"Rayagada",
	"Sambalpur",
	"Soro",
	"Sundargarh",
	"Talcher",
	"Titlagarh",
	"Karaikal",
	"Mahe",
	"Pondicherry",
	"Yanam",
	"Amritsar",
	"Barnala",
	"Batala",
	"Budhlada",
	"Chandigarh",
	"Dasua",
	"Dhuri",
	"Dinanagar",
	"Faridkot",
	"Fazilka",
	"Firozpur",
	"Giddarbaha",
	"Gobindgarh",
	"Gurdaspur",
	"Hoshiarpur",
	"Jagraon",
	"Jalalabad",
	"Jalandhar",
	"Jandiala",
	"Kapurthala",
	"Kartarpur",
	"Khanna",
	"Kharar",
	"Kurali",
	"Longowal",
	"Ludhiana",
	"Mansa",
	"Maur",
	"Moga",
	"Mohali",
	"Morinda",
	"Mukerian",
	"Muktsar",
	"Nabha",
	"Nakodar",
	"Nangal",
	"Nawanshahr",
	"Pathankot",
	"Patiala",
	"Patran",
	"Patti",
	"Phagwara",
	"Phillaur",
	"Qadian",
	"Raikot",
	"Rajpura",
	"Rupnagar",
	"Samana",
	"Sangrur",
	"Sujanpur",
	"Sunam",
	"Talwara",
	"Zira",
	"Bali",
	"Banswara",
	"Ajmer",
	"Alwar",
	"Bandikui",
	"Baran",
	"Barmer",
	"Bikaner",
	"Fatehpur",
	"Jaipur",
	"Jaisalmer",
	"Jodhpur",
	"Kota",
	"Lachhmangarh",
	"Ladnu",
	"Lakheri",
	"Lalsot",
	"Losal",
	"Makrana",
	"Malpura",
	"Mandalgarh",
	"Mandawa",
	"Mangrol",
	"Nadbai",
	"Nagar",
	"Nagaur",
	"Nargund",
	"Nasirabad",
	"Nathdwara",
	"Navalgund",
	"Nawalgarh",
	"Neem-Ka-Thana",
	"Nelamangala",
	"Nimbahera",
	"Nipani",
	"Nohar",
	"Nokha",
	"Pali",
	"Phalodi",
	"Phulera",
	"Pilani",
	"Pilibanga",
	"Pindwara",
	"Pratapgarh",
	"Raisinghnagar",
	"Rajakhera",
	"Rajaldesar",
	"Rajgarh(Alwar)",
	"Rajgarh(Churu",
	"Rajsamand",
	"Ratangarh",
	"Rawatbhata",
	"Rawatsar",
	"Reengus",
	"Sadri",
	"Sagwara",
	"Sambhar",
	"Sangaria",
	"Shahpura",
	"Sheoganj",
	"Sikar",
	"Sirohi",
	"Sojat",
	"Sujangarh",
	"Sumerpur",
	"Suratgarh",
	"Taranagar",
	"Tonk",
	"Udaipur",
	"Gangtok",
	"Calcutta",
	"Arakkonam",
	"Arcot",
	"Aruppukkottai",
	"Bhavani",
	"Chengalpattu",
	"Chennai",
	"Coimbatore",
	"Coonoor",
	"Cuddalore",
	"Dharmapuri",
	"Dindigul",
	"Erode",
	"Gudalur",
	"Kanchipuram",
	"Karaikudi",
	"Karur",
	"Lalgudi",
	"Madurai",
	"Nagapattinam",
	"Nagercoil",
	"Namagiripettai",
	"Namakkal",
	"Nandivaram-Guduvancheri",
	"Natham",
	"Neyveli",
	"Padmanabhapuram",
	"Palani",
	"Palladam",
	"Pallapatti",
	"Pallikonda",
	"Panruti",
	"Pattukkottai",
	"Perambalur",
	"Peravurani",
	"Periyakulam",	
	"Pollachi",
	"Polur",
	"Ponneri",
	"Pudukkottai",	
	"Puliyankudi",
	"Rajapalayam",
	"Ramanathapuram",
	"Rasipuram",
	"Salem",
	"Sankarankoil",
	"Sankari",
	"Sathyamangalam",
	"Sattur",
	"Shenkottai",
	"Sholavandan",
	"Sholingur",
	"Sirkali",
	"Sivaganga",
	"Sivagiri",
	"Sivakasi",
	"Srivilliputhur",
	"Surandai",
	"Suriyampalayam",
	"Tenkasi",
	"Thammampatti",
	"Thanjavur",
	"Tharamangalam",
	"Tharangambadi",
	"TheniAllinagaram",
	"Thirumangalam",
	"Thirunindravur",
	"Thiruparappu",
	"Thirupuvanam",
	"Thiruthuraipoondi",
	"Thiruvallur",
	"Thiruvarur",
	"Thoothukudi",
	"Thuraiyur",
	"Tindivanam",
	"Tiruchendur",
	"Tiruchengode",
	"Tiruchirappalli",
	"Tirukalukundram",
	"Tirukkoyilur",
	"Tirunelveli",
	"Tirupathur",
	"Tiruppur",
	"Tiruttani",
	"Tiruvannamalai",
	"Tiruvethipuram",
	"Tittakudi",
	"Udhagamandalam",
	"Udumalaipettai",
	"Unnamalaikadai",
	"Usilampatti",
	"Uthamapalayam",
	"Uthiramerur",
	"Vadakkuvalliyur",
	"Vadalur",
	"Vadipatti",
	"Valparai",
	"Vandavasi",
	"Vaniyambadi",
	"Vedaranyam",
	"Vellakoil",
	"Vellore",
	"Vikramasingapuram",
	"Viluppuram",
	"Virudhachalam",
	"Virudhunagar",
	"Viswanatham",
	"Agartala",
	"Badharghat",
	"Dharmanagar",
	"Indranagar",
	"Jogendranagar",
	"Kailasahar",
	"Khowai",
	"Pratapgarh",
	"Udaipur",
	"Achhnera",
	"Adari",
	"Agra",
	"Aligarh",
	"Allahabad",
	"Amroha",
	"Azamgarh",
	"Bahraich",
	"Ballia",
	"Balrampur",
	"Banda",
	"Bareilly",
	"Chandausi",
	"Dadri",
	"Deoria",
	"Etawah",
	"Fatehabad",
	"Fatehpur",
	"GreaterNoida",
	"Hamirpur",
	"Hardoi",
	"Jajmau",
	"Jaunpur",
	"Jhansi",
	"Kalpi",
	"Kanpur",
	"Kota",
	"Laharpur",
	"Lakhimpur",
	"LalGopalganjNindaura",
	"Lalganj",
	"Lalitpur",
	"Lar",
	"Loni",
	"Lucknow",
	"Mathura",
	"Meerut",
	"Modinagar",
	"Muradnagar",
	"Nagina",
	"Najibabad",
	"Nakur",
	"Nanpara",
	"Naraura",
	"NaugawanSadat",
	"Nautanwa",
	"Nawabganj",
	"Nehtaur",
	"NOIDA",
	"Noorpur",
	"Obra",
	"Orai",
	"Padrauna",
	"PaliaKalan",
	"Parasi",
	"Phulpur",
	"Pihani",
	"Pilibhit",
	"Pilkhuwa",
	"Powayan",
	"Pukhrayan",
	"Puranpur",
	"Purquazi",
	"Purwa",
	"RaeBareli",
	"Rampur",
	"RampurManiharan",
	"Rasra",
	"Rath",
	"Renukoot",
	"Reoti",
	"Robertsganj",
	"Rudauli",
	"Rudrapur",
	"Sadabad",
	"Safipur",
	"Saharanpur",
	"Sahaspur",
	"Sahaswan",
	"Sahawar",
	
	"Saidpur",
	"Sambhal",
	
	"Samthar",
	"Sandi",
	"Sandila",
	"Sardhana",
	"Seohara",
	"Shahabad",
	"Shahabad",
	"Shahganj",
	"Shahjahanpur",
	"Shamli",
	"Shamsabad",
	"Shamsabad",
	"Sherkot",
	"Shikarpur",
	"Shikohabad",
	"Shishgarh",
	"Siana",
	"Sikanderpur",
	"Sikandrabad",
	"Sirsaganj",
	"Sirsi",
	"Sitapur",
	"Soron",
	"Suar",
	"Sultanpur",
	"Sumerpur",
	"Tanda",
	"Thakurdwara",
	"Tilhar",
	"Tulsipur",
	"Tundla",
	"Unnao",
	"Utraula",
	"Varanasi",
	"Vrindavan",
	"Zaidpur",
	"Zamania",
	"Almora",
	"Bazpur",
	"Chamba",
	"Dehradun",
	"Haldwani",
	"Haridwar",
	"Jaspur",
	"Kashipur",
	"kichha",
	"Kotdwara",
	"Manglaur",
	"Mussoorie",
	"Nagla",
	
	"Pauri",
	"Pithoragarh",
	"Ramnagar",
	"Rishikesh",
	"Roorkee",
	"Rudrapur",
	"Sitarganj",
	"Tehri",
	"Muzaffarnagar",
	"Adra",
	"Arambagh",
	"Asansol",
	"Baharampur",
	"Bally",
	"Balurghat",
	"Bankura",
	"Barakar",
	"Barasat",
	"Bardhaman",
	"Chinsura",
	"Contai",
	"Darjeeling",
	"Durgapur",
	"Haldia",
	"Howrah",
	"Islampur",
	"Jhargram",
	"Kharagpur",
	"Kolkata",
	"Mainaguri",
	"Mal",
	"Mathabhanga",
	"Medinipur",
	"Memari",
	"Monoharpur",
	"Murshidabad",
	"Nabadwip",
	"Naihati",
	"Panchla",
	"Pandua",
	"Purulia",
	"Raghunathpur",
	"Raiganj",
	"Ranaghat",
	"Sainthia",
	"Santipur",
	"Siliguri",
	"Sonamukhi",
	"Suri",
	"Taki",
	"Tamluk",
	"Tarakeswar",
	"Chikmagalur",
	"Dharwad",
	"Gadag",
}
deleteOldDatabase()
var dates string
dates=getDates()
for i := 0;  i<len(city); i++ {
	//for j:=0;j<len(dates);j++{
	fmt.Println(city[i],dates)
	url := "http://api.apixu.com/v1/forecast.json?key=41dba5f315b740348dd111839190807&q="+city[i]+"&dt="+dates
	go fetcher(url)
//	}
	
	
}
}



