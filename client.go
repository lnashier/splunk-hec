package hec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lnashier/goarc/x/env"
	xhttp "github.com/lnashier/goarc/x/http"
	"net/http"
)

type Client struct {
	hc         *xhttp.Client
	token      string
	appHost    string
	index      string
	source     string
	sourceType string
}

func NewClient(opt ...ClientOpt) *Client {
	opts := defaultClientOpts
	opts.apply(opt)

	return &Client{
		hc: xhttp.NewClient(
			xhttp.WithHost(opts.host),
			xhttp.WithTimeout(10),
			xhttp.WithEncoder(json.Marshal),
			xhttp.WithDecoder(json.Unmarshal),
		),
		token:      opts.token,
		appHost:    env.Hostname(),
		index:      opts.index,
		source:     opts.source,
		sourceType: opts.sourceType,
	}
}

func (c Client) Send(ctx context.Context, events ...any) (map[string]any, error) {
	if len(events) < 1 {
		return nil, nil
	}

	var payloads eventPayloads
	for _, event := range events {
		payloads = append(payloads, &eventPayload{
			Host:       c.appHost,
			Index:      c.index,
			Source:     c.source,
			SourceType: c.sourceType,
			Event:      event,
		})
	}

	encodedPayload, err := payloads.Marshal()
	if err != nil {
		return nil, fmt.Errorf("error marshaling to events: %v", err)
	}

	header := make(http.Header)
	header.Set("Authorization", fmt.Sprintf("Splunk %s", c.token))

	req, err := c.hc.NewRequest(ctx, http.MethodPost, "/services/collector/event", header, bytes.NewReader(encodedPayload))
	if err != nil {
		return nil, fmt.Errorf("error building event api request: %v", err)
	}

	result := make(map[string]any)
	_, err = c.hc.DoDecoded(req, &result, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling event api: %v", err)
	}
	return result, nil
}

type eventPayloads []*eventPayload

func (ps eventPayloads) Marshal() ([]byte, error) {
	var pbs []byte
	for _, p := range ps {
		pb, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}
		pbs = append(pbs, pb...)
		pbs = append(pbs, []byte("\n")...)
	}
	return pbs, nil
}

type eventPayload struct {
	Host       string `json:"host"`
	Index      string `json:"index"`
	Source     string `json:"source"`
	SourceType string `json:"sourcetype"`
	Event      any    `json:"event"`
}
