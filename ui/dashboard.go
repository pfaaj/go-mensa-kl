// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Binary textdemo displays a couple of Text widgets.
// Exist when 'q' is pressed.
package main

import (
	"context"
	"fmt"
	"go-mensa/mensa"
	"go-mensa/weather"
	"strings"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/text"
)

var weatherInfo = weather.Info()
var plan = mensa.GetMensaPlan()
var i = 0
var j = 0
var cat = 0
var date = ""
var showNextDay = false
var language = "de"

func countPommesMeals(plan mensa.Plans) (int, int) {
	total := 0
	pommes := 0
	for _, meals := range plan.AllMeals {
		for _, meal := range meals.Meals {
			if strings.Contains(meal, "Pommes") {
				pommes = pommes + 1
			}
			total = total + 1
		}
	}
	return total, pommes
}

// writeLines writes a line of text to the text widget every delay.
// Exits when the context expires.
func writeLines(ctx context.Context, t *text.Text, delay time.Duration) {
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			if showNextDay == true {
				j = (j + 1) % len(plan.AllMeals)
				i = 0
				cat = 0
				showNextDay = false

				if err := t.Write(fmt.Sprintf("Date: %s\n\n",
					plan.AllMeals[j].Date)); err != nil {
					panic(err)
				}
			}

			if i <= len(plan.AllMeals[j].Meals)-1 {
				category := plan.AllMeals[j].Categories[cat]
				if category == "1" || category == "2" {
					category = "Ausgabe " + category
				}
				if err := t.Write(fmt.Sprintf("%s\n\n", category),
					text.WriteCellOpts(cell.FgColor(cell.ColorRGB24(255, 215, 0)))); err != nil {

					panic(err)
				}
				meal := strings.TrimLeft(plan.AllMeals[j].Meals[i], " ")
				if language == "en" {
					meal = mensa.Translate(meal)
				}

				if err := t.Write(fmt.Sprintf("%s\n\n", meal)); err != nil {
					panic(err)
				}

				i = (i + 1)
				cat = (cat + 1) % len(plan.AllMeals[j].Categories)
			}

		case <-ctx.Done():
			return
		}
	}
}

// playBarChart sets the values for the bar chart and draws it
// Exits when the context expires.
func playBarChart(ctx context.Context, bc *barchart.BarChart) {
	var values []int
	total, pommes := countPommesMeals(plan)
	values = append(values, pommes)
	values = append(values, total-pommes)

	if err := bc.Values(values, total); err != nil {
		panic(err)
	}
}

func main() {

	//start showing the plan for the current day
	for idx := 0; idx < len(plan.AllMeals); idx++ {
		if mensa.IsDateToday(plan.AllMeals[idx].Date) {
			j = idx
			break
		}

	}

	t, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())
	borderless, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}

	bc, err := barchart.New(
		barchart.BarColors([]cell.Color{
			cell.ColorBlue,
			cell.ColorRed,
			cell.ColorYellow,
			cell.ColorBlue,
			cell.ColorGreen,
			cell.ColorRed,
		}),
		barchart.ValueColors([]cell.Color{
			cell.ColorRed,
			cell.ColorYellow,
			cell.ColorBlue,
			cell.ColorGreen,
			cell.ColorRed,
			cell.ColorBlue,
		}),
		barchart.ShowValues(),
		barchart.BarWidth(12),
		barchart.Labels([]string{
			"Pommes",
			"Ke Pommes",
		}),
	)

	playBarChart(ctx, bc)

	if err := borderless.Write("Öffnungszeiten:\n\n" + plan.OpeningTimes); err != nil {
		panic(err)
	}

	unicode, err := text.New()
	if err != nil {
		panic(err)
	}
	if err := unicode.Write(weather.Main(weatherInfo)); err != nil {
		panic(err)
	}

	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}

	buffet := plan.Buffet
	if language == "en" {
		buffet = mensa.Translate(buffet)
	}

	if err := wrapped.Write(buffet, text.WriteCellOpts(cell.FgColor(cell.ColorRGB24(124, 252, 0)))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("\n\nThema der Woche: "+plan.BuffetDescription, text.WriteCellOpts(cell.FgColor(cell.ColorRGB24(124, 252, 0)))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("\n\nPreise: "+plan.BuffetPrices, text.WriteCellOpts(cell.FgColor(cell.ColorRGB24(124, 252, 0)))); err != nil {
		panic(err)
	}

	rolled, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		panic(err)
	}

	if err := rolled.Write(fmt.Sprintf("Date: %s\n\n",
		plan.AllMeals[j].Date)); err != nil {
		panic(err)
	}

	go writeLines(ctx, rolled, 500*time.Millisecond)

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT; N TO SEE NEXT DAY"),
		container.SplitVertical(
			container.Left(
				container.SplitHorizontal(
					container.Top(
						container.SplitHorizontal(
							container.Top(
								container.SplitVertical(
									container.Left(
										container.PlaceWidget(borderless),
									),
									container.Right(
										container.Border(linestyle.Light),
										container.BorderTitle("Wetter"),
										container.PlaceWidget(unicode),
									),
								),
							),
							container.Bottom(
								container.Border(linestyle.Light),
								container.BorderTitle("Statistik der Gerichte der Woche."),
								container.PlaceWidget(bc),
							),
						),
					),
					container.Bottom(
						container.Border(linestyle.Light),
						container.BorderTitle("Buffet"),
						container.PlaceWidget(wrapped),
					),
				),
			),
			container.Right(
				container.Border(linestyle.Light),
				container.BorderTitle("Dein Mensa Plan"),
				container.PlaceWidget(rolled),
			),
		),
	)
	if err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		} else if k.Key == 'n' || k.Key == 'N' {
			rolled.Reset()
			showNextDay = true
		}

	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter)); err != nil {
		panic(err)
	}
}
