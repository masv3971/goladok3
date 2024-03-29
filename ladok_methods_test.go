package goladok3

import (
	"context"
	"fmt"
	"testing"

	"github.com/masv3971/goladok3/ladokmocks"
	"github.com/masv3971/goladok3/ladoktypes"
	"github.com/stretchr/testify/assert"
)

func TestCheckPermission(t *testing.T) {
	var (
		uid = "11111111-2222-0000-0000-000000000000"
	)
	type serverStatusCode struct {
		egna, profil int
	}
	type serverReply struct {
		egna, profil []byte
	}
	type serverURL struct {
		egna, profil string
	}
	tts := []struct {
		name             string
		serverStatusCode serverStatusCode
		serverReply      serverReply
		serverURL        serverURL
		have             Permissions
		want             interface{}
	}{
		{
			name: "OK",
			serverURL: serverURL{
				egna:   "/kataloginformation/anvandarbehorighet/egna",
				profil: fmt.Sprintf("/kataloginformation/behorighetsprofil/%s", uid),
			},
			serverStatusCode: serverStatusCode{200, 200},
			serverReply: serverReply{
				egna:   ladokmocks.JSONKataloginformationEgna,
				profil: ladokmocks.JSONKataloginformationBehorighetsprofil,
			},
			have: Permissions{61001: "rattighetsniva.las", 90019: "rattighetsniva.las"},
			want: nil,
		},
		{
			name: "Missing id 0 with permission las",
			serverURL: serverURL{
				egna:   "/kataloginformation/anvandarbehorighet/egna",
				profil: fmt.Sprintf("/kataloginformation/behorighetsprofil/%s", uid),
			},
			have: Permissions{61001: "rattighetsniva.las", 8888: "rattighetsniva.las"},
			want: ladoktypes.PermissionErrors{{
				Msg:                 "Missing ladok permission",
				MissingPermissionID: 8888,
				PermissionLevel:     "rattighetsniva.las",
			}},
			serverStatusCode: serverStatusCode{200, 200},
			serverReply: serverReply{
				egna:   ladokmocks.JSONKataloginformationEgna,
				profil: ladokmocks.JSONKataloginformationBehorighetsprofil,
			},
		},
		{
			name: "Empty input Permissions map",
			serverURL: serverURL{
				egna:   "/kataloginformation/anvandarbehorighet/egna",
				profil: fmt.Sprintf("/kataloginformation/behorighetsprofil/%s", uid),
			},
			have: Permissions{},
			want: ladoktypes.PermissionErrors{
				{
					Msg: "No permissions provided",
				}},
			serverStatusCode: serverStatusCode{200, 200},
			serverReply: serverReply{
				egna:   ladokmocks.JSONKataloginformationEgna,
				profil: ladokmocks.JSONKataloginformationBehorighetsprofil,
			},
		},
		{
			name: "Ladok does not have any permissions",
			serverURL: serverURL{
				egna:   "/kataloginformation/anvandarbehorighet/egna",
				profil: fmt.Sprintf("/kataloginformation/behorighetsprofil/%s", uid),
			},
			have: Permissions{61001: "rattighetsniva.las", 90019: "rattighetsniva.las"},
			want: ladoktypes.PermissionError{
				Msg: "No permission found in ladok",
			},
			serverStatusCode: serverStatusCode{200, 200},
			serverReply: serverReply{
				egna:   ladokmocks.JSONKataloginformationEgna,
				profil: ladokmocks.JSONKataloginformationBehorighetsprofilNoPermissions,
			},
		},
		{
			name: "Egna does not respond",
			serverURL: serverURL{
				egna:   "/kataloginformation/anvandarbehorighet/egna",
				profil: fmt.Sprintf("/kataloginformation/behorighetsprofil/%s", uid),
			},
			have: Permissions{61001: "rattighetsniva.las", 90019: "rattighetsniva.las"},
			want: &ladoktypes.LadokError{
				FelUID:          "c0f52d2c-3a5f-11ec-aa00-acd34b504da7",
				Felkategori:     "commons.fel.kategori.applikationsfel",
				FelkategoriText: "Generellt fel i applikationen",
				Meddelande:      "java.lang.NullPointerException null",
				Link:            []interface{}{},
			},
			serverStatusCode: serverStatusCode{500, 200},
			serverReply: serverReply{
				egna:   ladokmocks.JSONErrors500,
				profil: ladokmocks.JSONKataloginformationBehorighetsprofil,
			},
		},
		{
			name: "Profil does not respond",
			serverURL: serverURL{
				egna:   "/kataloginformation/anvandarbehorighet/egna",
				profil: fmt.Sprintf("/kataloginformation/behorighetsprofil/%s", uid),
			},
			have: Permissions{61001: "rattighetsniva.las", 90019: "rattighetsniva.las"},
			want: &ladoktypes.LadokError{
				Detaljkod:       "commons.domain.uid",
				DetaljkodText:   "Unik identifierare",
				FelUID:          "14c837fd-3a60-11ec-aa00-acd34b504da7",
				Felgrupp:        "commons.fel.grupp.felaktigt_format",
				FelgruppText:    "Felaktigt format",
				Felkategori:     "commons.fel.kategori.valideringsfel",
				FelkategoriText: "Valideringsfel",
				Meddelande:      "Uid: 6daf0d1e-114f-11ec-95ca-f52940734df",
				Link:            []interface{}{},
			},
			serverStatusCode: serverStatusCode{200, 500},
			serverReply: serverReply{
				egna:   ladokmocks.JSONKataloginformationEgna,
				profil: ladokmocks.JSONErrorsValideringsFel,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			mux, server, client := mockSetup(t, ladoktypes.EnvIntTestAPI)
			defer server.Close()

			mockGenericEndpointServer(t, mux, ContentTypeKataloginformationJSON, "GET", tt.serverURL.egna, tt.serverReply.egna, tt.serverStatusCode.egna)
			mockGenericEndpointServer(t, mux, ContentTypeKataloginformationJSON, "GET", tt.serverURL.profil, tt.serverReply.profil, tt.serverStatusCode.profil)

			err := client.CheckPermission(context.TODO(), tt.have)

			switch tt.want.(type) {
			case *ladoktypes.LadokError:
				want := tt.want.(*ladoktypes.LadokError)
				assert.Equal(t, want, err)
			case *ladoktypes.PermissionErrors:
				want := tt.want.(*ladoktypes.PermissionErrors)
				assert.Equal(t, want, err)
			case nil:
				assert.Equal(t, nil, err)
			}
		})
	}
}

func TestComparePermission(t *testing.T) {
	client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "localhost")

	type have struct {
		ladok, my int64
	}
	tts := []struct {
		name string
		have have
		want bool
	}{
		{
			name: "Equal permissions",
			have: have{
				ladok: 6,
				my:    6,
			},
			want: true,
		},
		{
			name: "Ladok require less permissions then what's provided",
			have: have{
				ladok: 4,
				my:    6,
			},
			want: true,
		},
		{
			name: "Ladok require better permissions then what's provided",
			have: have{
				ladok: 6,
				my:    4,
			},
			want: false,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := client.comparePermission(context.TODO(), tt.have.ladok, tt.have.my)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPermissionUnify(t *testing.T) {
	type have struct {
		ladok ladoktypes.KataloginformationBehorighetsprofil
		my    Permissions
	}
	tts := []struct {
		name string
		have have
		want map[int64]map[string]int64
	}{
		{
			name: "Same permissions",
			have: have{
				ladok: ladoktypes.KataloginformationBehorighetsprofil{
					Systemaktiviteter: []ladoktypes.Systemaktiviteter{
						{
							ID:             51001,
							Rattighetsniva: "rattighetsniva.las",
						},
					},
				},
				my: map[int64]string{
					51001: "rattighetsniva.las",
				},
			},
			want: map[int64]map[string]int64{
				51001: {
					"ladok": 4,
					"my":    4,
				},
			},
		},
		{
			name: "Not the same permissions",
			have: have{
				ladok: ladoktypes.KataloginformationBehorighetsprofil{
					Systemaktiviteter: []ladoktypes.Systemaktiviteter{
						{
							ID:             41001,
							Rattighetsniva: "rattighetsniva.lokal",
						},
					},
				},
				my: map[int64]string{
					61001: "rattighetsniva.las",
				},
			},
			want: map[int64]map[string]int64{
				41001: {
					"ladok": 6,
				},
				61001: {
					"my": 4,
				},
			},
		},
		{
			name: "The same permission with different permission levels.",
			have: have{
				ladok: ladoktypes.KataloginformationBehorighetsprofil{
					Systemaktiviteter: []ladoktypes.Systemaktiviteter{

						{
							ID:             81001,
							Rattighetsniva: "rattighetsniva.las",
						},
					},
				},
				my: map[int64]string{
					81001: "rattighetsniva.lokal",
				},
			},
			want: map[int64]map[string]int64{
				81001: {
					"ladok": 4,
					"my":    6,
				},
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "localhost")
			permissions, err := client.permissionUnify(context.TODO(), tt.have.ladok, tt.have.my)
			assert.NoError(t, err)

			assert.Equal(t, tt.want, permissions)
		})
	}
}

func TestTranslateID(t *testing.T) {
	tts := []struct {
		name string
		have int64
		want string
	}{
		{
			name: "OK",
			have: 61001,
			want: "studentinformation.lasa",
		},
		{
			name: "Undefined",
			have: 00000,
			want: "undefined",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			// Test function
			client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "localhost")
			got := client.translateID(tt.have)
			assert.Equal(t, tt.want, got)
		})
	}
}
