# Go-mensa-kl
> A terminal dashboard (and a crawler) for TUK Mensa's meal plan.

![](ui/mensa.png)

## Getting Started

The easiest way to use go-mensa-kl is to download its binary release (only Linux supported)

### Usage example
For retrieving the meal plan in German simply run:
```
./go-mensa
```
or use the -lang argument for any other language code supported by yandex translate api ( still buggy: you might need to minimize and maximize the console window for the translation to show properly), e.g.

```
./go-mensa -lang en # will translate the meal's text to English
```

### Prerequisites

If you want to run the latest code: install all dependencies recursively by using the go get command, go to the main directory of the project and simply run:

```
  go get ./...
```  

to run the dashboard go to ui/ and run
```
  go run dashboard.go
```  



## Built With

* [Colly](https://github.com/gocolly/colly) - Elegant Scraper and Crawler Framework for Golang
* [Termdash](https://github.com/mum4k/termdash) - Terminal based dashboard.

* [openweathermap](https://github.com/briandowns/openweathermap) - Go (golang) package for use with openweathermap.org's API 
* [go-yandex-translate](https://github.com/dafanasev/go-yandex-translate) - Go (golang) Yandex Translate API wrapper
