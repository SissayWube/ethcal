package ethcal

import (
	"testing"
	"time"
)

func TestKnownDates(t *testing.T) {
	cases := []struct {
		ec Date
		gc time.Time
	}{
		{Date{2000, 1, 1}, time.Date(2007, 9, 12, 0, 0, 0, 0, time.UTC)},
		{Date{2016, 1, 1}, time.Date(2023, 9, 12, 0, 0, 0, 0, time.UTC)},
		{Date{2017, 1, 1}, time.Date(2024, 9, 11, 0, 0, 0, 0, time.UTC)},
		{Date{2018, 1, 1}, time.Date(2025, 9, 11, 0, 0, 0, 0, time.UTC)},
		{Date{2015, 13, 6}, time.Date(2023, 9, 11, 0, 0, 0, 0, time.UTC)},
		{Date{2016, 4, 29}, time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)},
		{Date{1, 1, 1}, time.Date(8, 8, 27, 0, 0, 0, 0, time.UTC)},
	}
	for _, c := range cases {
		got, err := c.ec.Gregorian()
		if err != nil || !got.Equal(c.gc) {
			t.Errorf("Gregorian(%v) = %v, %v; want %v", c.ec, got, err, c.gc)
		}
		if back := FromTime(c.gc); back != c.ec {
			t.Errorf("FromTime(%v) = %v; want %v", c.gc, back, c.ec)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	start := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 80000; i++ {
		g := start.AddDate(0, 0, i)
		e := FromTime(g)
		if err := e.Validate(); err != nil {
			t.Fatalf("%v -> invalid %v: %v", g, e, err)
		}
		back, _ := e.Gregorian()
		if !back.Equal(g) {
			t.Fatalf("round trip failed: %v -> %v -> %v", g, e, back)
		}
	}
}

func TestValidation(t *testing.T) {
	bad := []Date{{2015, 13, 7}, {2016, 13, 7}, {2016, 14, 1}, {2016, 0, 1}, {2016, 1, 31}, {0, 1, 1}}
	for _, d := range bad {
		if d.Validate() == nil {
			t.Errorf("%v should be invalid", d)
		}
	}
	if (Date{2015, 13, 6}).Validate() != nil {
		t.Error("2015 is a leap year; Pagume 6 should be valid")
	}
}

// ---- Weekday ----

func TestWeekday(t *testing.T) {
	// Cross-check: every EC date's weekday must match its Gregorian weekday.
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 3000; i++ {
		g := start.AddDate(0, 0, i)
		ec := FromTime(g)
		ecWd := ec.Weekday()
		goWd := g.Weekday()
		// Segno=0→Monday=time.Monday(1), ..., Ihud=6→Sunday=time.Sunday(0)
		var expected time.Weekday
		if ecWd == Ihud {
			expected = time.Sunday
		} else {
			expected = time.Weekday(ecWd + 1)
		}
		if goWd != expected {
			t.Fatalf("%v (EC %v): Weekday()=%v (mapped to %v) but Go says %v",
				g.Format("2006-01-02"), ec, ecWd, expected, goWd)
		}
	}
}

func TestWeekdayNames(t *testing.T) {
	if Segno.String() != "Segno" {
		t.Errorf("Segno.String() = %q", Segno.String())
	}
	if Ihud.Amharic() != "እሁድ" {
		t.Errorf("Ihud.Amharic() = %q", Ihud.Amharic())
	}
	if Weekday(-1).String() != "" || Weekday(7).String() != "" {
		t.Error("out-of-range weekday should return empty string")
	}
	if Weekday(-1).Amharic() != "" || Weekday(7).Amharic() != "" {
		t.Error("out-of-range weekday amharic should return empty string")
	}
}

// ---- Equal ----

func TestEqual(t *testing.T) {
	a := Date{2016, 1, 1}
	b := Date{2016, 1, 1}
	c := Date{2016, 1, 2}
	if !a.Equal(b) {
		t.Error("same dates should be equal")
	}
	if a.Equal(c) {
		t.Error("different dates should not be equal")
	}
}

// ---- String / FormatLong / FormatAmharic ----

func TestString(t *testing.T) {
	d := Date{2016, 1, 5}
	if s := d.String(); s != "2016-01-05" {
		t.Errorf("String() = %q", s)
	}
	d2 := Date{1, 1, 1}
	if s := d2.String(); s != "0001-01-01" {
		t.Errorf("String() = %q for year 1", s)
	}
}

func TestFormatLong(t *testing.T) {
	d := Date{2017, 1, 1}
	if s := d.FormatLong(); s != "1 Meskerem 2017" {
		t.Errorf("FormatLong() = %q", s)
	}
}

func TestFormatAmharic(t *testing.T) {
	d := Date{2017, 1, 1}
	if s := d.FormatAmharic(); s != "መስከረም 1, 2017" {
		t.Errorf("FormatAmharic() = %q", s)
	}
}

// ---- MonthName ----

func TestMonthName(t *testing.T) {
	if MonthName(1) != "Meskerem" {
		t.Errorf("month 1 = %q", MonthName(1))
	}
	if MonthName(13) != "Pagume" {
		t.Errorf("month 13 = %q", MonthName(13))
	}
	if MonthName(0) != "" || MonthName(14) != "" {
		t.Error("out-of-range month should return empty string")
	}
}

func TestMonthNameAmharic(t *testing.T) {
	if MonthNameAmharic(1) != "መስከረም" {
		t.Errorf("month 1 amharic = %q", MonthNameAmharic(1))
	}
	if MonthNameAmharic(0) != "" || MonthNameAmharic(14) != "" {
		t.Error("out-of-range month should return empty string")
	}
}

// ---- Constructors ----

func TestNew(t *testing.T) {
	d, err := New(2016, 1, 1)
	if err != nil || d != (Date{2016, 1, 1}) {
		t.Errorf("New(2016,1,1) = %v, %v", d, err)
	}
	_, err = New(0, 1, 1)
	if err == nil {
		t.Error("New(0,1,1) should fail")
	}
}

func TestFromYMD(t *testing.T) {
	d, err := FromYMD(2023, time.September, 12)
	if err != nil || d != (Date{2016, 1, 1}) {
		t.Errorf("FromYMD(2023, Sep, 12) = %v, %v", d, err)
	}
	_, err = FromYMD(2023, time.February, 30)
	if err == nil {
		t.Error("Feb 30 should be invalid")
	}
}

func TestToday(t *testing.T) {
	d := Today()
	if err := d.Validate(); err != nil {
		t.Errorf("Today() = %v is invalid: %v", d, err)
	}
	if d.Year < 2010 {
		t.Errorf("Today() year = %d; seems too low", d.Year)
	}
}

// ---- Gregorian conversions ----

func TestGregorianYMD(t *testing.T) {
	y, m, d, err := (Date{2016, 1, 1}).GregorianYMD()
	if err != nil || y != 2023 || m != time.September || d != 12 {
		t.Errorf("GregorianYMD() = %d-%v-%d, %v", y, m, d, err)
	}
	_, _, _, err = (Date{0, 1, 1}).GregorianYMD()
	if err == nil {
		t.Error("year 0 should fail")
	}
}

func TestGregorian_Invalid(t *testing.T) {
	_, err := (Date{0, 1, 1}).Gregorian()
	if err == nil {
		t.Error("Gregorian on year 0 should fail")
	}
	_, err = (Date{2016, 14, 1}).Gregorian()
	if err == nil {
		t.Error("Gregorian on month 14 should fail")
	}
}

// ---- AddDays ----

func TestAddDays(t *testing.T) {
	d := Date{2016, 1, 1}
	// Add 0
	r, err := d.AddDays(0)
	if err != nil || r != d {
		t.Errorf("AddDays(0) = %v, %v", r, err)
	}
	// Add 30 -> month 2 day 1
	r, _ = d.AddDays(30)
	if r != (Date{2016, 2, 1}) {
		t.Errorf("AddDays(30) = %v; want 2016-02-01", r)
	}
	// Add 365 (non-leap year) -> next new year
	r, _ = d.AddDays(365)
	if r != (Date{2017, 1, 1}) {
		t.Errorf("AddDays(365) = %v; want 2017-01-01", r)
	}
	// Subtract
	r, _ = d.AddDays(-1)
	if r != (Date{2015, 13, 6}) { // 2015 is leap, Pagume has 6 days
		t.Errorf("AddDays(-1) from 2016-01-01 = %v; want 2015-13-06", r)
	}
	// Invalid source
	_, err = (Date{0, 1, 1}).AddDays(1)
	if err == nil {
		t.Error("AddDays on invalid date should fail")
	}
}

// ---- AddMonths ----

func TestAddMonths(t *testing.T) {
	cases := []struct {
		name   string
		start  Date
		months int
		want   Date
	}{
		{"forward 1", Date{2016, 1, 15}, 1, Date{2016, 2, 15}},
		{"forward 12 clamps Pagume", Date{2016, 1, 15}, 12, Date{2016, 13, 5}}, // Pagume has only 5 days in non-leap 2016
		{"forward 13 wraps to next year", Date{2016, 1, 15}, 13, Date{2017, 1, 15}},
		{"backward 1 from month 1", Date{2016, 1, 15}, -1, Date{2015, 13, 6}}, // Pagume clamp: day 15 -> 6 (leap year 2015)
		{"backward 1 from month 2", Date{2016, 2, 15}, -1, Date{2016, 1, 15}},
		{"backward 13", Date{2016, 1, 15}, -13, Date{2015, 1, 15}},
		{"to Pagume non-leap clamp", Date{2016, 12, 30}, 1, Date{2016, 13, 5}}, // 2016 is non-leap, Pagume has 5 days
		{"to Pagume leap no clamp", Date{2015, 12, 6}, 1, Date{2015, 13, 6}},   // 2015 is leap, Pagume has 6 days
		{"zero months", Date{2016, 6, 15}, 0, Date{2016, 6, 15}},
		{"large forward", Date{2016, 1, 1}, 26, Date{2018, 1, 1}}, // 26 months = 2 years
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.start.AddMonths(c.months)
			if err != nil {
				t.Fatalf("AddMonths(%d): %v", c.months, err)
			}
			if got != c.want {
				t.Errorf("(%v).AddMonths(%d) = %v; want %v", c.start, c.months, got, c.want)
			}
		})
	}
}

func TestAddMonths_Invalid(t *testing.T) {
	_, err := (Date{0, 1, 1}).AddMonths(1)
	if err == nil {
		t.Error("AddMonths on invalid date should fail")
	}
	_, err = (Date{1, 1, 1}).AddMonths(-1)
	if err != ErrInvalidYear {
		t.Errorf("AddMonths before epoch: got %v; want ErrInvalidYear", err)
	}
}

// ---- AddYears ----

func TestAddYears(t *testing.T) {
	cases := []struct {
		name  string
		start Date
		years int
		want  Date
	}{
		{"forward 1", Date{2016, 6, 15}, 1, Date{2017, 6, 15}},
		{"backward 1", Date{2016, 6, 15}, -1, Date{2015, 6, 15}},
		{"Pagume 6 leap to non-leap clamps", Date{2015, 13, 6}, 1, Date{2016, 13, 5}},
		{"Pagume 5 non-leap to leap stays", Date{2016, 13, 5}, -1, Date{2015, 13, 5}},
		{"Pagume 6 leap to leap stays", Date{2015, 13, 6}, 4, Date{2019, 13, 6}},
		{"zero years", Date{2016, 6, 15}, 0, Date{2016, 6, 15}},
		{"regular month unaffected", Date{2015, 6, 30}, 1, Date{2016, 6, 30}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.start.AddYears(c.years)
			if err != nil {
				t.Fatalf("AddYears(%d): %v", c.years, err)
			}
			if got != c.want {
				t.Errorf("(%v).AddYears(%d) = %v; want %v", c.start, c.years, got, c.want)
			}
		})
	}
}

func TestAddYears_Invalid(t *testing.T) {
	_, err := (Date{0, 1, 1}).AddYears(1)
	if err == nil {
		t.Error("AddYears on invalid date should fail")
	}
	_, err = (Date{1, 1, 1}).AddYears(-1)
	if err != ErrInvalidYear {
		t.Errorf("AddYears before epoch: got %v; want ErrInvalidYear", err)
	}
}
