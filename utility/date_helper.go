package utility

import "time"

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
)

func InAPIDateFormat(date time.Time) string {
	return date.Format(apiDateLayout)
}
