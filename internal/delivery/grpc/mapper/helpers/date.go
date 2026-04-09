package helpers

import (
	"time"

	"google.golang.org/genproto/googleapis/type/date"
)

func TimeToPbDate(t time.Time) *date.Date {
	utc := t.UTC()

	return &date.Date{
		Year:  int32(utc.Year()),
		Month: int32(utc.Month()),
		Day:   int32(utc.Day()),
	}
}

func TimePtrToPbDate(t *time.Time) *date.Date {
	if t == nil {
		return nil
	}

	utc := t.UTC()

	return &date.Date{
		Year:  int32(utc.Year()),
		Month: int32(utc.Month()),
		Day:   int32(utc.Day()),
	}
}

func PbDateToTime(d *date.Date) time.Time {
	if d == nil {
		return time.Time{}
	}

	return time.Date(
		int(d.Year),
		time.Month(d.Month),
		int(d.Day),
		0, 0, 0, 0,
		time.UTC,
	)
}

func PbDateToTimePtr(d *date.Date) *time.Time {
	if d == nil {
		return nil
	}

	t := time.Date(
		int(d.Year),
		time.Month(d.Month),
		int(d.Day),
		0, 0, 0, 0,
		time.UTC,
	)

	return &t
}
