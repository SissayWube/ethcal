# ethcal

Convert dates between the Ethiopian calendar (EC) and Gregorian calendar (GC) in Go. No dependencies.

```go
import "github.com/SissayWube/ethcal"

// EC → GC
t, _ := ethcal.Date{Year: 2016, Month: 1, Day: 1}.Gregorian()
fmt.Println(t.Format("2006-01-02")) // 2023-09-12

// GC → EC
ec := ethcal.FromTime(time.Date(2024, 9, 11, 0, 0, 0, 0, time.UTC))
fmt.Println(ec.FormatLong())    // 1 Meskerem 2017
fmt.Println(ec.FormatAmharic()) // መስከረም 1, 2017
fmt.Println(ec.Weekday())       // Segno
```

## API

| Function / Method | Description |
|---|---|
| `New(y, m, d)` | Create and validate an Ethiopian date |
| `FromTime(t)` | Convert `time.Time` → Ethiopian date |
| `FromYMD(y, m, d)` | Convert Gregorian year/month/day → Ethiopian date |
| `Today()` | Current Ethiopian date (local time zone) |
| `TodayUTC()` | Current Ethiopian date (UTC) |
| `TodayEAT()` | Current Ethiopian date in East Africa Time (`EAT`, UTC+3) |
| `TodayIn(loc)` | Current Ethiopian date in the given `*time.Location` |
| `d.Gregorian()` | Ethiopian date → `time.Time` (midnight UTC) |
| `d.GregorianIn(loc)` | Ethiopian date → `time.Time` (midnight in `*time.Location`) |
| `d.GregorianYMD()` | Ethiopian date → Gregorian year, month, day |
| `d.AddDays(n)` | Add or subtract days |
| `d.AddMonths(n)` | Add or subtract months (clamps day if needed) |
| `d.AddYears(n)` | Add or subtract years (clamps Pagume 6 if needed) |
| `d.Weekday()` | Ethiopian day of the week (`Segno`..`Ihud`) |
| `d.FormatLong()` | Long form: `"1 Meskerem 2000"` |
| `d.FormatAmharic()` | Amharic: `"መስከረም 1, 2000"` |
| `d.String()` | ISO-style: `"2000-01-01"` |
| `d.Equal(other)` | Compare two dates |
| `d.Validate()` | Check date validity |
| `IsLeapYear(y)` | Ethiopian leap year check |
| `DaysInMonth(y, m)` | Days in a given month |
| `MonthName(m)` | English month name |
| `MonthNameAmharic(m)` | Amharic month name |
Other helpers: `New`, `Validate`, `IsLeapYear`, `DaysInMonth`, `MonthName`, `AddDays`, `Today`, `ToGregorianYMD`, `FromGregorianYMD`.