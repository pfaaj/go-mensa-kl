package mensa

import (
	"fmt"

	"strings"

	"github.com/gocolly/colly"
)

// MensaPlan stores information about a mensa plan
type MensaPlan struct {
	Meals      []string
	Categories []string
	Date       string
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

//GetMensaPlan returns a mensa plan
func GetMensaPlan() (plans []MensaPlan) {
	c := colly.NewCollector()
    
	// Find and meal links
	c.OnHTML("div[class=dailyplan]", func(e *colly.HTMLElement) {
		category := e.ChildText("div[class=c10l]")
		title := e.ChildText("div[class=c90r]")
		var plan MensaPlan
		plan.Meals = filterNonempty(strings.Split(title, "\n"), false)
		fmt.Printf("Meals: %q \n", plan.Meals)
		plan.Categories = filterNonempty(strings.Split(category, "\n"), true)
		fmt.Printf("Categories: %q\n", plan.Categories)
		plan.Date = e.ChildText("h5")
		fmt.Println("Date: " + plan.Date)
		plans = append(plans, plan)
	})

	c.OnHTML("div[class=buffet]", func(e *colly.HTMLElement) {
		category := "Buffet"
		dishes := e.ChildText("span[class]")
	    //plan.Meals = append(plan.Meals, dishes)
		fmt.Printf("Dishes: %q ", dishes)
		fmt.Println("Category: " + category)
		fmt.Println("Description: " + e.ChildText("h5"))
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", agent)
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.studierendenwerk-kaiserslautern.de/kaiserslautern/essen-und-trinken/tu-kaiserslautern/mensa/")


	return

}
