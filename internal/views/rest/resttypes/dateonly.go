package resttypes

import (
	"fmt"
	"strings"
	"time"
)

type DateOnly time.Time

func (d *DateOnly) unmarshalString(str string) error {
	t, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return err
	}

	*d = DateOnly(t)
	return nil
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	return d.unmarshalString(strings.Trim(string(b), "\""))
}

func (d *DateOnly) UnmarshalParam(str string) error {
	return d.unmarshalString(str)
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	value := fmt.Sprintf("\"%s\"", time.Time(d).Format(time.DateOnly))
	return []byte(value), nil
}
