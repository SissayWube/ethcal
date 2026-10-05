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
	"strconv"
	"strings"
	"time"
)

// Julian Day Number offset of the Ethiopian (Amete Mihret) epoch.
const jdnEpoch = 1723856

// Julian Day Number of the Unix epoch (1970-01-01).
const jdnUnixEpoch = 2440588

// Errors returned by this package.
var (
	ErrInvalidMonth      = errors.New("ethcal: month must be between 1 and 13")
	ErrInvalidDay        = errors.New("ethcal: day is out of range for the given month")
	ErrInvalidYear       = errors.New("ethcal: year must be 1 or greater")
	ErrInvalidDateFormat = errors.New("ethcal: date must be in YYYY-MM-DD format")
	ErrInvalidGregorian  = errors.New("ethcal: invalid Gregorian date")
)

// Month constants for the Ethiopian calendar.
const (
	Meskerem = 1
	Tikimt   = 2
	Hidar    = 3
	Tahsas   = 4
	Tir      = 5
	Yekatit  = 6
	Megabit  = 7
	Miazia   = 8
	Genbot   = 9
	Sene     = 10
	Hamle    = 11
	Nehase   = 12
	Pagume   = 13
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

// TimeWeekday converts an Ethiopian Weekday to time.Weekday.
func TimeWeekday(w Weekday) time.Weekday {
	switch w {
	case Ihud:
		return time.Sunday
	default:
		return time.Weekday(w + 1)
	}
}

// WeekdayFromTime converts time.Weekday to an Ethiopian Weekday.
func WeekdayFromTime(w time.Weekday) Weekday {
	switch w {
	case time.Sunday:
		return Ihud
	default:
		return Weekday(w - 1)
	}
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
	if year < 1 {
		return 0, ErrInvalidYear
	}
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

// NewET validates and returns an Ethiopian Date.
func NewET(year, month, day int) (Date, error) {
	d := Date{year, month, day}
	return d, d.Validate()
}

// New is an alias for NewET.
func New(year, month, day int) (Date, error) {
	return NewET(year, month, day)
}

// NewGC converts a Gregorian calendar date to its Ethiopian equivalent.
// It returns an error if the Gregorian date does not exist or precedes
// the Ethiopian epoch (1 Meskerem 1 EC / 27 August 8 AD GC).
func NewGC(year int, month time.Month, day int) (Date, error) {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || t.Month() != month || t.Day() != day {
		return Date{}, ErrInvalidGregorian
	}
	d := FromTime(t)
	if err := d.Validate(); err != nil {
		return Date{}, err
	}
	return d, nil
}

// FromYMD is an alias for NewGC.
func FromYMD(year int, month time.Month, day int) (Date, error) {
	return NewGC(year, month, day)
}

// FromTime converts a time.Time to an Ethiopian Date using the calendar date
// in t's local time zone.
func FromTime(t time.Time) Date {
	y, m, d := t.Date()
	days := floorDiv(time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix(), 86400)
	return jdnToEthiopic(int(days) + jdnUnixEpoch)
}

// Parse parses a date in "YYYY-MM-DD" format.
func Parse(s string) (Date, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 3 {
		return Date{}, ErrInvalidDateFormat
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return Date{}, ErrInvalidDateFormat
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return Date{}, ErrInvalidDateFormat
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return Date{}, ErrInvalidDateFormat
	}
	return NewET(year, month, day)
}

func Today() Date { return FromTime(time.Now()) }

func TodayUTC() Date { return FromTime(time.Now().UTC()) }

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

func (d Date) IsZero() bool {
	return d == Date{}
}

func (d Date) Equal(other Date) bool {
	return d.Year == other.Year && d.Month == other.Month && d.Day == other.Day
}

func (d Date) Compare(other Date) int {
	if d.Year != other.Year {
		if d.Year < other.Year {
			return -1
 		}
		return 1
	}
	if d.Month != other.Month {
		if d.Month < other.Month {
			return -1
		}
		return 1
	}
	if d.Day != other.Day {
		if d.Day < other.Day {
			return -1
		}
		return 1
	}
	return 0
}

func (d Date) Before(other Date) bool {
	return d.Compare(other) < 0
}

func (d Date) After(other Date) bool {
	return d.Compare(other) > 0
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (d Date) MarshalText() ([]byte, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}
	return []byte(d.String()), nil
}

func (d *Date) UnmarshalText(text []byte) error {
	parsed, err := Parse(string(text))
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

func (d Date) FormatLong() string {
	return fmt.Sprintf("%d %s %d", d.Day, MonthName(d.Month), d.Year)
}

func (d Date) FormatAmharic() string {
	return fmt.Sprintf("%s %d, %d", MonthNameAmharic(d.Month), d.Day, d.Year)
}

func (d Date) Weekday() Weekday {
	jdn := ethiopicToJDN(d.Year, d.Month, d.Day)
	return Weekday(mod(jdn, 7))
}

func (d Date) MonthName() string {
	return MonthName(d.Month)
}

func (d Date) MonthNameAmharic() string {
	return MonthNameAmharic(d.Month)
}

func (d Date) DaysInMonth() (int, error) {
	return DaysInMonth(d.Year, d.Month)
}

func (d Date) IsLeapYear() bool {
	return IsLeapYear(d.Year)
}

// Gregorian returns the date as midnight UTC in the Gregorian calendar.
func (d Date) Gregorian() (time.Time, error) {
	if err := d.Validate(); err != nil {
		return time.Time{}, err
	}
	jdn := ethiopicToJDN(d.Year, d.Month, d.Day)
	return time.Unix(int64(jdn-jdnUnixEpoch)*86400, 0).UTC(), nil
}

func (d Date) GC() (time.Time, error) {
	return d.Gregorian()
}

func (d Date) GregorianYMD() (int, time.Month, int, error) {
	t, err := d.Gregorian()
	if err != nil {
		return 0, 0, 0, err
	}
	y, m, dd := t.Date()
	return y, m, dd, nil
}

func (d Date) GCYMD() (int, time.Month, int, error) {
	return d.GregorianYMD()
}

func (d Date) AddDays(n int) (Date, error) {
	if err := d.Validate(); err != nil {
		return Date{}, err
	}
	res := jdnToEthiopic(ethiopicToJDN(d.Year, d.Month, d.Day) + n)
	if res.Year < 1 {
		return Date{}, ErrInvalidYear
	}
	return res, nil
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

// AddYears adds n years to d. If d is Pagume 6 in a leap year and the target year
// is non-leap, the day is clamped to Pagume 5.
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

// Sub returns the number of calendar days between d and other (d - other).
func (d Date) Sub(other Date) (int, error) {
	if err := d.Validate(); err != nil {
		return 0, err
	}
	if err := other.Validate(); err != nil {
		return 0, err
	}
	return ethiopicToJDN(d.Year, d.Month, d.Day) - ethiopicToJDN(other.Year, other.Month, other.Day), nil
}

func ethiopicToJDN(year, month, day int) int {
	return jdnEpoch + 365 + 365*(year-1) + floorDiv(year, 4) + 30*month + day - 31
}

// jdnToEthiopic converts a Julian Day Number to an Ethiopian Date using a 1461-day
// 4-year cycle. Shifting the epoch by 365 days places the cycle's leap day on day 1460,
// allowing r/1460 to absorb it into Pagume of year 3 without special-case branching.
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

