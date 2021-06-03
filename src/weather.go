package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"fmt"
	"flag"
	"os"
)

type openweather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func exportBodyStr(inScan *bufio.Scanner) string {

	var buffer string

	for i:=0 ; inScan.Scan() ; i++ {
		buffer += inScan.Text()
	}

	return buffer
}

func jsonParser(inStr string, key string) string {

	var ret string

	res := openweather{}
	json.Unmarshal([]byte(inStr), &res)

	switch key {
	case "weather":
		ret = res.Weather[0].Main
	case "humidity":
		ret = fmt.Sprintf("%d", res.Main.Humidity)
	case "temperature":
		ret = fmt.Sprintf("%f", res.Main.Temp - 273)
	}

	return ret
}

func main() {

	cityPtr := flag.String("city", "Seoul", "City name")
	countryPtr := flag.String("country", "KR", "Country name")
	apiKey := flag.String("key", "API KEY", "Open Weather API Key")

	flag.Parse()

	if flag.NFlag()  != 3 {
		flag.Usage()
		os.Exit(1)
	}

	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s,%s&appid=%s", *cityPtr, *countryPtr, *apiKey)

	resp, err := http.Get(apiUrl)

	if err != nil{
		panic(err)
	}

	if resp.Status == "200 OK" {
		scanner := bufio.NewScanner(resp.Body)

		var buffer = exportBodyStr(scanner)
		fmt.Printf("\n------ Today Weather ------\n\n")
		fmt.Printf("Weather : %s\n", jsonParser(buffer, "weather"))
		fmt.Printf("Humidity : %s\n", jsonParser(buffer, "humidity"))
		fmt.Printf("Temperature : %s\n", jsonParser(buffer, "temperature"))
		fmt.Printf("---------------------------\n")
	}

}
