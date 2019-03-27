package cron

import "testing"

func TestParse1(t *testing.T) {
	input := "*"
	want := uint64(0)
	bound := &bound{0, 59}

	for i := bound.min; i <= bound.max; i++ {
		want |= 1 << uint64(i)
	}

	bits, err := parseField(input, bound)

	if bits != want {
		t.Errorf("want: %b, get: %b, err: %v", want, bits, err)
	}

	return
}

func TestParse2(t *testing.T) {
	input := "*"
	bound := &bound{1, 31}
	want := uint64(0)

	for i := bound.min; i <= bound.max; i++ {
		want |= 1 << uint64(i)
	}

	bits, err := parseField(input, bound)
	if bits != want {
		t.Errorf("want: %b, get: %b, err: %v", want, bits, err)
	}
}

func TestParse3(t *testing.T) {
	input := "*12312312"
	bound := &bound{0, 23}

	bits, err := parseField(input, bound)

	if err == nil {
		t.Errorf("err is NIL, bits: %v", bits)
	}

	t.Logf("err: %v", err)
}

func TestParse4(t *testing.T) {
	input := "*/2"
	bound := &bound{0, 23}

	want := uint64(0)

	for i := bound.min; i <= bound.max; i += 2 {
		want |= 1 << uint64(i)
	}

	bits, err := parseField(input, bound)
	if bits != want {
		t.Errorf("want: %b, get: %b, err: %v", want, bits, err)
	}
}

func TestParse5(t *testing.T) {
	input := "*/100"
	bound := &bound{1, 31}

	bits, err := parseField(input, bound)
	if err == nil {
		t.Errorf("err is not NIL, bits: %v", bits)
	}

	t.Logf("err: %v", err)
}

func TestParse6(t *testing.T) {
	input := "*/6"
	bound := &bound{10, 31}

	bits, err := parseField(input, bound)
	if err == nil {
		t.Errorf("err is not NIL, bits: %v", bits)
	}

	t.Logf("err: %v", err)
}

func TestParse7(t *testing.T) {
	input := "*/0"
	bound := &bound{0, 6}

	bits, err := parseField(input, bound)
	if err == nil {
		t.Errorf("err is not NIL, bits: %v", bits)
	}

	t.Logf("err: %v", err)
}

func TestParse8(t *testing.T) {
	input := "*/"
	bound := &bound{0, 6}

	bits, err := parseField(input, bound)
	if err == nil {
		t.Errorf("err is not NIL, bits: %v", bits)
	}

	t.Logf("err: %v", err)
}

func TestParse9(t *testing.T) {
	input := "*/hello"
	bound := &bound{0, 6}

	bits, err := parseField(input, bound)
	if err == nil {
		t.Errorf("err is not NIL, bits: %v", bits)
	}

	t.Logf("err: %v", err)
}

func TestParse10(t *testing.T) {
	input := "5,10,13,17,23"
	bound := &bound{0, 23}

	want := uint64(0)

	for _, num := range []uint64{5, 10, 13, 17, 23} {
		want |= 1 << num
	}

	bits, err := parseField(input, bound)
	if err != nil {
		t.Errorf("want: %b, get: %b, err: %v", want, bits, err)
	}
}

func TestParse11(t *testing.T) {
	input := "0-9"
	bound := &bound{0, 23}

	want := uint64(0)

	for _, num := range []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		want |= 1 << num
	}

	bits, err := parseField(input, bound)
	if err != nil {
		t.Errorf("want: %b, get: %b, err: %v", want, bits, err)
	}
}

func TestParse12(t *testing.T) {
	input := "0-9,10,11,12,13"
	bound := &bound{0, 23}

	want := uint64(0)

	for _, num := range []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13} {
		want |= 1 << num
	}

	bits, err := parseField(input, bound)
	if err != nil {
		t.Errorf("want: %b, get: %b, err: %v", want, bits, err)
	}
}

func TestParse13(t *testing.T) {
	input := "0-9,0,1,2,3,4"
	bound := &bound{0, 23}

	want := uint64(0)

	for _, num := range []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		want |= 1 << num
	}

	bits, err := parseField(input, bound)
	if err != nil {
		t.Errorf("want: %b, get: %b, err: %v", want, bits, err)
	}
}

func TestParse14(t *testing.T) {
	input := "0-9,100,200"
	bound := &bound{0, 23}

	want := uint64(0)

	for _, num := range []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		want |= 1 << num
	}

	bits, err := parseField(input, bound)
	if err == nil {
		t.Errorf("err is not NIL, bits: %v", bits)
	}

	t.Logf("err: %v", err)
}

func TestParse15(t *testing.T) {
	input := "0-100,0,1,2,3,4"
	bound := &bound{0, 23}

	want := uint64(0)

	for _, num := range []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		want |= 1 << num
	}

	bits, err := parseField(input, bound)
	if err == nil {
		t.Errorf("err is not NIL, bits: %v", bits)
	}

	t.Logf("err: %v", err)
}

func TestParse16(t *testing.T) {
	input := "0-10,a,b,c,"
	bound := &bound{0, 23}

	want := uint64(0)

	for _, num := range []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		want |= 1 << num
	}

	bits, err := parseField(input, bound)
	if err == nil {
		t.Errorf("err is not NIL, bits: %v", bits)
	}

	t.Logf("err: %v", err)
}

func TestParse17(t *testing.T) {
	input := "0-7,,,,"
	bound := &bound{0, 23}

	want := uint64(0)

	for _, num := range []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		want |= 1 << num
	}

	bits, err := parseField(input, bound)
	if err != nil {
		t.Errorf("want: %b, get: %b, err: %v", want, bits, err)
	}

	t.Logf("want: %b, get: %v", want, bits)
}
