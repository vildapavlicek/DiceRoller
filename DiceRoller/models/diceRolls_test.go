package models

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParseHTML(t *testing.T) {
	t.Run("Parse 10 numbers", func(t *testing.T) {
		rolls := DiceRolls{
			Data: []int{},
		}

		status := http.StatusOK
		body := []byte(`<!DOCTYPE html>
		<html>
		<body>
		
		<p>My first paragraph.
		<img src="dice2.png" alt="2" />
		<img src="dice4.png" alt="4" />
		<img src="dice4.png" alt="4" />
		<img src="dice3.png" alt="3" />
		<img src="dice1.png" alt="1" />
		<img src="dice4.png" alt="4" />
		<img src="dice2.png" alt="2" />
		<img src="dice5.png" alt="5" />
		<img src="dice4.png" alt="4" />
		<img src="dice5.png" alt="5" />
		</p>
		
		</body>
		</html> `)

		server := newTestHTTPServer(status, body)
		defer server.Close()

		client := http.Client{}
		response, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Test failed, with error: %s", err)
		}

		rolls.ParseHTML(response)
		got := rolls.Data
		want := []int{2, 4, 4, 3, 1, 4, 2, 5, 4, 5}

		assetrEqualsParsedHTML(t, want, got)

	})

	t.Run("Parse 5 numbers", func(t *testing.T) {
		rolls := DiceRolls{
			Data: []int{},
		}

		status := http.StatusOK
		body := []byte(`<!DOCTYPE html>
		<html>
		<body>
		
		<p>My first paragraph.
		<img src="dice2.png" alt="2" />
		<img src="dice4.png" alt="4" />
		<img src="dice4.png" alt="4" />
		<img src="dice3.png" alt="3" />
		<img src="dice1.png" alt="1" />
		</p>
		
		</body>
		</html> `)

		server := newTestHTTPServer(status, body)
		defer server.Close()

		client := http.Client{}
		response, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Test failed, with error: %s", err)
		}

		rolls.ParseHTML(response)
		got := rolls.Data
		want := []int{2, 4, 4, 3, 1}

		assetrEqualsParsedHTML(t, want, got)

	})
}

func newTestHTTPServer(status int, body []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(body)
	}))
}

func assetrEqualsParsedHTML(t *testing.T, want, got []int) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Parsed result doesn't match expected one, got '%v', want '%v'", got, want)
	}
}
