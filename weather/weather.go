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

	info = fmt.Sprintf("Temperature: %f %s \n\n",
		w.Main.Temp, w.Weather[0].Description)

	info += fmt.Sprintf("Max temp: %f\n\nMin temp: %f\n\n",
		w.Main.TempMax, w.Main.TempMin)

	return
}

//Info returns current weather data of Kaiserslautern
func Info() *openweathermap.CurrentWeatherData {
	w, err := owm.NewCurrent("C", "en", key) // fahrenheit (imperial) with Russian output
	if err != nil {
		log.Fatalln(err)
	}

	w.CurrentByName("Kaiserslautern")
	return w
}
