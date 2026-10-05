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
		got, err := c.ec.ToGregorian()
		if err != nil || !got.Equal(c.gc) {
			t.Errorf("ToGregorian(%v) = %v, %v; want %v", c.ec, got, err, c.gc)
		}
		if back := FromGregorian(c.gc); back != c.ec {
			t.Errorf("FromGregorian(%v) = %v; want %v", c.gc, back, c.ec)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	start := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 80000; i++ {
		g := start.AddDate(0, 0, i)
		e := FromGregorian(g)
		if err := e.Validate(); err != nil {
			t.Fatalf("%v -> invalid %v: %v", g, e, err)
		}
		back, _ := e.ToGregorian()
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
