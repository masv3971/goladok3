package goladok3

import (
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
)

type feedService struct {
	client  *Client
	service string
}

func (s *feedService) acceptHeader() string {
	return ladokAcceptHeader[s.service]["xml"]
}

type feedRecent struct {
	XMLName xml.Name `xml:"feed"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Title   struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"title"`
	Link []struct {
		Text string `xml:",chardata"`
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
	ID        FeedID `xml:"id"`
	Generator struct {
		Text string `xml:",chardata"`
		URI  string `xml:"uri,attr"`
	} `xml:"generator"`
	Updated string `xml:"updated"`
	Entry   []struct {
		Text     string `xml:",chardata"`
		Category struct {
			Text  string `xml:",chardata"`
			Term  string `xml:"term,attr"`
			Label string `xml:"label,attr"`
		} `xml:"category"`
		ID      string `xml:"id"`
		Updated string `xml:"updated"`
		Content struct {
			Text                             string                 `xml:",chardata"`
			Type                             string                 `xml:"type,attr"`
			AnvandareAndradEvent             *anvandareEvent        `xml:"AnvandareAndradEvent,omitempty"`
			AnvandareSkapadEvent             *anvandareEvent        `xml:"AnvandareSkapadEvent,omitempty"`
			KontaktuppgifterEvent            *kontaktuppgifterEvent `xml:"KontaktuppgifterEvent,omitempty"`
			ExternPartEvent                  *externPartEvent       `xml:"ExternPartEvent,omitempty"`
			LokalStudentEvent                *lokalStudentEvent     `xml:"LokalStudentEvent,omitempty"`
			ResultatPaModulAttesteratEvent   *resultatEvent         `xml:"ResultatPaModulAttesteratEvent,omitempty"`
			ResultatPaHelKursAttesteratEvent *resultatEvent         `xml:"ResultatPaHelKursAttesteratEvent,omitempty"`
		} `xml:"content"`
	} `xml:"entry"`
}

type eventContext struct {
	Text         string `xml:",chardata"`
	AnvandareUID string `xml:"AnvandareUID"`
	Anvandarnamn string `xml:"Anvandarnamn"`
	LarosateID   string `xml:"LarosateID"`
}

// SuperFeed is a made up type, in order to make unstructured data to structured data.
type SuperFeed struct {
	ID          int           `json:"id,omitempty"`
	SuperEvents []*SuperEvent `json:"super_events,omitempty"`
}

// SuperEventContext is a made up type
type SuperEventContext struct {
	AnvandareUID string `json:"anvandare_uid"`
	Anvandarnamn string `json:"anvandarnamn"`
	LarosateID   string `json:"larosate_id"`
}

// SuperPostadress is a made up type
type SuperPostadress struct {
	Land             string `json:"land"`
	PostadressTyp    string `json:"postadress_typ"`
	Postnummer       string `json:"postnummer"`
	Postort          string `json:"postort"`
	Utdelningsadress string `json:"utdelningsadress"`
	CareOf           string `json:"care_of"`
}

// SuperBeslut is a made up type
type SuperBeslut struct {
	BeslutUID         string `json:"beslut_uid"`
	Beslutsdatum      string `json:"beslutsdatum"`
	Beslutsfattare    string `json:"beslutsdattare"`
	BeslutsfattareUID string `json:"beslutsfattare_uid"`
}

// SuperResultat is a made up type
type SuperResultat struct {
	BetygsgradID       string `json:"betygsgrad_id"`
	BetygsskalaID      string `json:"betygsskala_id"`
	Examinationsdatum  string `json:"examinationsdatum"`
	GiltigSomSlutbetyg string `json:"giltig_som_slutbetyg"`
	OmfattningsPoang   string `json:"omfattnings_poang"`
	PrestationsPoang   string `json:"prestations_poang"`
	ResultatUID        string `json:"resultat_uid"`
}

// SuperEvent is a made up type consists of all the aviable ladok attributes
type SuperEvent struct {
	EventTypeName         string            `json:"event_type_name"`
	EntryID               string            `json:"entry_id"`
	EventContext          SuperEventContext `json:"event_context"`
	HandelseUID           string            `json:"handelse_uid"`
	ID                    string            `json:"id"`
	AnvandareUID          string            `json:"anvandare_uid"`
	Efternamn             string            `json:"efternamn"`
	Email                 string            `json:"email"`
	Fornamn               string            `json:"fornamn"`
	Handelsetyp           string            `json:"handelsetype"`
	StudentUID            string            `json:"student_uid"`
	Postadresser          []SuperPostadress `json:"postadresser"`
	Telefonnummer         string            `json:"telefonnummer"`
	Beslut                SuperBeslut       `json:"beslut"`
	Resultat              SuperResultat     `json:"resultat"`
	UtbildningsinstansUID string            `json:"utbildningsinstans_uid"`
	Anvandarnamnet        string            `json:"anvandarnamnet"`
	EventTyp              string            `json:"event_typ"`
	Giltighetsperiod      string            `json:"giltighetsperiod"`
	Kod                   string            `json:"kod"`
	LandID                string            `json:"land_id"`
	TypAvExternPartID     string            `json:"typ_av_extern_part_id"`
	KursUID               string            `json:"kurs_uid"`
	KursinstansUID        string            `json:"kursinstans_uid"`
	KurstillfalleUID      string            `json:"kurstillfalle_uid"`
	ExterntStudentUID     string            `json:"externt_student_uid"`
	Fodelsedata           string            `json:"fodelsedata"`
	Kon                   string            `json:"kon"`
	Personnummer          string            `json:"personnummer"`
}

type anvandareEvent struct {
	Text           string       `xml:",chardata"`
	Ki             string       `xml:"ki,attr"`
	Base           string       `xml:"base,attr"`
	Dap            string       `xml:"dap,attr"`
	Events         string       `xml:"events,attr"`
	HandelseUID    string       `xml:"HandelseUID"`
	EventContext   eventContext `xml:"EventContext"`
	AnvandareUID   string       `xml:"AnvandareUID"`
	Anvandarnamnet string       `xml:"Anvandarnamnet"`
	Efternamn      string       `xml:"Efternamn"`
	Fornamn        string       `xml:"Fornamn"`
	Email          string       `xml:"Email"`
}

func (a *anvandareEvent) parse(eventTypeName, entryID string) *SuperEvent {
	s := &SuperEvent{
		EntryID:       entryID,
		EventTypeName: eventTypeName,
		HandelseUID:   a.HandelseUID,
		EventContext: SuperEventContext{
			AnvandareUID: a.EventContext.AnvandareUID,
			Anvandarnamn: a.EventContext.Anvandarnamn,
			LarosateID:   a.EventContext.LarosateID,
		},
		AnvandareUID:   a.AnvandareUID,
		Anvandarnamnet: a.Anvandarnamnet,
		Efternamn:      a.Efternamn,
		Fornamn:        a.Fornamn,
	}
	return s
}

type kontaktuppgifterEvent struct {
	Text         string       `xml:",chardata"`
	Si           string       `xml:"si,attr"`
	Base         string       `xml:"base,attr"`
	Dap          string       `xml:"dap,attr"`
	Events       string       `xml:"events,attr"`
	HandelseUID  string       `xml:"HandelseUID"`
	EventContext eventContext `xml:"EventContext"`
	Handelsetyp  string       `xml:"Handelsetyp"`
	Epostadress  string       `xml:"Epostadress"`
	Postadresser []struct {
		Text             string `xml:",chardata"`
		Land             string `xml:"Land"`
		PostadressTyp    string `xml:"PostadressTyp"`
		Postnummer       string `xml:"Postnummer"`
		Postort          string `xml:"Postort"`
		Utdelningsadress string `xml:"Utdelningsadress"`
		CareOf           string `xml:"CareOf"`
	} `xml:"Postadresser"`
	StudentUID    string `xml:"StudentUID"`
	Telefonnummer string `xml:"Telefonnummer"`
}

func (e *kontaktuppgifterEvent) parse(entryID string) *SuperEvent {
	superAdresser := []SuperPostadress{}

	for _, adress := range e.Postadresser {
		ad := SuperPostadress{
			Land:             adress.Land,
			PostadressTyp:    adress.PostadressTyp,
			Postnummer:       adress.Postnummer,
			Postort:          adress.Postort,
			Utdelningsadress: adress.Utdelningsadress,
			CareOf:           adress.CareOf,
		}
		superAdresser = append(superAdresser, ad)
	}

	s := &SuperEvent{
		EventTypeName: KontaktuppgifterEventName,
		EntryID:       entryID,
		HandelseUID:   e.HandelseUID,
		EventContext: SuperEventContext{
			AnvandareUID: e.EventContext.AnvandareUID,
			Anvandarnamn: e.EventContext.Anvandarnamn,
			LarosateID:   e.EventContext.LarosateID,
		},
		Handelsetyp:   e.Handelsetyp,
		Email:         e.Epostadress,
		StudentUID:    e.StudentUID,
		Telefonnummer: e.Telefonnummer,
		Postadresser:  superAdresser,
	}
	return s
}

type lokalStudentEvent struct {
	Text         string `xml:",chardata"`
	Si           string `xml:"si,attr"`
	Base         string `xml:"base,attr"`
	Dap          string `xml:"dap,attr"`
	Events       string `xml:"events,attr"`
	HandelseUID  string `xml:"HandelseUID"`
	EventContext struct {
		Text         string `xml:",chardata"`
		AnvandareUID string `xml:"AnvandareUID"`
		Anvandarnamn string `xml:"Anvandarnamn"`
		LarosateID   string `xml:"LarosateID"`
	} `xml:"EventContext"`
	Handelsetyp       string `xml:"Handelsetyp"`
	Efternamn         string `xml:"Efternamn"`
	ExterntStudentUID string `xml:"ExterntStudentUID"`
	Fodelsedata       string `xml:"Fodelsedata"`
	Fornamn           string `xml:"Fornamn"`
	Kon               string `xml:"Kon"`
	Personnummer      string `xml:"Personnummer"`
	StudentUID        string `xml:"StudentUID"`
}

func (e *lokalStudentEvent) parse(entryID string) *SuperEvent {
	s := &SuperEvent{
		EntryID:       entryID,
		EventTypeName: LokalStudentEventName,
		EventContext: SuperEventContext{
			AnvandareUID: e.EventContext.AnvandareUID,
			Anvandarnamn: e.EventContext.Anvandarnamn,
			LarosateID:   e.EventContext.LarosateID,
		},
		HandelseUID:       e.HandelseUID,
		Efternamn:         e.Efternamn,
		Fornamn:           e.Fornamn,
		Handelsetyp:       e.Handelsetyp,
		StudentUID:        e.StudentUID,
		ExterntStudentUID: e.ExterntStudentUID,
		Fodelsedata:       e.Fodelsedata,
		Kon:               e.Kon,
		Personnummer:      e.Personnummer,
	}
	return s
}

type externPartEvent struct {
	Text         string       `xml:",chardata"`
	Ki           string       `xml:"ki,attr"`
	Base         string       `xml:"base,attr"`
	Dap          string       `xml:"dap,attr"`
	Events       string       `xml:"events,attr"`
	HandelseUID  string       `xml:"HandelseUID"`
	EventContext eventContext `xml:"EventContext"`
	Benamningar  struct {
		Text      string `xml:",chardata"`
		Benamning []struct {
			Chardata string `xml:",chardata"`
			Sprakkod string `xml:"Sprakkod"`
			Text     string `xml:"Text"`
		} `xml:"Benamning"`
	} `xml:"Benamningar"`
	Beskrivningar struct {
		Text      string `xml:",chardata"`
		Benamning struct {
			Chardata string `xml:",chardata"`
			Sprakkod string `xml:"Sprakkod"`
			Text     string `xml:"Text"`
		} `xml:"Benamning"`
	} `xml:"Beskrivningar"`
	EventTyp          string `xml:"EventTyp"`
	Giltighetsperiod  string `xml:"Giltighetsperiod"`
	ID                string `xml:"Id"`
	Kod               string `xml:"Kod"`
	LandID            string `xml:"LandID"`
	TypAvExternPartID string `xml:"TypAvExternPartID"`
}

func (e *externPartEvent) parse(entryID string) *SuperEvent {
	s := &SuperEvent{
		EntryID:       entryID,
		EventTypeName: "ExternPartEvent",
		HandelseUID:   e.HandelseUID,
		EventContext: SuperEventContext{
			AnvandareUID: e.EventContext.AnvandareUID,
			Anvandarnamn: e.EventContext.Anvandarnamn,
			LarosateID:   e.EventContext.LarosateID,
		},
		EventTyp:          e.EventTyp,
		Giltighetsperiod:  e.Giltighetsperiod,
		ID:                e.ID,
		Kod:               e.Kod,
		LandID:            e.LandID,
		TypAvExternPartID: e.TypAvExternPartID,
	}
	return s
}

type resultatEvent struct {
	Text         string `xml:",chardata"`
	Rr           string `xml:"rr,attr"`
	Base         string `xml:"base,attr"`
	Dap          string `xml:"dap,attr"`
	Events       string `xml:"events,attr"`
	HandelseUID  string `xml:"HandelseUID"`
	EventContext struct {
		Text         string `xml:",chardata"`
		AnvandareUID string `xml:"AnvandareUID"`
		Anvandarnamn string `xml:"Anvandarnamn"`
		LarosateID   string `xml:"LarosateID"`
	} `xml:"EventContext"`
	Beslut struct {
		Text              string `xml:",chardata"`
		BeslutUID         string `xml:"BeslutUID"`
		Beslutsdatum      string `xml:"Beslutsdatum"`
		Beslutsfattare    string `xml:"Beslutsfattare"`
		BeslutsfattareUID string `xml:"BeslutsfattareUID"`
	} `xml:"Beslut"`
	KursUID          string `xml:"KursUID"`
	KursinstansUID   string `xml:"KursinstansUID"`
	KurstillfalleUID string `xml:"KurstillfalleUID"`
	Resultat         struct {
		Text               string `xml:",chardata"`
		BetygsgradID       string `xml:"BetygsgradID"`
		BetygsskalaID      string `xml:"BetygsskalaID"`
		Examinationsdatum  string `xml:"Examinationsdatum"`
		GiltigSomSlutbetyg string `xml:"GiltigSomSlutbetyg"`
		OmfattningsPoang   string `xml:"OmfattningsPoang"`
		PrestationsPoang   string `xml:"PrestationsPoang"`
		ResultatUID        string `xml:"ResultatUID"`
	} `xml:"Resultat"`
	StudentUID            string `xml:"StudentUID"`
	UtbildningsinstansUID string `xml:"UtbildningsinstansUID"`
}

func (e *resultatEvent) parse(eventTypeName, entryID string) *SuperEvent {
	s := &SuperEvent{
		EntryID:       entryID,
		EventTypeName: eventTypeName,
		HandelseUID:   e.HandelseUID,
		EventContext: SuperEventContext{
			AnvandareUID: e.EventContext.AnvandareUID,
			Anvandarnamn: e.EventContext.Anvandarnamn,
			LarosateID:   e.EventContext.LarosateID,
		},
		Beslut: SuperBeslut{
			BeslutUID:         e.Beslut.BeslutUID,
			Beslutsdatum:      e.Beslut.Beslutsdatum,
			Beslutsfattare:    e.Beslut.Beslutsfattare,
			BeslutsfattareUID: e.Beslut.BeslutsfattareUID,
		},
		KursUID:          e.KursUID,
		KursinstansUID:   e.KursinstansUID,
		KurstillfalleUID: e.KurstillfalleUID,
		Resultat: SuperResultat{
			BetygsgradID:       e.Resultat.BetygsgradID,
			BetygsskalaID:      e.Resultat.BetygsskalaID,
			Examinationsdatum:  e.Resultat.Examinationsdatum,
			GiltigSomSlutbetyg: e.Resultat.GiltigSomSlutbetyg,
			OmfattningsPoang:   e.Resultat.OmfattningsPoang,
			PrestationsPoang:   e.Resultat.PrestationsPoang,
			ResultatUID:        e.Resultat.ResultatUID,
		},
		StudentUID:            e.StudentUID,
		UtbildningsinstansUID: e.UtbildningsinstansUID,
	}
	return s
}

func (s *feedService) feedURL() (string, error) {
	env, err := s.client.environment()
	if err != nil {
		return "", err
	}

	switch env {
	case envIntTestAPI:
		return "handelser/feed", nil
	default:
		return "uppfoljning/feed", nil
	}
}

func (s *feedService) atomReader(ctx context.Context, param string) (*SuperFeed, *http.Response, error) {
	envURL, err := s.feedURL()
	if err != nil {
		return nil, nil, err
	}

	reply := &feedRecent{}
	resp, err := s.client.call(ctx, s.acceptHeader(), http.MethodGet, envURL, param, nil, reply)
	if err != nil {
		return nil, resp, err
	}

	superFeed, err := reply.parse()
	if err != nil {
		return nil, nil, err
	}

	return superFeed, resp, nil

}

// Recent atom feed .../feed/recent gets the most recent publiced feed
func (s *feedService) Recent(ctx context.Context) (*SuperFeed, *http.Response, error) {
	superFeed, resp, err := s.atomReader(ctx, "recent")
	if err != nil {
		return nil, resp, err
	}

	return superFeed, resp, nil
}

// Historical atom feed .../feed/{id} gets feed of {id}
func (s *feedService) Historical(ctx context.Context, id int) (*SuperFeed, *http.Response, error) {
	superFeed, resp, err := s.atomReader(ctx, strconv.Itoa(id))
	if err != nil {
		return nil, resp, err
	}

	return superFeed, resp, nil

}

// First atom feed .../feed/first gets the first publiced feed
func (s *feedService) First(ctx context.Context) (*SuperFeed, *http.Response, error) {
	superFeed, resp, err := s.atomReader(ctx, "first")
	if err != nil {
		return nil, resp, err
	}

	return superFeed, resp, nil

}

func (f *feedRecent) parse() (*SuperFeed, error) {
	superFeed := &SuperFeed{}
	feedID, err := f.ID.trim().int()
	if err != nil {
		return nil, err
	}
	superFeed.ID = feedID

	for _, entry := range f.Entry {
		if entry.Content.AnvandareAndradEvent != nil {
			event := entry.Content.AnvandareAndradEvent.parse(AnvandareAndradEventName, entry.ID)
			superFeed.SuperEvents = append(superFeed.SuperEvents, event)
			continue
		}

		if entry.Content.AnvandareSkapadEvent != nil {
			event := entry.Content.AnvandareSkapadEvent.parse(AnvandareSkapadEventName, entry.ID)
			superFeed.SuperEvents = append(superFeed.SuperEvents, event)
			continue
		}

		if entry.Content.ExternPartEvent != nil {
			event := entry.Content.ExternPartEvent.parse(entry.ID)
			superFeed.SuperEvents = append(superFeed.SuperEvents, event)
			continue
		}

		if entry.Content.KontaktuppgifterEvent != nil {
			event := entry.Content.KontaktuppgifterEvent.parse(entry.ID)
			superFeed.SuperEvents = append(superFeed.SuperEvents, event)
			continue
		}

		if entry.Content.ResultatPaModulAttesteratEvent != nil {
			event := entry.Content.ResultatPaModulAttesteratEvent.parse(ResultatPaModulAttesteratEventName, entry.ID)
			superFeed.SuperEvents = append(superFeed.SuperEvents, event)
			continue
		}

		if entry.Content.ResultatPaHelKursAttesteratEvent != nil {
			event := entry.Content.ResultatPaHelKursAttesteratEvent.parse(ResultatPaHelKursAttesteratEventName, entry.ID)
			superFeed.SuperEvents = append(superFeed.SuperEvents, event)
			continue
		}

		if entry.Content.LokalStudentEvent != nil {
			event := entry.Content.LokalStudentEvent.parse(entry.ID)
			superFeed.SuperEvents = append(superFeed.SuperEvents, event)
			continue
		}
	}

	return superFeed, nil
}
