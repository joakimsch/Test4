package main

import (
	"html/template"
	"net/http"
	"io/ioutil"
	"log"
	"path"
	"encoding/xml"
	"fmt"
	"strconv"
	"os"
)

type Temps struct {
	XMLName xml.Name `xml:"temperature"`
	Value   string   `xml:"value,attr"`
	Unit    string   `xml:"unit,attr"`
}
type Ws struct {
	Mps  string `xml:"mps,attr"`
	Name string `xml:"name,attr"`
}
type Symbol struct {
	XMLName xml.Name `xml:"symbol"`
	Name string `xml:"name,attr"`
}
type Times struct {
	XMLName     xml.Name `xml:"time"`
	From        string   `xml:"from,attr"`
	To          string   `xml:"to,attr"`
	Period      string   `xml:"period,attr"`
	Temperature Temps    `xml:"temperature"`
	WindSpeed   Ws       `xml:"windSpeed"`
	Symbol      Symbol   `xml:"symbol"`
}
type Forecasts struct {
	XMLName xml.Name `xml:"forecast"`
	Tabular []Times  `xml:"tabular>time"`
}
type weatherdata struct {
	XMLName  xml.Name  `xml:"weatherdata"`
	Forecast Forecasts `xml:"forecast"`
}

func main()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fp := path.Join("Oblig4/templates/index.html")

		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, "index"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/kristiansand", func(w http.ResponseWriter, r *http.Request) {
		weatherSuggestion(w, r, "Vest-Agder/Kristiansand/Kristiansand")
	})
	http.ListenAndServe(":8080", nil)
}

func weatherSuggestion(w http.ResponseWriter, r *http.Request, extension string) {

	url := "http://www.yr.no/sted/Norge/" + extension + "/varsel.xml"

	var XMLweather = ""
	if getWdFromUrl, err := getXML(url); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		log.Println("Received XML:")
		log.Println(getWdFromUrl)
		XMLweather = getWdFromUrl
	}

	fp := path.Join("OBLIG4/templates/suggestion.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var weatherData weatherdata
	xmlErr := xml.Unmarshal([]byte(XMLweather), &weatherData)
	if xmlErr != nil {
		log.Fatal(xmlErr)
	}
	var weatherType = weatherData.Forecast.Tabular[0].Symbol.Name
	var totalTemp = convertStringToInt(weatherData.Forecast.Tabular[0].Temperature.Value)
	var totalWs = convertStringToFloat(weatherData.Forecast.Tabular[0].WindSpeed.Mps)

	println(weatherType)
	println(totalWs)
	println(totalTemp)
	if err := tmpl.Execute(w, getMessage(totalTemp, totalWs, weatherType)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func getMessage(temp int, ws float64, wt string) (string) {
	var message string
	if wt == "Regn" || wt == "Kraftig regn" || wt == "Kraftige regnbyger" {
		message = "I dag kan du late som du bor i Bergen"
		if temp > 15 {
			message += "PLASKEPARTY!!!!"
		} else if temp > 10 {
			message += "Regnjakke. Alt du trenger"
		} else if temp > 0 {
			message += " Perfekt vær å sove i"
		} else {
			message += " Arnold Schwarzenegger i Batman & Robin. Get it?"
		}
	} else if wt == "Skyet" || wt == "Lettskyet" || wt == "Delvis skyet" {
		message = "Solen er litt shy i dag"
		if temp > 20 {
			message += "Alt annet enn shorts i dag er ulovlig, jfr UD-21"
		} else if temp > 10 {
			message += " Nice temp, do what you want, brother"
		} else if temp > 0 {
			message += "Surt. Bare surt"
		} else {
			message += " You gonna freez boy!"
		}
	} else if wt == "Klarvær" || wt == "Sol" {
		message = "Suns out, guns out!."
		if temp > 20 {
			message += " Varmere enn djevelsens baller i dag."
		} else if temp > 10 {
			message += "Alt er lov i dag"
		} else if temp > 0 {
			message += " Ufyselig kaldt, men du overlever."
		} else {
			message += "Sibirtilstander i dag"
		}
	} else {
		message = ""
		if temp > 20 {
			message += "Spådom: Air Condition er din bestevenn i dag."
		} else if temp > 10 {
			message += "Hvis du er kul tar du på shorts, er du smart tar du på litt mer klær."
		} else if temp > 0 {
			message += "You gonna freez boy!."
		}
		if ws > 10 {
			message += " Det blåser mer enn Stavanger en sommer dag, på med allværsjakke!."
		}
	}
	return message
}
func convertStringToFloat(toFloat string) (float64) {
	f, err := strconv.ParseFloat(toFloat, 64)
	if err == nil {
		fmt.Printf("Type: %T \n", f)

		fmt.Println("Value:", f)
	}
	return f
}
func convertStringToInt(toInt string) (int) {
	convertedInt, err := strconv.Atoi(toInt)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return convertedInt
}
func getXML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Read body: %v", err)
	}

	return string(data), nil
}