package weather

import (
	"fmt"
	"log"

	"github.com/briandowns/openweathermap"
	owm "github.com/briandowns/openweathermap"
)

var key = "affc427170d00667a5a5381ac0fc8e70"

//Main returns temp and humidity info
func Main(w *openweathermap.CurrentWeatherData) (info string) {

	info = fmt.Sprintf("Kaiserslautern\n%s\n%d °C\n\n",
		w.Weather[0].Description, int(w.Main.Temp))

	info += fmt.Sprintf("Max temp: %d °C\n\nMin temp: %d °C\n\n",
		int(w.Main.TempMax), int(w.Main.TempMin))

	return
}

//Info returns current weather data of Kaiserslautern
func Info() *openweathermap.CurrentWeatherData {
	w, err := owm.NewCurrent("C", "de", key) // fahrenheit (imperial) with Russian output
	if err != nil {
		log.Fatalln(err)
	}

	w.CurrentByName("Kaiserslautern")
	return w
}
