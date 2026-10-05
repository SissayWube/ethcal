# ethcal

Convert dates between the Ethiopian calendar (EC) and Gregorian calendar (GC) in Go. No dependencies.

```go
import "github.com/sisaywube/ethcal"

// EC -> GC
t, _ := ethcal.Date{Year: 2016, Month: 1, Day: 1}.ToGregorian()
fmt.Println(t.Format("2006-01-02")) // 2023-09-12

// GC -> EC
ec := ethcal.FromGregorian(time.Date(2024, 9, 11, 0, 0, 0, 0, time.UTC))
fmt.Println(ec.Format())        // 1 Meskerem 2017
fmt.Println(ec.FormatAmharic()) // መስከረም 1, 2017
```

Other helpers: `New`, `Validate`, `IsLeapYear`, `DaysInMonth`, `MonthName`, `AddDays`, `Today`, `ToGregorianYMD`, `FromGregorianYMD`.