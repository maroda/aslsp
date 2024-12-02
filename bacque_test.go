package main

import (
	"bytes"
	"testing"
)

// For mocking FindIP
type mockIPConfig struct {
	ExtIPwPort string
}

func (ipc *mockIPConfig) EgressIP() (string, error) {
	return "10.10.10.28", nil
}

func TestFindIP(t *testing.T) {
	mockIP := &mockIPConfig{ExtIPwPort: "0.0.0.0:0"}
	want := "10.10.10.28"
	got, err := FindIP(mockIP)

	assertString(t, got, want)
	assertError(t, err, nil)
}

func TestEgIP(t *testing.T) {
	e := "8.8.8.8:80"     // external address
	want := "10.10.10.28" // known internal interface
	got, err := EgIP(e)

	assertString(t, got, want)
	assertError(t, err, nil)
}

func TestLocalCMD(t *testing.T) {

	// Use LocalCMD to run a very simple shell command
	// The 'jot' command is an old-school BSD command for creating lists.
	// If tests fail because of this command not being present,
	// then that's a dependency we may not want.
	t.Run("simple shell: jot", func(t *testing.T) {
		app := "jot"
		arg := "1"
		buffer := bytes.Buffer{}

		// Use LocalCMD to run a shell command
		//
		// this function only returns an error
		// because it writes to a buffer as its output
		// so make sure to test the error
		err := LocalCMD(&buffer, app, arg)

		// reach into that buffer method where the string is put
		got := buffer.String()
		want := "1"

		assertString(t, got, want)
		assertError(t, err, nil)
	})

	/*
		t.Run("fake date", func(t *testing.T) {
			app := "date"
			arg := "+%Y%m%d%H%M"

			buffer := bytes.Buffer{}
			current := time.Now()

			// Use Golang to get the time
			timeS, werr := fmt.Printf("%d%02d%02d%02d%02d",
				current.Year(), current.Month(), current.Day(),
				current.Hour(), current.Minute())
			if werr != nil {
				t.Errorf("couldn't get time stamp from Go")
			}
			fmt.Println(timeS)

			// Use LocalCMD to get the time
			//
			// this function only returns an error
			// because it writes to a buffer as its output
			// so make sure to test the error
			err := LocalCMD(&buffer, app, arg)

			// reach into that buffer method where the string is put
			got := buffer.String()
			// show what it should be
			// i need to get the value in timeS to be truncated
			// so it can match
			// so right now the test only works with the manual time entered here :D
			want := "202407191645\n"

			if got != want {
				t.Errorf("got %q want %q", got, want)
			}

			// The error should be nil
			assertError(t, err, nil)
		})
	*/
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}

func assertString(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
