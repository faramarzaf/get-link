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
	BASE_URL = "https://dls2.top-movies2filmha.click/DonyayeSerial/series"
)

func main() {

	allSeries := getAllSeries()

	fmt.Println("Welcome to get link app.")
	scanner := bufio.NewScanner(os.Stdin)
	for _, seriesName := range allSeries {
		seasonsCount := getSeasonsCount(seriesName)
		fmt.Printf("name: %s , seasons: %v\n", seriesName, seasonsCount)
	}
	var urls []string
	var downloadUrl string

	for {
		fmt.Println("Enter name: ")
		scanner.Scan()
		seriesName := scanner.Text()

		fmt.Println("Enter season: ")
		scanner.Scan()
		season, _ := strconv.Atoi(scanner.Text())
		formattedSeason := fmt.Sprintf("%02d", season)

		printQualityBySeriesNameAndSeasonNumber(seriesName, formattedSeason)
		fmt.Println("Enter quality:")
		scanner.Scan()
		quality := scanner.Text()

		downloadUrl = getDownloadUrlByQualityAndSeasonNumber(seriesName, quality, formattedSeason)

		c := colly.NewCollector()
		c.OnHTML(".list tr", func(e *colly.HTMLElement) {
			url := e.ChildAttr("a", "href")
			if strings.HasSuffix(url, ".mkv") {
				urls = append(urls, downloadUrl+url+"\n")
			}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.Request, "\nError:", err)
		})

		c.OnScraped(func(r *colly.Response) {
		})

		extensions.RandomUserAgent(c)
		c.Visit(downloadUrl)

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

func getAllSeries() []string {
	var allSeries []string
	c := colly.NewCollector()
	c.OnHTML(".list tr", func(e *colly.HTMLElement) {
		series := e.ChildAttr("a", "href")
		formattedName := strings.Replace(series, "/", "", -1)
		if formattedName != "" && formattedName != ".." { // this is the parent directory column in table of website
			allSeries = append(allSeries, formattedName)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.Request, "\nError:", err)
		os.Exit(-1)
	})

	c.OnScraped(func(r *colly.Response) {
	})

	extensions.RandomUserAgent(c)
	c.Visit(BASE_URL)

	return allSeries
}

func getSeasonsCount(seriesName string) (seasons int) {
	count := 0
	seriesHomePageUrl := BASE_URL + "/" + seriesName + "/Soft.Sub/"

	//var failedFixedUrlRecords string

	c := colly.NewCollector()
	c.OnHTML(".list tr", func(e *colly.HTMLElement) {
		season := e.ChildAttr("a", "href")
		if strings.HasPrefix(season, "S") {
			count++
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.Request, "\nError:", err)
		//	failedFixedUrlRecords = BASE_URL + "/" + seriesName + "/"
	})

	c.OnScraped(func(r *colly.Response) {
	})

	extensions.RandomUserAgent(c)
	c.Visit(seriesHomePageUrl)

	return count
}

func getDownloadUrlByQualityAndSeasonNumber(seriesName, quality, season string) string {
	return BASE_URL + "/" + seriesName + "/Soft.Sub/S" + season + "/" + quality + "/"
}

func printQualityBySeriesNameAndSeasonNumber(name, formattedSeason string) {
	var qualities []string
	c := colly.NewCollector()
	c.OnHTML(".list tr", func(e *colly.HTMLElement) {
		quality := e.ChildAttr("a", "href")
		formattedQuality := strings.Replace(quality, "/", "", -1)
		qualities = append(qualities, formattedQuality+"\n")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.Request, "\nError:", err)
	})

	c.OnScraped(func(r *colly.Response) {
	})

	extensions.RandomUserAgent(c)
	c.Visit(BASE_URL + "/" + name + "/" + "Soft.Sub/" + "S" + formattedSeason)
	fmt.Println(qualities)
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
