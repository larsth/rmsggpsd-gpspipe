package daemon

import (
	"fmt"
	"strconv"
	"strings"
)

func errorTest(got error, want error) (s string, ok bool) {
	var (
		wantStr string
		gotStr  string
	)
	if want == nil {
		wantStr = `<nil>`
	} else {
		wantStr = want.Error()
	}
	if got == nil {
		gotStr = `<nil>`
	} else {
		gotStr = got.Error()
	}
	if strings.Compare(wantStr, gotStr) != 0 {
		format := "%s%s'\n, %s%s'"
		s1 := `Got the error '`
		s2 := `... but want the error '`
		s = fmt.Sprintf(format, s1, gotStr, s2, wantStr)
		ok = false
	} else {
		s = ""
		ok = true
	}
	return
}

func errorTestI(got error, want error, i int, tdVarName string) (s string, ok bool) {
	var (
		sTest string
	)
	sTest, ok = errorTest(got, want)
	if !ok {
		format := "%s %s[%d]:\t\n%s"
		s1 := `Test of type error ocurred in`
		s = fmt.Sprintf(format, s1, tdVarName, i, sTest)
	} else {
		s = ""
	}
	return
}

func byteSliceTest(got []byte, want []byte, i int) (ok bool, s string) {
	var (
		wantStr string
		gotStr  string
	)
	if want == nil {
		wantStr = "<empty>"
	} else {
		wantStr = string(want)
		if len(wantStr) == 0 {
			wantStr = "<empty>"
		}
	}
	if got == nil {
		gotStr = "<empty>"
	} else {
		gotStr = string(got)
		if len(gotStr) == 0 {
			gotStr = "<empty>"
		}
	}
	if strings.Compare(wantStr, gotStr) != 0 {
		format := "%s%s\n, %s%s\"\n%s %d]"
		s1 := `Got the byte slice "`
		s2 := `"... but want the byte slice "`
		s3 := `The test error ocurred in test: tdNew[`
		s = fmt.Sprintf(format, s1, gotStr, s2, wantStr, s3, i)
		ok = false
	} else {
		s = ""
		ok = true
	}
	return
}

const float64TruncatePrecision = 4

func isSameFloat64(a, b float64) (aTxt string, bTxt string, ok bool) {
	aTxt = strconv.FormatFloat(a, 'f', float64TruncatePrecision, 64)
	bTxt = strconv.FormatFloat(b, 'f', float64TruncatePrecision, 64)
	ok = (strings.Compare(aTxt, bTxt) == 0)
	return
}

func isSameFloat64Test(want, got float64) (s string, ok bool) {
	var (
		wantTxt string
		gotTxt  string
	)
	gotTxt = strconv.FormatFloat(got, 'f', float64TruncatePrecision, 64)
	if wantTxt, gotTxt, ok = isSameFloat64(want, got); !ok {
		format1 := "\tWant: %f  (truncated to: %s)."
		format2 := "\n\tGot %f (truncated to: %s)."
		format3 := "is not the same float64 value (precision: %d)"
		format := format1 + format2 + format3
		s := fmt.Sprintf(format,
			want, wantTxt,
			got, gotTxt,
			float64TruncatePrecision)
		return s, false
	}
	return "", true
}

func isSameFloat64TestI(want, got float64, i int, tdVarName string) (s string, ok bool) {
	if s, ok = isSameFloat64Test(want, got); !ok {
		t := fmt.Sprintf("\n%s[%d]:\n%s", tdVarName, i, s)
		return t, ok
	}
	return s, ok
}
