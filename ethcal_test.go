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

func TestMonthConstants(t *testing.T) {
	if Meskerem != 1 || Pagume != 13 {
		t.Errorf("expected Meskerem=1 and Pagume=13, got %d and %d", Meskerem, Pagume)
	}
}

func TestWeekday(t *testing.T) {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 3000; i++ {
		g := start.AddDate(0, 0, i)
		ec := FromTime(g)
		ecWd := ec.Weekday()
		goWd := g.Weekday()
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

func TestWeekdayConversions(t *testing.T) {
	for w := Segno; w <= Ihud; w++ {
		gw := TimeWeekday(w)
		back := WeekdayFromTime(gw)
		if back != w {
			t.Errorf("round trip failed for %v: got %v", w, back)
		}
	}
	if TimeWeekday(Ihud) != time.Sunday {
		t.Errorf("TimeWeekday(Ihud) = %v; want Sunday", TimeWeekday(Ihud))
	}
	if WeekdayFromTime(time.Sunday) != Ihud {
		t.Errorf("WeekdayFromTime(Sunday) = %v; want Ihud", WeekdayFromTime(time.Sunday))
	}
}

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

func TestComparisons(t *testing.T) {
	d1 := Date{2016, 1, 1}
	d2 := Date{2016, 1, 2}
	d3 := Date{2016, 2, 1}
	d4 := Date{2017, 1, 1}

	if d1.Compare(d1) != 0 || !d1.Equal(d1) {
		t.Errorf("expected %v to compare equal to itself", d1)
	}
	if d1.Compare(d2) >= 0 || !d1.Before(d2) || d1.After(d2) {
		t.Errorf("expected %v before %v", d1, d2)
	}
	if d2.Compare(d1) <= 0 || !d2.After(d1) || d2.Before(d1) {
		t.Errorf("expected %v after %v", d2, d1)
	}
	if d1.Compare(d3) >= 0 || !d1.Before(d3) {
		t.Errorf("expected %v before %v", d1, d3)
	}
	if d3.Compare(d1) <= 0 || !d3.After(d1) {
		t.Errorf("expected %v after %v", d3, d1)
	}
	if d1.Compare(d4) >= 0 || !d1.Before(d4) {
		t.Errorf("expected %v before %v", d1, d4)
	}
	if d4.Compare(d1) <= 0 || !d4.After(d1) {
		t.Errorf("expected %v after %v", d4, d1)
	}

	zero := Date{}
	if !zero.IsZero() {
		t.Error("zero date should return true for IsZero")
	}
	if d1.IsZero() {
		t.Error("initialized date should return false for IsZero")
	}
}

func TestSub(t *testing.T) {
	d1 := Date{2016, 1, 1}
	d2 := Date{2016, 1, 10}

	diff, err := d2.Sub(d1)
	if err != nil || diff != 9 {
		t.Errorf("Sub: got %d, %v; want 9, nil", diff, err)
	}

	diffNeg, err := d1.Sub(d2)
	if err != nil || diffNeg != -9 {
		t.Errorf("Sub neg: got %d, %v; want -9, nil", diffNeg, err)
	}

	diffZero, err := d1.Sub(d1)
	if err != nil || diffZero != 0 {
		t.Errorf("Sub zero: got %d, %v; want 0, nil", diffZero, err)
	}

	if _, err := (Date{0, 1, 1}).Sub(d1); err == nil {
		t.Error("Sub with invalid receiver should return error")
	}
	if _, err := d1.Sub(Date{0, 1, 1}); err == nil {
		t.Error("Sub with invalid arg should return error")
	}
}

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

func TestDateConvenienceMethods(t *testing.T) {
	d := Date{2015, 13, 6}
	if d.MonthName() != "Pagume" {
		t.Errorf("MonthName() = %q; want Pagume", d.MonthName())
	}
	if d.MonthNameAmharic() != "ጳጉሜን" {
		t.Errorf("MonthNameAmharic() = %q; want ጳጉሜን", d.MonthNameAmharic())
	}
	days, err := d.DaysInMonth()
	if err != nil || days != 6 {
		t.Errorf("DaysInMonth() = %d, %v; want 6, nil", days, err)
	}
	if !d.IsLeapYear() {
		t.Error("2015 is leap year, IsLeapYear should return true")
	}
}

func TestDaysInMonth(t *testing.T) {
	if _, err := DaysInMonth(0, 1); err != ErrInvalidYear {
		t.Errorf("expected ErrInvalidYear for year 0, got %v", err)
	}
	if _, err := DaysInMonth(2016, 0); err != ErrInvalidMonth {
		t.Errorf("expected ErrInvalidMonth for month 0, got %v", err)
	}
	if _, err := DaysInMonth(2016, 14); err != ErrInvalidMonth {
		t.Errorf("expected ErrInvalidMonth for month 14, got %v", err)
	}
	if n, err := DaysInMonth(2016, 5); err != nil || n != 30 {
		t.Errorf("DaysInMonth(2016, 5) = %d, %v; want 30, nil", n, err)
	}
	if n, err := DaysInMonth(2016, 13); err != nil || n != 5 {
		t.Errorf("DaysInMonth(2016, 13) non-leap = %d, %v; want 5, nil", n, err)
	}
	if n, err := DaysInMonth(2015, 13); err != nil || n != 6 {
		t.Errorf("DaysInMonth(2015, 13) leap = %d, %v; want 6, nil", n, err)
	}
}

func TestNewET(t *testing.T) {
	d, err := NewET(2016, 1, 1)
	if err != nil || d != (Date{2016, 1, 1}) {
		t.Errorf("NewET(2016,1,1) = %v, %v", d, err)
	}
	d2, err := New(2016, 1, 1)
	if err != nil || d2 != d {
		t.Errorf("New alias mismatch: %v vs %v", d2, d)
	}
	if _, err := NewET(0, 1, 1); err == nil {
		t.Error("NewET(0,1,1) should fail")
	}
}

func TestNewGC(t *testing.T) {
	d, err := NewGC(2023, time.September, 12)
	if err != nil || d != (Date{2016, 1, 1}) {
		t.Errorf("NewGC(2023, Sep, 12) = %v, %v", d, err)
	}
	d2, err := FromYMD(2023, time.September, 12)
	if err != nil || d2 != d {
		t.Errorf("FromYMD alias mismatch: %v vs %v", d2, d)
	}
	if _, err := NewGC(2023, time.February, 30); err == nil {
		t.Error("Feb 30 should be invalid")
	}
	if _, err := NewGC(1, time.January, 1); err == nil {
		t.Error("pre-epoch Gregorian date should return error")
	}
}

func TestParse(t *testing.T) {
	valid := []struct {
		in   string
		want Date
	}{
		{"2016-01-01", Date{2016, 1, 1}},
		{"2015-13-06", Date{2015, 13, 6}},
		{"0001-01-01", Date{1, 1, 1}},
	}
	for _, tc := range valid {
		got, err := Parse(tc.in)
		if err != nil || got != tc.want {
			t.Errorf("Parse(%q) = %v, %v; want %v, nil", tc.in, got, err, tc.want)
		}
	}

	invalid := []string{
		"2016-01",
		"2016-01-01-01",
		"abcd-01-01",
		"2016-ab-01",
		"2016-01-cd",
		"2016-14-01",
		"2016-13-06", // 2016 is non-leap
		"0000-01-01",
	}
	for _, in := range invalid {
		if _, err := Parse(in); err == nil {
			t.Errorf("Parse(%q) expected error, got nil", in)
		}
	}
}

func TestTextMarshaling(t *testing.T) {
	d := Date{2016, 1, 1}
	b, err := d.MarshalText()
	if err != nil || string(b) != "2016-01-01" {
		t.Errorf("MarshalText: got %s, %v; want 2016-01-01", b, err)
	}

	var d2 Date
	if err := d2.UnmarshalText([]byte("2016-01-01")); err != nil || d2 != d {
		t.Errorf("UnmarshalText: got %v, %v; want %v", d2, err, d)
	}

	if _, err := (Date{0, 1, 1}).MarshalText(); err == nil {
		t.Error("MarshalText on invalid date should return error")
	}
	if err := d2.UnmarshalText([]byte("invalid")); err == nil {
		t.Error("UnmarshalText on invalid date string should return error")
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

	dUTC := TodayUTC()
	if err := dUTC.Validate(); err != nil {
		t.Errorf("TodayUTC() = %v is invalid: %v", dUTC, err)
	}
}

func TestGregorianYMD(t *testing.T) {
	d := Date{2016, 1, 1}
	y, m, day, err := d.GregorianYMD()
	if err != nil || y != 2023 || m != time.September || day != 12 {
		t.Errorf("GregorianYMD() = %d-%v-%d, %v", y, m, day, err)
	}
	y2, m2, day2, err := d.GCYMD()
	if err != nil || y2 != y || m2 != m || day2 != day {
		t.Errorf("GCYMD() mismatch: %d-%v-%d, %v", y2, m2, day2, err)
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
	_, err = (Date{0, 1, 1}).GC()
	if err == nil {
		t.Error("GC on year 0 should fail")
	}
	_, err = (Date{2016, 14, 1}).Gregorian()
	if err == nil {
		t.Error("Gregorian on month 14 should fail")
	}
}

func TestAddDays(t *testing.T) {
	d := Date{2016, 1, 1}
	r, err := d.AddDays(0)
	if err != nil || r != d {
		t.Errorf("AddDays(0) = %v, %v", r, err)
	}
	r, _ = d.AddDays(30)
	if r != (Date{2016, 2, 1}) {
		t.Errorf("AddDays(30) = %v; want 2016-02-01", r)
	}
	r, _ = d.AddDays(365)
	if r != (Date{2017, 1, 1}) {
		t.Errorf("AddDays(365) = %v; want 2017-01-01", r)
	}
	r, _ = d.AddDays(-1)
	if r != (Date{2015, 13, 6}) {
		t.Errorf("AddDays(-1) from 2016-01-01 = %v; want 2015-13-06", r)
	}
	if _, err := (Date{0, 1, 1}).AddDays(1); err == nil {
		t.Error("AddDays on invalid date should fail")
	}
	if _, err := (Date{1, 1, 1}).AddDays(-1); err != ErrInvalidYear {
		t.Errorf("AddDays underflow before year 1: got %v, want ErrInvalidYear", err)
	}
}

func TestAddMonths(t *testing.T) {
	cases := []struct {
		name   string
		start  Date
		months int
		want   Date
	}{
		{"forward 1", Date{2016, 1, 15}, 1, Date{2016, 2, 15}},
		{"forward 12 clamps Pagume", Date{2016, 1, 15}, 12, Date{2016, 13, 5}},
		{"forward 13 wraps to next year", Date{2016, 1, 15}, 13, Date{2017, 1, 15}},
		{"backward 1 from month 1", Date{2016, 1, 15}, -1, Date{2015, 13, 6}},
		{"backward 1 from month 2", Date{2016, 2, 15}, -1, Date{2016, 1, 15}},
		{"backward 13", Date{2016, 1, 15}, -13, Date{2015, 1, 15}},
		{"to Pagume non-leap clamp", Date{2016, 12, 30}, 1, Date{2016, 13, 5}},
		{"to Pagume leap no clamp", Date{2015, 12, 6}, 1, Date{2015, 13, 6}},
		{"zero months", Date{2016, 6, 15}, 0, Date{2016, 6, 15}},
		{"large forward", Date{2016, 1, 1}, 26, Date{2018, 1, 1}},
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
