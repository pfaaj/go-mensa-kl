package mensa

import (
	"fmt"

	"os"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"

	"time"
)

//CrawlInfo stores info related to the scraping
type CrawlInfo struct {
	CrawledAt time.Time
}

// Plan stores information about a mensa plan for a day
type Plan struct {
	Meals      []string
	Categories []string
	Date       string
}

// Plans stores information about a mensa plan for a week with buffet
type Plans struct {
	Buffet            string
	BuffetDescription string
	BuffetPrices      string
	AllMeals          []Plan
	AtriumMeals       []Plan
	OpeningTimes      string
}

const agent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0"

// ChildTexts returns the stripped text content of all the matching
// element's attributes.
func ChildTexts(e *colly.HTMLElement, goquerySelector string) []string {
	var res []string
	e.DOM.Find(goquerySelector).Each(func(_ int, s *goquery.Selection) {

		res = append(res, strings.TrimSpace(s.Text()))
	})
	return res
}

func getMeals(e *colly.HTMLElement, plans *Plans) {
	category := e.ChildText("div[class=c10l]")
	title := e.ChildText("div[class=c90r]")
	var plan Plan
	plan.Meals = filterNonempty(strings.Split(title, "\n"), false)
	for i, meal := range plan.Meals {
		meal = strings.Replace(meal, "Studenten", "\nStudenten", -1)
		plan.Meals[i] = meal
	}
	plan.Categories = filterNonempty(strings.Split(category, "\n"), true)
	plan.Date = e.ChildText("h5")
	plans.AllMeals = append(plans.AllMeals, plan)
}

func getAtriumMeals(e *colly.HTMLElement, plans *Plans) {

	var plan Plan

	plan.Meals = ChildTexts(e, "div[class=c90r]")

	plan.Date = e.ChildText("h5")

	for i := 0; i < len(plan.Meals); i++ {
		plan.Categories = append(plan.Categories, "Atrium")
	}

	for i, meal := range plan.Meals {
		meal = strings.Replace(meal, "Studenten", "\nStudenten", -1)
		plan.Meals[i] = meal
	}

	plans.AtriumMeals = append(plans.AtriumMeals, plan)

}

//GetMensaPlan returns a mensa plan
func GetMensaPlan() (plans Plans) {

	res := CrawlInfo{}

	if _, err := os.Stat("info.json"); err == nil {
		res = readInfo()

	} else if os.IsNotExist(err) {
		writeInfo()
		res = readInfo()
	}

	isoYear, isoWeek := time.Now().ISOWeek()
	storedYear, storedWeek := res.CrawledAt.ISOWeek()

	weekday := time.Now().Weekday()

	if isoYear > storedYear || (isoYear == storedYear && isoWeek > storedWeek) ||
		(isoWeek == storedWeek && weekday.String() == "Saturday") {
		//purge cache
		os.RemoveAll("./cache")
		//store new time of cache creation
		writeInfo()
	}

	c := colly.NewCollector(
		colly.CacheDir("./cache"),
	)

	atrium := colly.NewCollector(
		colly.CacheDir("./cache"),
	)

	// get all meals
	c.OnHTML("div[class=dailyplan]", func(e *colly.HTMLElement) {
		getMeals(e, &plans)
	})

	// get atriums meals
	atrium.OnHTML("div[class=dailyplan]", func(e *colly.HTMLElement) {
		getAtriumMeals(e, &plans)
	})

	c.OnHTML("div[class=buffet]", func(e *colly.HTMLElement) {
		dishes := e.ChildText("span[class]")

		plans.Buffet = dishes
		if plans.BuffetDescription == "" {
			plans.BuffetDescription = e.ChildText("h5")
			plans.BuffetDescription = strings.Replace(plans.BuffetDescription, ")", ")\n", -1)

		}
		if plans.BuffetPrices == "" {
			plans.BuffetPrices = e.ChildText("div[class=c40r]")
		}

	})

	c.OnHTML("div[class=widget]", func(e *colly.HTMLElement) {

		if e.ChildText("h5[class=widget_header]") == "Öffnungszeiten" {

			opening := e.ChildText("p[class=widget_list]")

			opening = strings.Replace(opening, ".B", ".\n\nB", -1)
			opening = strings.Replace(opening, "rA", "r\n\nA", -1)
			opening = strings.Replace(opening, "sv", "s v", -1)
			opening = strings.Replace(opening, ":m", ": m", -1)
			plans.OpeningTimes = opening
		}

	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", agent)
		fmt.Println("Visited", r.URL)
	})

	atrium.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", agent)
		fmt.Println("Visited", r.URL)
	})

	c.Visit("https://www.studierendenwerk-kaiserslautern.de/kaiserslautern/essen-und-trinken/tu-kaiserslautern/mensa/")
	atrium.Visit("https://www.studierendenwerk-kaiserslautern.de/kaiserslautern/essen-und-trinken/tu-kaiserslautern/mensaria-atrium/")

	return

}
