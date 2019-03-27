package cron

import (
	"testing"
	"time"
)

func TestSchedule1(t *testing.T) {
	schedule, _ := Parse("* * * * *")
	now, _ := time.Parse(time.RFC3339, "2018-01-01T10:00:00+08:00")
	want, _ := time.Parse(time.RFC3339, "2018-01-01T10:01:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Logf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}

func TestSchedule2(t *testing.T) {
	schedule, _ := Parse("* * * * *")
	now, _ := time.Parse(time.RFC3339, "2018-01-01T10:00:59+08:00")
	want, _ := time.Parse(time.RFC3339, "2018-01-01T10:01:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Errorf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}

func TestSchedule3(t *testing.T) {
	schedule, _ := Parse("5 * * * *")
	now, _ := time.Parse(time.RFC3339, "2018-01-01T10:00:59+08:00")
	want, _ := time.Parse(time.RFC3339, "2018-01-01T10:05:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Errorf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}

func TestSchedule4(t *testing.T) {
	schedule, _ := Parse("0 * * * *")
	now, _ := time.Parse(time.RFC3339, "2018-01-01T10:00:00+08:00")
	want, _ := time.Parse(time.RFC3339, "2018-01-01T11:00:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Errorf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}

func TestSchedule5(t *testing.T) {
	schedule, _ := Parse("20-30 * * * *")
	now, _ := time.Parse(time.RFC3339, "2018-01-01T10:00:00+08:00")
	want, _ := time.Parse(time.RFC3339, "2018-01-01T10:20:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Errorf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}

func TestSchedule6(t *testing.T) {
	schedule, _ := Parse("20-30,5,6,7,8,9 * * * *")
	now, _ := time.Parse(time.RFC3339, "2018-01-01T10:00:00+08:00")
	want, _ := time.Parse(time.RFC3339, "2018-01-01T10:05:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Errorf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}

func TestSchedule7(t *testing.T) {
	schedule, _ := Parse("20-30,5,6,7,8,9 23 * * *")
	now, _ := time.Parse(time.RFC3339, "2018-01-01T10:00:00+08:00")
	want, _ := time.Parse(time.RFC3339, "2018-01-01T23:05:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Errorf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}

func TestSchedule8(t *testing.T) {
	schedule, _ := Parse("20-30,5,6,7,8,9 23 2 * *")
	now, _ := time.Parse(time.RFC3339, "2018-01-01T10:00:00+08:00")
	want, _ := time.Parse(time.RFC3339, "2018-01-02T23:05:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Errorf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}

func TestSchedule9(t *testing.T) {
	schedule, _ := Parse("20-30,5,6,7,8,9 23 2 5 *")
	now, _ := time.Parse(time.RFC3339, "2018-01-01T10:00:00+08:00")
	want, _ := time.Parse(time.RFC3339, "2018-05-02T23:05:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Errorf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}

func TestSchedule10(t *testing.T) {
	schedule, _ := Parse("20-30,5,6,7,8,9 23 1 * 2")
	now, _ := time.Parse(time.RFC3339, "2018-01-01T10:00:00+08:00")
	want, _ := time.Parse(time.RFC3339, "2018-05-01T23:05:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Errorf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}

func TestSchedule11(t *testing.T) {
	schedule, _ := Parse("0 0 29 2 *")
	now, _ := time.Parse(time.RFC3339, "2008-02-29T10:00:00+08:00")
	want, _ := time.Parse(time.RFC3339, "2012-02-29T00:00:00+08:00")

	next := schedule.Next(now)

	if next != want {
		t.Logf("year: %v", next.Year())
		t.Errorf("now: %v, want: %v, get: %v", now, want, next)
	}

	return
}
