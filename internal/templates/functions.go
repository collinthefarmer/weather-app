package templates

import "time"

func AsDateInputValue(t time.Time) string {
	return t.Format(time.DateTime)
}
