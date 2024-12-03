package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Used here for OS integration tests
var (
	localIP  = "10.10.10.95"
	googleIP = "8.8.8.8:80"
)

// For mocking FindIP
type mockIPConfig struct {
	ExtIPwPort string
}

func (ipc *mockIPConfig) EgressIP() (string, error) {
	return localIP, nil
}

func TestFindIP(t *testing.T) {
	t.Run("Matches an IP", func(t *testing.T) {
		mockIP := &mockIPConfig{ExtIPwPort: "0.0.0.0:0"}
		got, err := FindIP(mockIP)

		assertString(t, got, localIP)
		assertError(t, err, nil)
	})
}

// Integration test: TCP/IP
func TestIPConfig_EgressIP(t *testing.T) {
	t.Run("Retrieves the egress IP address", func(t *testing.T) {
		// interface struct
		mockIP := &mockIPConfig{ExtIPwPort: googleIP}
		got, err := mockIP.EgressIP()
		assertString(t, got, localIP)
		assertError(t, err, nil)
	})
}

// Integration test: TCP/IP
func TestEgIP(t *testing.T) {
	t.Run("Retrieves the egress IP address", func(t *testing.T) {
		got, err := EgIP(googleIP)
		assertString(t, got, localIP)
		assertError(t, err, nil)
	})

	t.Run("Throws an error with no port given", func(t *testing.T) {
		e := "8.8.8.8" // any address minus the port
		_, err := EgIP(e)
		assertGotError(t, err)
	})
}

// TODO: This is problematic because of the 'jot' command
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
}

// Integration test: Network
func TestRequestIP(t *testing.T) {
	t.Run("Returns the httptest IP", func(t *testing.T) {
		// This is in the network created by httptest
		want := "RequestIP=192.0.2.1\n"

		r := httptest.NewRequest(http.MethodGet, "/dt", nil)
		w := httptest.NewRecorder()
		err := RequestIP(w, r)

		assertError(t, err, nil)
		assertStatus(t, w.Code, http.StatusOK)
		assertResponseBody(t, w.Body.String(), want)

		err = r.Body.Close()
		assertError(t, err, nil)
	})

	t.Run("Returns an error for no port", func(t *testing.T) {
		// mock client for extracting a request IP address
		r := httptest.NewRequest(http.MethodGet, "/dt", nil)
		// Set the response RemoteAddr to something without a port
		r.RemoteAddr = "1.2.3.4"
		w := httptest.NewRecorder()
		got := RequestIP(w, r)

		assertGotError(t, got)

		err := r.Body.Close()
		assertError(t, err, nil)
	})
}

func TestPing(t *testing.T) {
	want := "pong"

	// mock client to ping()
	r := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	ping(w, r)

	assertStatus(t, w.Code, http.StatusOK)
	assertResponseBody(t, w.Body.String(), want)
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("got error %q want %q", got, want)
	}
}

func assertGotError(t testing.TB, got error) {
	t.Helper()
	if errors.Is(got, nil) {
		t.Errorf("expected an error but got nil (%q)", got)
	}
}

func assertString(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
