package goladok3

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestKataloginformation(t *testing.T) {
	var (
		client = mockNewClient(t, envProdAPI, "")
	)
	tts := []struct {
		name       string
		url        string
		payload    []byte
		statusCode int
		param      string
		reply      interface{}
		fn         interface{}
	}{
		{
			name:       "GetAnvandareAutentiserad 200",
			url:        "/kataloginformation/anvandare/autentiserad",
			payload:    jsonAutentiserad,
			statusCode: 200,
			reply:      &AnvandareAutentiserad{},
			param:      "",
			fn:         client.Kataloginformation.GetAnvandareAutentiserad,
		},
		{
			name:       "GetAnvandareAutentiserad 500",
			url:        "/kataloginformation/anvandare/autentiserad",
			payload:    jsonErrors500,
			statusCode: 500,
			reply: &Errors{Ladok: &LadokError{
				FelUID:          "c0f52d2c-3a5f-11ec-aa00-acd34b504da7",
				Felkategori:     "commons.fel.kategori.applikationsfel",
				FelkategoriText: "Generellt fel i applikationen",
				Meddelande:      "java.lang.NullPointerException null",
				Link:            []interface{}{},
			}},
			param: "",
			fn:    client.Kataloginformation.GetAnvandareAutentiserad,
		},
		{
			name:       "GetAnvandarbehorighetEgna 200",
			url:        "/kataloginformation/anvandarbehorighet/egna",
			payload:    jsonEgna,
			statusCode: 200,
			reply:      &KataloginformationAnvandarbehorighetEgna{},
			param:      "",
			fn:         client.Kataloginformation.GetAnvandarbehorighetEgna,
		},
		{
			name:       "GetAnvandarbehorighetEgna 500",
			url:        "/kataloginformation/anvandarbehorighet/egna",
			payload:    jsonErrors500,
			statusCode: 500,
			reply: &Errors{Ladok: &LadokError{
				FelUID:          "c0f52d2c-3a5f-11ec-aa00-acd34b504da7",
				Felkategori:     "commons.fel.kategori.applikationsfel",
				FelkategoriText: "Generellt fel i applikationen",
				Meddelande:      "java.lang.NullPointerException null",
				Link:            []interface{}{},
			}},
			param: "",
			fn:    client.Kataloginformation.GetAnvandarbehorighetEgna,
		},
		{
			name:       "GetBehorighetsprofil 200",
			url:        "/kataloginformation/behorighetsprofil",
			payload:    jsonProfil,
			statusCode: 200,
			reply:      &KataloginformationBehorighetsprofil{},
			param:      uuid.NewString(),
			fn:         client.Kataloginformation.GetBehorighetsprofil,
		},
		{
			name:    "GetBehorighetsprofil 500",
			url:     "/kataloginformation/behorighetsprofil",
			payload: jsonErrors500,
			reply: &Errors{Ladok: &LadokError{
				FelUID:          "c0f52d2c-3a5f-11ec-aa00-acd34b504da7",
				Felkategori:     "commons.fel.kategori.applikationsfel",
				FelkategoriText: "Generellt fel i applikationen",
				Meddelande:      "java.lang.NullPointerException null",
				Link:            []interface{}{},
			}},
			param: uuid.NewString(),
			fn:    client.Kataloginformation.GetBehorighetsprofil,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			mux, server, _ := mockSetup(t, envIntTestAPI)
			client.url = server.URL

			mockGenericEndpointServer(t, mux, contentTypeKataloginformationJSON, "GET", tt.url, tt.param, tt.payload, tt.statusCode)

			err := json.Unmarshal(tt.payload, tt.reply)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			switch tt.fn.(type) {
			case func(context.Context) (*AnvandareAutentiserad, *http.Response, error):
				f := tt.fn.(func(context.Context) (*AnvandareAutentiserad, *http.Response, error))
				switch tt.statusCode {
				case 200:
					reply, _, err := f(context.TODO())
					if !assert.NoError(t, err) {
						t.Fatal(err)
					}

					if !assert.Equal(t, tt.reply, reply, "Should be equal") {
						t.FailNow()
					}
				case 500:
					_, _, err = f(context.TODO())
					assert.Equal(t, err, tt.reply.(*Errors))
				}
			case func(context.Context) (*KataloginformationAnvandarbehorighetEgna, *http.Response, error):
				f := tt.fn.(func(context.Context) (*KataloginformationAnvandarbehorighetEgna, *http.Response, error))
				switch tt.statusCode {
				case 200:
					reply, _, err := f(context.TODO())
					if !assert.NoError(t, err) {
						t.Fatal(err)
					}

					if !assert.Equal(t, tt.reply, reply, "Should be equal") {
						t.FailNow()
					}
				case 500:
					_, _, err = f(context.TODO())
					assert.Equal(t, err, tt.reply.(*Errors))
				}
			case func(context.Context, *GetBehorighetsprofilerCfg) (*KataloginformationBehorighetsprofil, *http.Response, error):
				f := tt.fn.(func(context.Context, *GetBehorighetsprofilerCfg) (*KataloginformationBehorighetsprofil, *http.Response, error))
				switch tt.statusCode {
				case 200:
					reply, _, err := f(context.TODO(), &GetBehorighetsprofilerCfg{UID: tt.param})
					if !assert.NoError(t, err) {
						t.Fatal(err)
					}

					if !assert.Equal(t, tt.reply, reply, "Should be equal") {
						t.FailNow()
					}
				case 500:
					_, _, err = f(context.TODO(), &GetBehorighetsprofilerCfg{UID: tt.param})
					assert.Equal(t, err, tt.reply.(*Errors))
				}
			default:
				t.Fatalf("ERROR No function signature found! %T", tt.fn)
			}

			server.Close() // Close server after each run
		})
	}
}