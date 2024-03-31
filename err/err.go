package err

import (
	"github.com/sqweek/dialog"
	// "fyne.io/fyne/v2/dialog"
)

func PanicError(e error) {
	if e != nil {
		panic(e)
	}
}

func ShowError(s string) {
	go func() {
		dialog.Message("%s", s).Title("Error").Error()
	}()
}
