package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestFailBacque(t *testing.T) {
	t.Run("Bacque failure is written with DateTime", func(t *testing.T) {
		buffer := bytes.Buffer{}
		prefix := "DateTime"

		err := failBacque(&buffer)
		assertError(t, err, nil)

		if !strings.HasPrefix(buffer.String(), prefix) {
			t.Errorf("Failure should start with '%s' but got '%s'", prefix, buffer.String())
		}
	})
}

// Integration test: OS date
func TestGetDate(t *testing.T) {
	t.Run("Go time matches time retrieved", func(t *testing.T) {
		// check the function's OS-sourced datetime
		// to Go's internal datetime with to-the-hour accuracy
		want := time.Now().Format("2006010215")
		got, err := getDate("date", "+%Y%m%d%H")

		assertError(t, err, nil)
		assertString(t, got, want)
	})

	t.Run("Error is thrown for failed configuration", func(t *testing.T) {
		// use a typo to produce a failure with the executable
		_, err := getDate("cate", "+%Y%m%d%H")

		assertGotError(t, err)
	})
}
