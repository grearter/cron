package cron

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ParseError struct {
	Field  string
	Reason string
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("%s: %s", pe.Reason, pe.Field)
}

type bound struct {
	min int64
	max int64
}

var (
	minuteBound = &bound{min: 0, max: 59}
	hourBound   = &bound{min: 0, max: 23}
	dayBound    = &bound{min: 1, max: 31}
	monthBound  = &bound{min: 1, max: 12}
	dowBound    = &bound{min: 0, max: 6}
)

type Schedule struct {
	Month     uint64 // bit1~bit12
	Day       uint64 // bit1~bit31
	DayOfWeek uint64 // bit0~bit6(周日～周六)
	Hour      uint64 // bit0~bit23
	Minute    uint64 // bit0~bit59
}

func (schedule *Schedule) Next(now time.Time) (next time.Time) {

	next = now

	second := next.Second()
	next = next.Add(time.Duration(60-second) * time.Second)

	// 确定minute
	minute := uint64(next.Minute())
	for ((1 << minute) & schedule.Minute) == 0 {
		next = next.Add(time.Minute)
		minute = uint64(next.Minute())
	}

	// 确定hour
	hour := uint64(next.Hour())
	for ((1 << hour) & schedule.Hour) == 0 {
		next = next.Add(time.Hour)
		hour = uint64(next.Hour())
	}

	// 确定day, month
	var (
		day   = uint64(next.Day())
		month = uint64(next.Month())
		dow   = uint64(next.Weekday())
	)
	for (((1 << day) & schedule.Day) == 0) || (((1 << month) & schedule.Month) == 0) || (((1 << dow) & schedule.DayOfWeek) == 0) {
		next = next.Add(time.Hour * 24)

		if next.After(now.Add(time.Hour * 24 * 366 * 4)) { // 最大向后找4年
			next = time.Unix(0, 0)
			break
		}

		day = uint64(next.Day())
		month = uint64(next.Month())
		dow = uint64(next.Weekday())
	}

	next.IsZero()
	return
}

/*
* * * * * "command to be executed"
- - - - -
| | | | |
| | | | ----- Day of week (0 - 6) (Sunday=0)
| | | ------- Month (1 - 12)
| | --------- Day of month (1 - 31)
| ----------- Hour (0 - 23)
------------- Minute (0 - 59)
*/
func Parse(scheduleStr string) (schedule *Schedule, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
			return
		}
	}()

	fields := strings.Fields(scheduleStr)

	if len(scheduleStr) <= 5 {
		err = fmt.Errorf("invalid schedule time: %v", scheduleStr)
		return
	}

	minute, err := parseField(fields[0], minuteBound)
	if err != nil {
		return
	}

	hour, err := parseField(fields[1], hourBound)
	if err != nil {
		return
	}

	day, err := parseField(fields[2], dayBound)
	if err != nil {
		return
	}

	month, err := parseField(fields[3], monthBound)
	if err != nil {
		return
	}

	dow, err := parseField(fields[4], dowBound)
	if err != nil {
		return
	}

	schedule = &Schedule{
		Minute:    minute,
		Hour:      hour,
		Day:       day,
		Month:     month,
		DayOfWeek: dow,
	}
	return
}

func parseField(field string, boundInfo *bound) (bits uint64, err error) {

	min, max := boundInfo.min, boundInfo.max

	if field == "*" {
		for i := min; i <= max; i++ {
			bits |= 1 << uint64(i)
		}
		return
	}

	return setBits(field, min, max)
}

// field format: *  1,2,3  0-9,10,11,12  */3
func setBits(field string, min, max int64) (bits uint64, err error) {

	subFields := strings.FieldsFunc(field, func(r rune) bool { return r == ',' })

	if len(subFields) <= 0 {
		err = &ParseError{
			Field:  field,
			Reason: "format error",
		}

		return
	}

	for _, subField := range subFields {
		if strings.Contains(subField, "-") {
			rangeBits, e := setRangeBits(subField, min, max)
			if e != nil {
				err = e
				return
			}

			bits |= rangeBits
		} else if strings.HasPrefix(subField, "*/") {
			modBits, e := setModBits(subField, min, max)
			if e != nil {
				err = e
				return
			}

			bits |= modBits
		} else {
			var iSubField int64
			iSubField, e := strconv.ParseInt(subField, 10, 64)

			if e != nil || iSubField > max || iSubField < min {
				err = &ParseError{
					Field:  field,
					Reason: "format error",
				}
				return
			}

			if iSubField > max {
				err = &ParseError{
					Field:  field,
					Reason: fmt.Sprintf("filed(%v) is greater than MaxValue(%v)", iSubField, max),
				}
				return
			}

			if iSubField < min {
				err = &ParseError{
					Field:  field,
					Reason: fmt.Sprintf("field(%v) is less than MinValue(%v)", iSubField, min),
				}
				return
			}

			bits |= 1 << uint64(iSubField)
		}
	}

	return
}

// field format: */2
func setModBits(field string, min, max int64) (bits uint64, err error) {
	if len(field) < 3 {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("format error"),
		}
		return
	}

	iSubField, e := strconv.ParseInt(field[2:], 10, 64)
	if e != nil {
		err = &ParseError{
			Field:  field,
			Reason: "format error",
		}

		return
	}

	if iSubField > max {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("mod(%v) is greater than MaxValue(%v)", iSubField, max),
		}
		return
	}

	if iSubField < min {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("mod(%v) is less than MinValue(%v)", iSubField, min),
		}
		return
	}

	if iSubField <= 0 {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("mod(%v) must be greater than 0", iSubField),
		}
		return
	}

	for i := min; i <= max; i += iSubField {
		bits |= 1 << uint64(i)
	}

	return
}

// field format: 0-9
func setRangeBits(field string, min, max int64) (bits uint64, err error) {

	subFields := strings.FieldsFunc(field, func(r rune) bool { return r == '-' })

	if len(subFields) != 2 {
		err = &ParseError{
			Field:  field,
			Reason: "format error",
		}
		return
	}

	start, end := subFields[0], subFields[1]

	iStart, e := strconv.ParseInt(start, 10, 64)
	if e != nil {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("left-range-value is not integer"),
		}
		return
	}

	iEnd, e := strconv.ParseInt(end, 10, 64)
	if e != nil {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("right-range-value is not integer"),
		}
		return
	}

	if iStart > iEnd {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("left-range-value must be less or equal than right-range-value"),
		}
		return
	}

	if iStart > max {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("left-range-range(%v) is greater than MaxValue(%v)", iStart, max),
		}
		return
	}

	if iStart < min {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("left-range-range(%v) is less than MinValuelue(%v)", iStart, min),
		}
		return
	}

	if iEnd > max {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("right-range-range(%v) is greater than MaxValue(%v)", iEnd, max),
		}
		return
	}

	if iEnd < min {
		err = &ParseError{
			Field:  field,
			Reason: fmt.Sprintf("right-range-range(%v) is less than MinValuelue(%v)", iEnd, min),
		}
		return
	}

	for i := iStart; i <= iEnd; i++ {
		bits |= 1 << uint64(i)
	}

	return
}
