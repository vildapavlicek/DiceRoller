package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetResponse(t *testing.T) {
	t.Run("200 OK test", func(t *testing.T) {
		want := http.StatusOK
		server := makeStatusServer(want)

		randomRoller := RandomOrgDiceRoller{
			client: &http.Client{},
			genURI: server.URL,
		}

		response, err := randomRoller.GetResponse("GET", nil)
		defer response.Body.Close()
		if err != nil {
			t.Fatalf("Failed to get response: '%s'", err)
		}

		got := response.StatusCode
		assertEqualsStatus(t, got, want)

	})

	t.Run("500 internal server error", func(t *testing.T) {
		status := http.StatusInternalServerError
		server := makeStatusServer(status)
		defer server.Close()

		randomRoller := RandomOrgDiceRoller{
			client: &http.Client{},
			genURI: server.URL,
		}

		_, err := randomRoller.GetResponse("GET", nil)
		//defer response.Body.Close()

		want := "Failed to get reponse 200 OK"
		got := err.Error()

		assertEqualsError(t, got, want)
	})

	t.Run("timeout error, 5s Sleep server, 1s Timeout", func(t *testing.T) {
		timeout := 1 * time.Second
		server := makeSleepServer()
		defer server.Close()

		randomRoller := RandomOrgDiceRoller{
			client: &http.Client{
				Timeout: timeout,
			},
			genURI: server.URL,
		}

		response, err := randomRoller.GetResponse("GET", nil)
		defer response.Body.Close()

		assertErroNotNil(t, err)

	})

}

func makeStatusServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))
}

func makeSleepServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
}

func assertEqualsStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("Got '%v', want '%v'", got, want)
	}
}

func assertEqualsError(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("Got '%v', want '%v'", got, want)
	}
}

func assertErroNotNil(t *testing.T, got error) {
	t.Helper()
	if got == nil {
		t.Fatalf("Wanted error, but got nil")
	}
}
