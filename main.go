package main

import (
	"log"

	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var baseURL = "https://tympanus.net/codrops/collective/"
var dropBaseURL = "https://tympanus.net/codrops/collective/collective-%v"

func main() {

	createResultsFile()
	latestFetched := getLatestFetchedDrop()
	fmt.Println("Latest drop downloaded:", latestFetched)

	doc, err := goquery.NewDocument(baseURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	latestTitle := doc.Find("h3").First().Text()
	issue := strings.Replace(latestTitle, "Collective #", "", -1)
	totalDrops, _ := strconv.ParseInt(strings.TrimSpace(issue), 10, 64)

	for i := int(totalDrops); i > 0; i-- {
	//for i := 130; i > 0; i-- {
		if i <= latestFetched {
			break
		}
		fmt.Printf("Getting collective #%v\n", i)
		getDrop(dropBaseURL, i)
	}
	fmt.Println("Up to date!")
}
func getDrop(url string, n int) {

	doc, err := goquery.NewDocument(fmt.Sprintf(url, n))
	if err != nil {
		log.Fatal(err)
		return
	}

	articles := []*Article{}

	doc.Find(".ct-coll-container .ct-coll-item").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a").First().Attr("href")
		a := &Article{
			Title:       s.Find("h2").First().Text(),
			Description: s.Find("p").First().Text(),
			Link:        link,
		}

		articles = append(articles, a)
	})

	writeArticlesToFile(n, articles)
}

func writeArticlesToFile(n int, articles []*Article) {
	file, err := os.OpenFile("./results.md", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	checkErr(err)
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("# Collective #%v\n", n))
	checkErr(err)

	for _, a := range articles {
		if a.Link != "" && a.Title != "" {
			_, err = file.WriteString(fmt.Sprintf("- [%v](%v): %v\n", a.Title, a.Link, a.Description))
		}
	}

	_, err = file.WriteString(fmt.Sprintf("\n"))
	checkErr(err)

	// save changes
	//err = file.Sync()
	//checkErr(err)
}

func createResultsFile() {
	_, err := os.Stat("./results.md")

	if os.IsNotExist(err) {
		file, err := os.Create("./results.md")
		checkErr(err)
		defer file.Close()
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getLatestFetchedDrop() int {
	var latest int
	file, err := os.Open("./results.md")
	checkErr(err)

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "Collective #") {
			issue := strings.Replace(s, "# Collective #", "", -1)
			totalDrops, _ := strconv.ParseInt(strings.TrimSpace(issue), 10, 64)
			latest = int(totalDrops)
			break
		}
	}

	return latest
}

type Article struct {
	Title       string
	Description string
	Link        string
}
