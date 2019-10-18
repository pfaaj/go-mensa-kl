package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

// MensaPlan stores information about a mensa plan
type MensaPlan struct {
	Title    string
	Category string
	Prices   string
	Date     string
}

const agent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0"

func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("div[class=dailyplan]", func(e *colly.HTMLElement) {
		category := e.ChildText("div[class=c10l]")
		title := e.ChildText("div[class=c90r]")
		fmt.Println("Title: " + title)
		fmt.Println("Category: " + category)
		fmt.Println("Date: " + e.ChildText("h5"))
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", agent)
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.studierendenwerk-kaiserslautern.de/kaiserslautern/essen-und-trinken/tu-kaiserslautern/mensa/")
}
