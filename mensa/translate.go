package mensa

import (
	"fmt"

	translate "github.com/dafanasev/go-yandex-translate"
)

//Translate translates text passed in to english
func Translate(str string, language string) (translated string) {

	tr := translate.New("trnsl.1.1.20191023T124920Z.63524b1f3817bdc2.1719c9be2a2e95a9ce652519943ee104fb9e0a56")

	translation, err := tr.Translate(language, str)
	translated = ""
	if err != nil {
		fmt.Println(err)
	} else {
		translated = translation.Result()
	}
	return
}
