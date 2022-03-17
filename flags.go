package main

import (
	"fmt"
	"time"
)

type DateValue struct {
	Date *time.Time
}

func (v *DateValue) String() string {
	if v.Date != nil {
		y, m, _ := v.Date.Date()
		return fmt.Sprintf("%s %d", time.Month(m), y)
	}
	return ""
}

func (v *DateValue) Set(s string) error {
	if t, err := time.Parse("Jan 2006", s); err != nil {
		return err
	} else {
		v.Date = &t
	}
	return nil
}
