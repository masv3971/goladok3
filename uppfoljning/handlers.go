package uppfoljning

import (
	"context"
	"fmt"
	"net/http"
)

// Service holds uppfoljning feed object
type Service struct {
	client      *Service
	contentType string
}

// FeedRecent atom feed /uppfoljning/feed/recent
func (s *Service) FeedRecent(ctx context.Context) (*SuperFeed, *http.Response, error) {
	env, err := s.client.environment()
	if err != nil {
		return nil, nil, err
	}

	var url string
	switch env {
	case envIntTestAPI:
		url = "/handelse/feed/recent"
	default:
		url = "/uppfoljning/feed/recent"
	}

	req, err := s.client.newRequest(
		ctx,
		"GET",
		fmt.Sprintf("%s", url),
		ladokAcceptHeader[s.contentType]["xml"],
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	reply := &FeedRecent{}
	resp, err := s.client.do(req, reply)
	if err != nil {
		return nil, resp, err
	}

	superFeed, err := reply.parse()
	if err != nil {
		return nil, nil, err
	}

	return superFeed, resp, nil
}