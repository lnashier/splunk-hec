package hec

type ManagerOpt func(*managerOpts)

type managerOpts struct {
	client         *Client
	bufferSize     int
	batchSize      int
	maxWaitSeconds int
}

func (s *managerOpts) apply(opts []ManagerOpt) {
	for _, o := range opts {
		o(s)
	}
}

var defaultManagerOpts = managerOpts{
	bufferSize:     10,
	batchSize:      10,
	maxWaitSeconds: 1,
}

func WithClient(c *Client) ManagerOpt {
	return func(s *managerOpts) {
		s.client = c
	}
}

func WithBufferSize(v int) ManagerOpt {
	return func(o *managerOpts) {
		o.bufferSize = v
	}
}

func WithBatchSize(v int) ManagerOpt {
	return func(o *managerOpts) {
		o.batchSize = v
	}
}

func WithMaxWaitSeconds(v int) ManagerOpt {
	return func(o *managerOpts) {
		o.maxWaitSeconds = v
	}
}
