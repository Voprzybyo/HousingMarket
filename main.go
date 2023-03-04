package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jasonlvhit/gocron"
)

type FlatData struct {
	Price           string `json:"Price"`
	Area            string `json:"Area"`
	Place           string `json:"Place"`
	PublicationDate string `json:"PublicationDate"`
	FetchDate       string `json:"FetchDate"`
	FetchHour       string `json:"FetchHour"`
	InflationRate   string `json:"InflationRate"`
}

func scrapeOLX(flatData []FlatData, page int) []FlatData {
	
	res, err := http.Get(fmt.Sprintf("https://www.olx.pl/d/nieruchomosci/mieszkania/krakow/?page=%d&search%%5Bfilter_float_price%%3Afrom%%5D=100000&view=list", page))
	CheckError(err)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		ErrorLogger.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckError(err)

	var fd FlatData
	currentTime := time.Now()
	inflationRate := addInflationRate()

	doc.Find(".css-1apmciz").Each(func(i int, s *goquery.Selection) {

		child1 := s.ChildrenFiltered(".css-u2ayx9")
		child2 := s.ChildrenFiltered(".css-odp1qd")

		fd.Price = child1.Find("p").Text()
		fd.Area = child2.Find("span").Text()
		fd.Place = child2.Find("p").Text()
		tempPlace := fd.Place

		// Write formatted data to structure
		fd.Price = fd.FormatPrice(fd.Price)
		fd.Area = fd.FormatSpace(fd.Area)
		fd.Place = fd.FormatPlace(fd.Place)
		fd.PublicationDate = fd.FormatDate(tempPlace)
		fd.FetchDate = currentTime.Format("2006-01-02")
		fd.FetchHour = strconv.Itoa(currentTime.Hour())
		fd.InflationRate = inflationRate

		flatData = append(flatData, fd)
	})
	return flatData
}

func addInflationRate() string {

	res, err := http.Get("https://tradingeconomics.com/poland/inflation-cpi")
	CheckError(err)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		ErrorLogger.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckError(err)

	var inflationRate string
	doc.Find("#ctl00_ContentPlaceHolder1_ctl00_ctl02_Panel1").Each(func(i int, s *goquery.Selection) {
		inflationRate = s.Find("table > tbody > tr > td:nth-child(2)").Text()
	})
	return inflationRate
}

func getAnnouncementNumber() int {

	res, err := http.Get("https://www.olx.pl/nieruchomosci/mieszkania/krakow/?search%5Bfilter_float_price%3Afrom%5D=100000&view=list")
	CheckError(err)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		ErrorLogger.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckError(err)

	var adNumberStr string
	doc.Find(".css-n9feq4").Each(func(i int, s *goquery.Selection) {
		adNumberStr = s.Find("h3 > div").Text()
	})

	splittedString := strings.Split(adNumberStr, " ")
	conv, _ := strconv.ParseInt(splittedString[len(splittedString)-2], 10, 32)
	ret := int(conv)
	return ret
}

func writeToJSON(flatData []FlatData) {

	filename := "Data.json"
	InfoLogger.Println("Writing data to", filename)

	for _, v := range flatData {

		jsonFile, _ := os.Open(filename)

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		data := []FlatData{}
		json.Unmarshal(byteValue, &data)

		data = append(data, v)

		dataBytes, _ := json.MarshalIndent(data, "", "    ")

		ioutil.WriteFile(filename, dataBytes, 0644)
	}
	InfoLogger.Println("Writing to", filename, "finished!")
}

func parseAndWrite() {
	page := 1
	var flatData []FlatData
	var flatDataRet []FlatData
	howManyDataAleradyParsed := 0
	AdNumber := getAnnouncementNumber()

	for {
		InfoLogger.Println("Parsing data from website | Iteration: ", page)
		howManyDataAleradyParsed += len(flatDataRet)
		flatDataRet = scrapeOLX(flatData, page)
		if howManyDataAleradyParsed > AdNumber-60 { // How many ads on one page
			break
		}
		writeToJSON(flatDataRet)
		AddToDb(flatDataRet)
		page++
	}
	InfoLogger.Println("Done! Time to sleep for a while... ")
}

func main() {

	InfoLogger.Println("Starting new scheduler and running task...")
	s := gocron.NewScheduler()
	s.Every(2).Hours().Do(parseAndWrite)
	<-s.Start()

}
