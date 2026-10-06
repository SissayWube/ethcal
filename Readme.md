# ethcal

Convert dates between the Ethiopian calendar (EC) and Gregorian calendar (GC) in Go. Zero dependencies.

## Installation

```bash
go get github.com/SissayWube/ethcal
```

## Quick Start

```go
package main

import (
	"fmt"
	"time"

	"github.com/SissayWube/ethcal"
)

func main() {
	// Ethiopian → Gregorian
	d, _ := ethcal.New(2016, 1, 1)
	t, _ := d.ToGregorian()
	fmt.Println(t.Format("2006-01-02")) // 2023-09-12

	// Gregorian → Ethiopian
	ec := ethcal.FromTime(time.Date(2024, 9, 11, 0, 0, 0, 0, time.UTC))
	fmt.Println(ec.Format())   // 1 Meskerem 2017
	fmt.Println(ec.Amharic())  // መስከረም 1, 2017
	fmt.Println(ec.Weekday())  // Segno

	// Today
	fmt.Println(ethcal.Today())             // local timezone
	fmt.Println(ethcal.TodayIn(ethcal.EAT)) // East Africa Time (UTC+3)
}
```

## API Reference

### Constructors & Conversions (GC → EC)

| Function | Description |
|---|---|
| `New(y, m, d)` | Create and validate an Ethiopian date |
| `FromTime(t)` | Convert `time.Time` → Ethiopian date |
| `FromGregorian(y, m, d)` | Convert Gregorian year, month, day → Ethiopian date |
| `Today()` | Current Ethiopian date in the local timezone |
| `TodayIn(loc)` | Current Ethiopian date in the given timezone (defaults to UTC if `nil`) |

### Conversions (EC → GC)

| Method | Description |
|---|---|
| `d.ToGregorian()` | Ethiopian date → `time.Time` (midnight UTC) |

### JavaScript & API Interoperability

| Function / Method | Description |
|---|---|
| `FromISO(s)` | Parse ISO 8601 string (e.g. JS `Date.toISOString()`) → Ethiopian date |
| `FromUnixMilli(ms)` | Parse Unix millisecond timestamp (JS `Date.getTime()`) → Ethiopian date |
| `d.ToISO()` | Ethiopian date → Gregorian ISO 8601 string (`"2023-09-12T00:00:00Z"`) |
| `d.UnixMilli()` | Ethiopian date → Unix milliseconds (midnight UTC, for JS `new Date(ms)`) |
| `d.MarshalJSON()` | JSON marshal as `"YYYY-MM-DD"` (Ethiopian) |
| `d.UnmarshalJSON()` | JSON unmarshal from `"YYYY-MM-DD"` (Ethiopian) |

### Date Arithmetic

| Method | Description |
|---|---|
| `d.AddDays(n)` | Add or subtract days |
| `d.AddMonths(n)` | Add or subtract months (clamps day if needed) |
| `d.AddYears(n)` | Add or subtract years (clamps Pagume 6 if needed) |

### Formatting & Inspection

| Method | Description |
|---|---|
| `d.String()` | ISO-style string: `"2016-01-01"` (implements `fmt.Stringer`) |
| `d.Format()` | Long form: `"1 Meskerem 2016"` |
| `d.Amharic()` | Amharic formatted string: `"መስከረም 1, 2016"` |
| `d.Weekday()` | Day of the week (`Segno`..`Ihud`) |
| `d.Validate()` | Check if date exists in the Ethiopian calendar |

### Calendar Utilities

| Function / Variable | Description |
|---|---|
| `EAT` | `*time.Location` for East Africa Time (UTC+3, Ethiopia's standard timezone) |
| `IsLeapYear(y)` | Report whether an Ethiopian year is a leap year (`year % 4 == 3`) |
| `DaysInMonth(y, m)` | Number of days in an Ethiopian month (30, or 5/6 for Pagume) |
| `MonthName(m)` | English transliterated month name (`Meskerem`..`Pagume`) |
| `MonthNameAmharic(m)` | Amharic month name (`መስከረም`..`ጳጉሜን`) |