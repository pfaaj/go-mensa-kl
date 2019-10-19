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
	"strings"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/text"
)

var plan = mensa.GetMensaPlan()
var i = 0
var j = 0
var cat = 0
var date = ""
var showNextDay = false

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
				if err := t.Write(fmt.Sprintf("Category: %s\n\nMeal: %s\n\n",
					plan.AllMeals[j].Categories[cat],
					strings.TrimLeft(plan.AllMeals[j].Meals[i], " "))); err != nil {
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

func main() {

	t, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())
	borderless, err := text.New()
	if err != nil {
		panic(err)
	}
	if err := borderless.Write("Text without border."); err != nil {
		panic(err)
	}

	unicode, err := text.New()
	if err != nil {
		panic(err)
	}
	if err := unicode.Write("你好，世界!"); err != nil {
		panic(err)
	}

	trimmed, err := text.New()
	if err != nil {
		panic(err)
	}
	if err := trimmed.Write("Trims lines that don't fit onto the canvas because they are too long for its width.."); err != nil {
		panic(err)
	}

	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}
	if err := wrapped.Write("Buffet: "+plan.Buffet, text.WriteCellOpts(cell.FgColor(cell.ColorRGB24(124, 252, 0)))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("\n\nBuffet theme: "+plan.BuffetDescription, text.WriteCellOpts(cell.FgColor(cell.ColorRGB24(124, 252, 0)))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("\n\nBuffet prices: "+plan.BuffetPrices, text.WriteCellOpts(cell.FgColor(cell.ColorRGB24(124, 252, 0)))); err != nil {
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
		container.BorderTitle("PRESS Q TO QUIT"),
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
										container.BorderTitle("你好，世界!"),
										container.PlaceWidget(unicode),
									),
								),
							),
							container.Bottom(
								container.Border(linestyle.Light),
								container.BorderTitle("Trims lines"),
								container.PlaceWidget(trimmed),
							),
						),
					),
					container.Bottom(
						container.Border(linestyle.Light),
						container.BorderTitle("Wraps lines at rune boundaries"),
						container.PlaceWidget(wrapped),
					),
				),
			),
			container.Right(
				container.Border(linestyle.Light),
				container.BorderTitle("Your Mensa Plan"),
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
