package mensa

import (
	"fmt"

	"os"

	"strings"

	"github.com/gocolly/colly"

	"encoding/json"

	"time"

	"io/ioutil"
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
	OpeningTimes      string
}

const agent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0"

func filterNonempty(ss []string, clean bool) (ret []string) {
	for _, s := range ss {
		noempty := strings.Replace(s, " ", "", -1)
		if len(noempty) >= 1 {
			if clean == true {
				ret = append(ret, noempty)
			} else {
				ret = append(ret, s)
			}

		}
	}
	return
}

func parseDate() {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
}

func insertRune(a []rune, x rune, i int) {
	a = append(a[:i], append([]rune{x}, a[i:]...)...)
}

func writeInfo() {

	info := CrawlInfo{}

	info.CrawledAt = time.Now()

	file, _ := json.MarshalIndent(info, "", " ")

	_ = ioutil.WriteFile("info.json", file, 0644)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readInfo() CrawlInfo {

	dat, err := ioutil.ReadFile("info.json")
	check(err)
	str := string(dat)
	res := CrawlInfo{}
	json.Unmarshal([]byte(str), &res)

	return res
}

func getTime(str string) {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
}

//GetMensaPlan returns a mensa plan
func GetMensaPlan() (plans Plans) {

	res := readInfo()

	isoYear, isoWeek := time.Now().ISOWeek()
	storedYear, storedWeek := res.CrawledAt.ISOWeek()

	if isoYear > storedYear || (isoYear == storedYear && isoWeek > storedWeek) {
		//purge cache
		os.RemoveAll("./cache")
		//store new time of cache creation
		writeInfo()
	}

	c := colly.NewCollector(
		colly.CacheDir("./cache"),
	)

	// Find and meal links
	c.OnHTML("div[class=dailyplan]", func(e *colly.HTMLElement) {
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
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.studierendenwerk-kaiserslautern.de/kaiserslautern/essen-und-trinken/tu-kaiserslautern/mensa/")

	return

}
