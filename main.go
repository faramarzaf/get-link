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
	BASE_URL = "https://dls6.top-movies2filmha.click/DonyayeSerial/series"
)

func main() {
	allSeries := getAllSeries()
	fmt.Println("Welcome to the get link app.")
	fmt.Printf("Total of series: %d\n", len(allSeries))

	scanner := bufio.NewScanner(os.Stdin)
	correctedUrls := printSeasonsBySeriesName(allSeries)

	if len(correctedUrls) != 0 {
		fmt.Println("Retrying for corrected series...")
		for _, url := range correctedUrls {
			seasonCount := retryScrapForCorrectedUrls(url)
			fmt.Println(url, " : ", seasonCount)
		}
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

		writeToFile(seriesName+strconv.Itoa(season)+".txt", urls)
		urls = nil

		fmt.Println("Continue? [y/n]")
		scanner.Scan()
		ans := scanner.Text()

		if strings.ToLower(ans) == "n" {
			os.Exit(0)
		}
	}

}

func printSeasonsBySeriesName(allSeries []string) []string {
	var correctedUrls []string
	for _, seriesName := range allSeries {
		seasonsCount, correctedUrl := getSeasonsCount(seriesName)
		if seasonsCount == 1 {
			fmt.Printf("'%s' has %v season\n", seriesName, seasonsCount)
		} else {
			fmt.Printf("'%s' has %v seasons\n", seriesName, seasonsCount)
		}
		if correctedUrl != "" {
			fmt.Printf("Error occured for seaons of %s. Correct url is %s\n", seriesName, correctedUrl)
			correctedUrls = append(correctedUrls, correctedUrl)
		}
	}
	return correctedUrls
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

func getSeasonsCount(seriesName string) (seasons int, correctedUrl string) {
	count := 0
	seriesHomePageUrl := BASE_URL + "/" + seriesName + "/Soft.Sub/"

	c := colly.NewCollector()
	c.OnHTML(".list tr", func(e *colly.HTMLElement) {
		season := e.ChildAttr("a", "href")
		if strings.HasPrefix(season, "S") {
			count++
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		correctedUrl = BASE_URL + "/" + seriesName + "/"
	})

	c.OnScraped(func(r *colly.Response) {
	})

	extensions.RandomUserAgent(c)
	c.Visit(seriesHomePageUrl)

	return count, correctedUrl
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

		if formattedQuality != ".." && formattedQuality != "Sub" {
			qualities = append(qualities, formattedQuality+"\n")
		}
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

func retryScrapForCorrectedUrls(url string) int {
	count := 0
	c := colly.NewCollector()
	c.OnHTML(".list tr", func(e *colly.HTMLElement) {
		season := e.ChildAttr("a", "href")
		if strings.HasPrefix(season, "S") {
			count++
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.Request, "\nError:", err)
	})

	c.OnScraped(func(r *colly.Response) {
	})

	extensions.RandomUserAgent(c)
	c.Visit(url)

	return count
}
