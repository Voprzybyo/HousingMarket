package main

import (
	"fmt"
	"strings"
	"time"
)

func (flat FlatData) FormatPrice(s string) string {

	if idx := strings.Index(s, " zł"); idx != -1 {
		s = s[:idx]
	}
	return s
}

func (flat FlatData) FormatSpace(s string) string {
	if idx := strings.Index(s, " m"); idx != -1 {
		s = s[:idx]
	}
	return s
}

func (flat FlatData) FormatPlace(s string) string {

	var ss []string
	if strings.Contains(s, ",") {
		ss = strings.Split(s, ", ")
		ss2 := strings.Split(ss[1], " -")
		return ss2[0]
	} else {
		ss = strings.Split(s, " -")
		return ss[0]
	}
}

func ConvertToDate(s string) string {

	if strings.Contains(s, "Dzisiaj") {
		currentTime := time.Now()
		s = fmt.Sprintf("%d-%d-%d",
			currentTime.Year(),
			currentTime.Month(),
			currentTime.Day())
	}

	// Dummy conversion
	if strings.Contains(s, "stycznia") {
		s = strings.Replace(s, " lutego ", "-01-", -1)
	}
	if strings.Contains(s, "lutego") {
		s = strings.Replace(s, " lutego ", "-02-", -1)
	}
	if strings.Contains(s, "marca") {
		s = strings.Replace(s, " marca ", "-03-", -1)
	}
	if strings.Contains(s, "kwietnia") {
		s = strings.Replace(s, " kwietnia ", "-04-", -1)
	}
	if strings.Contains(s, "maja") {
		s = strings.Replace(s, " maja ", "-05-", -1)
	}
	if strings.Contains(s, "czerwca") {
		s = strings.Replace(s, " czerwca ", "-06-", -1)
	}
	if strings.Contains(s, "lipca") {
		s = strings.Replace(s, " lipca ", "-07-", -1)
	}
	if strings.Contains(s, "sierpnia") {
		s = strings.Replace(s, " sierpnia ", "-08-", -1)
	}
	if strings.Contains(s, "września") {
		s = strings.Replace(s, " września ", "-09-", -1)
	}
	if strings.Contains(s, "października") {
		s = strings.Replace(s, " paździenika ", "-10-", -1)
	}
	if strings.Contains(s, "listopada") {
		s = strings.Replace(s, " listopada ", "-11-", -1)
	}
	if strings.Contains(s, "grudnia") {
		s = strings.Replace(s, " grudnia ", "-12-", -1)
	}

	s = strings.Replace(s, "Odświeżono dnia ", "", -1)
	return s
}

func (flat FlatData) FormatDate(s string) string {

	ss := strings.Split(s, " - ")
	ret := ConvertToDate(ss[1])
	return ret

}
