package hec

type ClientOpt func(*clientOpts)

type clientOpts struct {
	host       string
	token      string
	index      string
	source     string
	sourceType string
}

func (s *clientOpts) apply(opts []ClientOpt) {
	for _, o := range opts {
		o(s)
	}
}

var defaultClientOpts = clientOpts{
	sourceType: "_json",
}

func WithHost(host string) ClientOpt {
	return func(s *clientOpts) {
		s.host = host
	}
}

func WithToken(t string) ClientOpt {
	return func(o *clientOpts) {
		o.token = t
	}
}

func WithIndex(i string) ClientOpt {
	return func(o *clientOpts) {
		o.index = i
	}
}

func WithSource(s string) ClientOpt {
	return func(o *clientOpts) {
		o.source = s
	}
}

func WithSourceType(s string) ClientOpt {
	return func(o *clientOpts) {
		o.sourceType = s
	}
}
