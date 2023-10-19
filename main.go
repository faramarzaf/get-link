package main

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	Q_720  = "720p.x265"
	Q_1080 = "1080p.x265.10bit"
)

func main() {

	fmt.Println("Welcome to get link app.")
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter the name of series")
	scanner.Scan()
	inputName := scanner.Text()
	seriesName := getSeriesName(inputName)

	var urls []string
	var baseUrl string

	for {

		fmt.Println("Enter season: ")
		scanner.Scan()
		season, _ := strconv.Atoi(scanner.Text())
		formattedSeason := fmt.Sprintf("%02d", season)

		fmt.Println("Enter quality: (720 / 1080)")
		scanner.Scan()
		quality := scanner.Text()

		baseUrl = getBaseUrlByQuality(seriesName, quality, formattedSeason)

		c := colly.NewCollector()

		c.OnHTML(".list tr", func(e *colly.HTMLElement) {
			url := e.ChildAttr("a", "href")
			if strings.HasSuffix(url, ".mkv") {
				urls = append(urls, baseUrl+url+"\n")
			}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.Request, "\nError:", err)
			os.Exit(-1)
		})

		c.OnScraped(func(r *colly.Response) {
		})

		extensions.RandomUserAgent(c)
		c.Visit(baseUrl)

		writeToFile("data"+strconv.Itoa(season)+".txt", urls)
		urls = nil

		fmt.Println("Continue? [y/n]")
		scanner.Scan()
		ans := scanner.Text()

		if strings.ToLower(ans) == "n" {
			os.Exit(0)
		}
	}

}

func getSeriesName(inputName string) string {
	hasSpace := strings.Contains(inputName, " ")

	if hasSpace {
		return strings.Replace(inputName, " ", ".", -1)
	}
	return inputName
}

func getBaseUrlByQuality(seriesName, quality, season string) string {
	if quality == "720" {
		return "https://dls2.top-movies2filmha.click/DonyayeSerial/series/" + seriesName + "/Soft.Sub/S" + season + "/" + Q_720 + ".BluRay/"

	} else if quality == "1080" {
		return "https://dls2.top-movies2filmha.click/DonyayeSerial/series/" + seriesName + "/Soft.Sub/S" + season + "/" + Q_1080 + ".BluRay/"
	}

	return ""
}

func writeToFile(fileName string, data []string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, value := range data {
		_, err2 := f.WriteString(value)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	fmt.Println("done")
}
