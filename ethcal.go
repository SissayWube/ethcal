// Package ethcal converts dates between the Ethiopian calendar (EC) and the
// Gregorian calendar (GC) using Julian Day Numbers.
//
// The Ethiopian calendar has 13 months: twelve of 30 days and a final month,
// Pagume, of 5 days (6 in a leap year). A leap year is any year where
// year % 4 == 3. The new year (1 Meskerem) falls on 11 September in the
// Gregorian calendar, or 12 September in the year before a Gregorian leap year.
package ethcal

import (
	"errors"
	"fmt"
	"time"
)

// Julian Day Number offset of the Ethiopian (Amete Mihret) epoch.
const jdnEpoch = 1723856

// Julian Day Number of the Unix epoch (1970-01-01).
const jdnUnixEpoch = 2440588

// Errors returned by this package.
var (
	ErrInvalidMonth = errors.New("ethcal: month must be between 1 and 13")
	ErrInvalidDay   = errors.New("ethcal: day is out of range for the given month")
	ErrInvalidYear  = errors.New("ethcal: year must be 1 or greater")
)

// Weekday represents an Ethiopian day of the week.
type Weekday int

const (
	Segno    Weekday = iota // ሰኞ — Monday
	Maksegno                // ማክሰኞ — Tuesday
	Rob                     // ረቡዕ — Wednesday
	Hamus                   // ሐሙስ — Thursday
	Arb                     // ዓርብ — Friday
	Kidame                  // ቅዳሜ — Saturday
	Ihud                    // እሁድ — Sunday
)

var weekdayNamesEnglish = [7]string{
	"Segno", "Maksegno", "Rob", "Hamus", "Arb", "Kidame", "Ihud",
}

var weekdayNamesAmharic = [7]string{
	"ሰኞ", "ማክሰኞ", "ረቡዕ", "ሐሙስ", "ዓርብ", "ቅዳሜ", "እሁድ",
}

// String returns the English (transliterated) name of the weekday.
func (w Weekday) String() string {
	if w < Segno || w > Ihud {
		return ""
	}
	return weekdayNamesEnglish[w]
}

// Amharic returns the Amharic name of the weekday.
func (w Weekday) Amharic() string {
	if w < Segno || w > Ihud {
		return ""
	}
	return weekdayNamesAmharic[w]
}

// Date is a date in the Ethiopian calendar.
type Date struct {
	Year  int
	Month int // 1 (Meskerem) .. 13 (Pagume)
	Day   int
}

var monthsEnglish = [13]string{
	"Meskerem", "Tikimt", "Hidar", "Tahsas", "Tir", "Yekatit", "Megabit",
	"Miazia", "Genbot", "Sene", "Hamle", "Nehase", "Pagume",
}

var monthsAmharic = [13]string{
	"መስከረም", "ጥቅምት", "ኅዳር", "ታኅሣሥ", "ጥር", "የካቲት", "መጋቢት",
	"ሚያዝያ", "ግንቦት", "ሰኔ", "ሐምሌ", "ነሐሴ", "ጳጉሜን",
}

// IsLeapYear reports whether the Ethiopian year is a leap year.
func IsLeapYear(year int) bool {
	return mod(year, 4) == 3
}

// DaysInMonth returns the number of days in the given Ethiopian month.
func DaysInMonth(year, month int) (int, error) {
	switch {
	case month < 1 || month > 13:
		return 0, ErrInvalidMonth
	case month < 13:
		return 30, nil
	case IsLeapYear(year):
		return 6, nil
	default:
		return 5, nil
	}
}

// MonthName returns the English (transliterated) name of an Ethiopian month.
func MonthName(month int) string {
	if month < 1 || month > 13 {
		return ""
	}
	return monthsEnglish[month-1]
}

// MonthNameAmharic returns the Amharic name of an Ethiopian month.
func MonthNameAmharic(month int) string {
	if month < 1 || month > 13 {
		return ""
	}
	return monthsAmharic[month-1]
}

// New validates and returns an Ethiopian Date.
func New(year, month, day int) (Date, error) {
	d := Date{year, month, day}
	return d, d.Validate()
}

// FromTime converts a time.Time to an Ethiopian Date. Only the calendar
// date in t's own location is used.
func FromTime(t time.Time) Date {
	y, m, d := t.Date()
	days := floorDiv(time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix(), 86400)
	return jdnToEthiopic(int(days) + jdnUnixEpoch)
}

// FromYMD converts a Gregorian year, month, day to an Ethiopian Date.
// It returns an error if the Gregorian date does not exist (e.g. 30 Feb).
func FromYMD(year int, month time.Month, day int) (Date, error) {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || t.Month() != month || t.Day() != day {
		return Date{}, errors.New("ethcal: invalid Gregorian date")
	}
	return FromTime(t), nil
}

// Today returns the current Ethiopian date in the local time zone.
func Today() Date { return FromTime(time.Now()) }

// Validate checks that the date exists in the Ethiopian calendar.
func (d Date) Validate() error {
	if d.Year < 1 {
		return ErrInvalidYear
	}
	n, err := DaysInMonth(d.Year, d.Month)
	if err != nil {
		return err
	}
	if d.Day < 1 || d.Day > n {
		return ErrInvalidDay
	}
	return nil
}

// Equal reports whether d and other represent the same date.
func (d Date) Equal(other Date) bool {
	return d.Year == other.Year && d.Month == other.Month && d.Day == other.Day
}

// String formats the date as "YYYY-MM-DD".
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// FormatLong returns a long-form date such as "1 Meskerem 2000".
func (d Date) FormatLong() string {
	return fmt.Sprintf("%d %s %d", d.Day, MonthName(d.Month), d.Year)
}

// FormatAmharic returns a long-form date using the Amharic month name.
func (d Date) FormatAmharic() string {
	return fmt.Sprintf("%s %d, %d", MonthNameAmharic(d.Month), d.Day, d.Year)
}

// Weekday returns the Ethiopian day of the week for this date.
func (d Date) Weekday() Weekday {
	// JDN mod 7 gives: 0=Monday, 1=Tuesday, ..., 6=Sunday,
	// which aligns with the Weekday constants.
	jdn := ethiopicToJDN(d.Year, d.Month, d.Day)
	return Weekday(mod(jdn, 7))
}

// Gregorian converts the Ethiopian date to a time.Time (midnight UTC).
func (d Date) Gregorian() (time.Time, error) {
	if err := d.Validate(); err != nil {
		return time.Time{}, err
	}
	jdn := ethiopicToJDN(d.Year, d.Month, d.Day)
	return time.Unix(int64(jdn-jdnUnixEpoch)*86400, 0).UTC(), nil
}

// GregorianYMD converts the Ethiopian date to Gregorian year, month, day.
func (d Date) GregorianYMD() (int, time.Month, int, error) {
	t, err := d.Gregorian()
	if err != nil {
		return 0, 0, 0, err
	}
	y, m, dd := t.Date()
	return y, m, dd, nil
}

// AddDays returns the Ethiopian date n days after (or before, if negative) d.
func (d Date) AddDays(n int) (Date, error) {
	if err := d.Validate(); err != nil {
		return Date{}, err
	}
	return jdnToEthiopic(ethiopicToJDN(d.Year, d.Month, d.Day) + n), nil
}

// AddMonths returns the Ethiopian date n months after (or before) d.
// If the resulting month has fewer days than d.Day, the day is clamped to
// the last valid day of that month.
func (d Date) AddMonths(n int) (Date, error) {
	if err := d.Validate(); err != nil {
		return Date{}, err
	}
	totalMonths := (d.Year-1)*13 + (d.Month - 1) + n
	y := floorDiv(totalMonths, 13) + 1
	m := mod(totalMonths, 13) + 1
	if y < 1 {
		return Date{}, ErrInvalidYear
	}
	maxDay, _ := DaysInMonth(y, m)
	day := d.Day
	if day > maxDay {
		day = maxDay
	}
	return Date{y, m, day}, nil
}

// AddYears returns the Ethiopian date n years after (or before) d.
// If d is Pagume 6 in a leap year and the target year is not a leap year,
// the day is clamped to Pagume 5.
func (d Date) AddYears(n int) (Date, error) {
	if err := d.Validate(); err != nil {
		return Date{}, err
	}
	y := d.Year + n
	if y < 1 {
		return Date{}, ErrInvalidYear
	}
	day := d.Day
	if d.Month == 13 {
		maxDay, _ := DaysInMonth(y, 13)
		if day > maxDay {
			day = maxDay
		}
	}
	return Date{y, d.Month, day}, nil
}

// ---- internal helpers ----

func ethiopicToJDN(year, month, day int) int {
	return jdnEpoch + 365 + 365*(year-1) + floorDiv(year, 4) + 30*month + day - 31
}

func jdnToEthiopic(jdn int) Date {
	r := mod(jdn-jdnEpoch, 1461)
	n := r%365 + 365*(r/1460)
	return Date{
		Year:  4*floorDiv(jdn-jdnEpoch, 1461) + r/365 - r/1460,
		Month: n/30 + 1,
		Day:   n%30 + 1,
	}
}

func floorDiv[T int | int64](a, b T) T {
	q := a / b
	if (a%b != 0) && ((a < 0) != (b < 0)) {
		q--
	}
	return q
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

