package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DarioAF/ip-analyser/pkg/db"
	"github.com/DarioAF/ip-analyser/pkg/model"
	redis "github.com/go-redis/redis/v8"
)

func TestHealthHandlerWhenUP(t *testing.T) {
	var database db.DBInterface = &db.MockDB{PingMock: "PONG"}

	req := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req, database)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Errorf("health check must be 200, it was: %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("content type must be application/json, it was: %s", ct)
	}

	if string(body) != fmt.Sprint(`{"status":"UP & Running"}`) {
		t.Errorf("response for health must be UP & Running when redis is running, it was: %s", string(body))
	}
}

func TestHealthHandlerWhenDOWN(t *testing.T) {
	var database db.DBInterface = &db.MockDB{PingMock: ""}

	req := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req, database)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	st := resp.StatusCode
	if st != 500 {
		t.Errorf("health check must be 500, it was: %d", st)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("content type must be application/json, it was: %s", ct)
	}

	if string(body) != fmt.Sprint(`{"status":"something is wrong with the db, please check your connection with it"}`) {
		t.Errorf("response for health check was different: %s", body)
	}
}

func TestDistanceHandler(t *testing.T) {
	var database db.DBInterface = &db.MockDB{
		RetrieveMock: func(hash, key string) string {
			if hash != "statistics-AR" {
				t.Errorf("trying to access a unknown hash for distance statistics: %s", key)
				return ""
			}
			if key == "nearest" {
				return "BR-2821"
			}
			return "ES-10274"
		},
	}

	scenarios := map[string]model.Statistic{
		"nearest": model.Statistic{
			Country:  "BR",
			Distance: 2821,
		},
		"farthest": model.Statistic{
			Country:  "ES",
			Distance: 10274,
		},
	}

	for impl, expected := range scenarios {
		req := httptest.NewRequest("GET", "http://localhost:8080/"+impl, nil)
		w := httptest.NewRecorder()

		distanceHandler(w, req, database, impl)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		st := resp.StatusCode
		if st != 200 {
			t.Errorf("distance status code must be 200, it was: %d", st)
		}

		ct := resp.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("content type must be application/json, it was: %s", ct)
		}

		actualResponse, err := json.Marshal(expected)
		if err != nil {
			t.Errorf("There was an error marshaling our user! %err", err)
		}

		if string(body) != string(actualResponse) {
			t.Errorf("response for %s distance was different: %s", impl, body)
		}
	}
}

func TestCountryRequestsHandler(t *testing.T) {
	var database db.DBInterface = &db.MockDB{
		RetrieveAllScoresMock: func(key string) []redis.Z {
			if key == "trend-BR" {
				return []redis.Z{
					{Score: 100, Member: "2.2.2.2"},
				}
			}
			if key == "trend-ES" {
				return []redis.Z{
					{Score: 50, Member: "3.3.3.3"},
					{Score: 30, Member: "4.4.4.4"},
				}
			}
			return []redis.Z{}
		},
	}

	scenarios := map[string]int{
		"BR": 100,
		"ES": 40,
	}

	for country, expected := range scenarios {
		req := httptest.NewRequest("GET", "http://localhost:8080/avg-requests/"+country, nil)
		w := httptest.NewRecorder()

		countryRequestsHandler(w, req, database)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		st := resp.StatusCode
		if st != 200 {
			t.Errorf("avg-requests status must be 200, it was: %d", st)
		}

		ct := resp.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("content type must be application/json, it was: %s", ct)
		}

		if string(body) != fmt.Sprintf(`{"avg":%d}`, expected) {
			t.Errorf("response for /avg-requests/%s was %s, and we where expecting %d", country, body, expected)
		}
	}
}

func TestIPValidations(t *testing.T) {
	var database db.DBInterface = &db.MockDB{}

	scenarios := map[string]string{
		"aaa":                           `{"error":"invalid input: invalid character 'a' looking for beginning of value"}`,
		`{"ip": "9999.9999.9999.9999"}`: `{"error":"invalid ip"}`,
	}

	for post, expected := range scenarios {
		req := httptest.NewRequest("POST", "http://localhost:8080/user", ioutil.NopCloser(strings.NewReader(post)))
		w := httptest.NewRecorder()

		userHandler(w, req, database)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		st := resp.StatusCode
		if st != 400 {
			t.Errorf("avg-requests status must be 400, it was: %d", st)
		}

		ct := resp.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("content type must be application/json, it was: %s", ct)
		}

		if string(body) != expected {
			t.Errorf("response for invalid json was different from expected")
		}
	}
}
