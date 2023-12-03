package intgrt_afick

import (
	"fmt"
	asrt "github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCheckOutput(t *testing.T) {
	s := "# Hash database : 14 files scanned, 0 changed (new : 1; delete : 2; changed : 3; dangling : 4; exclude_suffix : 5; exclude_prefix : 6; exclude_re : 7; degraded : 8)"
	exp := &AfickCheckRes{Scanned: 14,
		New:            1,
		Delete:         2,
		Changed:        3,
		Dangling:       4,
		Exclude_suffix: 5,
		Exclude_prefix: 6,
		Exclude_re:     7,
		Degraded:       8}
	res, _ := parseCheckOutput([]byte(s))
	fmt.Printf("Parsed result:\t %#v \n expected: \t %#v", *res, *exp)
	asrt.Equal(t, *exp, *res, "parsed result not equal to expected")
}
