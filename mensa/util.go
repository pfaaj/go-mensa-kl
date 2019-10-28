package mensa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

//LatestRelease contains the url to the current release (change before deploying a new release to github)
const LatestRelease = "https://github.com/pfaaj/go-mensa-kl/releases/download/0.5-beta/go-mensa"

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

//ParseDate converts date string to time.Time
func ParseDate(date string) time.Time {

	tokens := strings.Split(date, ",")

	tokens = strings.Split(tokens[1], "-")

	tokens = strings.Split(tokens[0], ".")

	year, err := strconv.Atoi(strings.TrimSpace(tokens[2]))

	if err != nil {
		panic(err)
	}
	month, err := strconv.Atoi(strings.TrimSpace(tokens[1]))

	if err != nil {
		panic(err)
	}

	day, err := strconv.Atoi(strings.TrimSpace(tokens[0]))

	if err != nil {
		panic(err)
	}

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 651387237, time.UTC)

	return t
}

//IsDateToday tells if date is today
func IsDateToday(date string) bool {
	t := ParseDate(date)
	istoday := false
	current := time.Now()
	if t.Year() == current.Year() && t.Month() == current.Month() && t.Day() == current.Day() {
		istoday = true
	}

	return istoday
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
