# ethcal

Convert dates between the Ethiopian calendar (EC) and Gregorian calendar (GC) in Go. Zero dependencies.

```go
import "github.com/yourname/ethcal"

// Ethiopian date
d, err := ethcal.NewET(2016, ethcal.Meskerem, 1)

// EC → GC
t, _ := d.GC()
fmt.Println(t.Format("2006-01-02")) // 2023-09-12

// GC → EC
ec, err := ethcal.NewGC(2024, 9, 11)
fmt.Println(ec.FormatLong())    // 1 Meskerem 2017
fmt.Println(ec.FormatAmharic()) // መስከረም 1, 2017
fmt.Println(ec.Weekday())       // Segno

// Parse and compare
parsed, _ := ethcal.Parse("2016-01-01")
diffDays := ec.Sub(parsed)
```

## API

### Constructors & Parsing
| Function | Description |
|---|---|
| `NewET(y, m, d)` | Create and validate an Ethiopian date (`New` alias supported) |
| `NewGC(y, m, d)` | Convert Gregorian year, month, day → Ethiopian date (`FromYMD` alias supported) |
| `FromTime(t)` | Convert `time.Time` → Ethiopian date |
| `Parse(s)` | Parse `"YYYY-MM-DD"` into an Ethiopian `Date` |
| `Today()` | Current Ethiopian date in the local time zone |
| `TodayUTC()` | Current Ethiopian date in UTC |

### Date Methods
| Method | Description |
|---|---|
| `d.GC()` / `d.Gregorian()` | Convert Ethiopian date → `time.Time` |
| `d.GCYMD()` / `d.GregorianYMD()` | Convert Ethiopian date → Gregorian `(year, month, day)` |
| `d.AddDays(n)` | Add or subtract days |
| `d.AddMonths(n)` | Add or subtract months (clamps day if target month is shorter) |
| `d.AddYears(n)` | Add or subtract years (clamps Pagume 6 on non-leap years) |
| `d.Sub(other)` | Difference between dates in days (`d - other`) |
| `d.Compare(other)` | Returns -1 if `d < other`, 0 if equal, +1 if `d > other` |
| `d.Before(other)` | True if `d` is strictly before `other` |
| `d.After(other)` | True if `d` is strictly after `other` |
| `d.Equal(other)` | True if both dates represent the same year, month, and day |
| `d.IsZero()` | True if date is uninitialized (`Date{}`) |
| `d.Weekday()` | Day of the week (`Segno`..`Ihud`) |
| `d.MonthName()` | English name of the date's month |
| `d.MonthNameAmharic()` | Amharic name of the date's month |
| `d.DaysInMonth()` | Number of days in the date's month (30, or 5/6 for Pagume) |
| `d.IsLeapYear()` | True if the date's year is an Ethiopian leap year |
| `d.FormatLong()` | Formatted as `"1 Meskerem 2016"` |
| `d.FormatAmharic()` | Formatted as `"መስከረም 1, 2016"` |
| `d.String()` | ISO formatted as `"YYYY-MM-DD"` |
| `d.Validate()` | Validate date bounds |

### Helpers & Constants
- **Month constants**: `Meskerem` (1) through `Pagume` (13).
- **Weekday constants**: `Segno` (1) through `Ihud` (7).
- `IsLeapYear(year)` / `DaysInMonth(year, month)` / `MonthName(m)` / `MonthNameAmharic(m)`.
- `TimeWeekday(w)` converts Ethiopian `Weekday` to `time.Weekday`.
- `WeekdayFromTime(w)` converts `time.Weekday` to Ethiopian `Weekday`.
- Implements `encoding.TextMarshaler` and `encoding.TextUnmarshaler` for JSON/YAML serialization.