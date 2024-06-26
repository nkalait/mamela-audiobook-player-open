package merror

import "github.com/sqweek/dialog"

func PanicError(e error) {
	if e != nil {
		panic(e)
	}
}

func ShowError(msg string, e error) {
	if e != nil {
		go func() {
			dialog.Message("%s:\n%s", msg, e.Error()).Title("Error").Error()
		}()
	}
}
