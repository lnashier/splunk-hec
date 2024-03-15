package hec

import "time"

type ScanOpt func(*scanOpts)

type scanOpts struct {
	manager       *EventManager
	stdoutEnabled bool
	drainTime     time.Duration
}

func (s *scanOpts) apply(opts []ScanOpt) {
	for _, o := range opts {
		o(s)
	}
}

var defaultScanOpts = scanOpts{
	stdoutEnabled: true,
	drainTime:     time.Duration(10) * time.Second,
}

func WithEventManager(v *EventManager) ScanOpt {
	return func(s *scanOpts) {
		s.manager = v
	}
}

func WithStdoutEnabled(v bool) ScanOpt {
	return func(o *scanOpts) {
		o.stdoutEnabled = v
	}
}

func WithDrainTime(v time.Duration) ScanOpt {
	return func(o *scanOpts) {
		o.drainTime = v
	}
}
