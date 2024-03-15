package hec

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
)

// Scan reads input from stdin and publishes it to a specified Splunk HTTP Event Collector (HEC) and also optionally to stdout.
// See EventManager Client
func Scan(opt ...ScanOpt) error {
	opts := defaultScanOpts
	opts.apply(opt)

	em := opts.manager

	if err := em.Start(); err != nil {
		return err
	}

	// Defer stopping the manager after a specified drain time.
	defer func() {
		time.Sleep(time.Duration(math.Max(10, opts.drainTime.Seconds())))
		em.Stop()
	}()

	publishFn := func(e string) {}

	if opts.stdoutEnabled {
		publishFn = func(e string) {
			fmt.Println(e)
		}
	}

	publishFn1 := publishFn
	publishFn = func(e string) {
		publishFn1(e)
		ee := make(map[string]any)
		if err := json.Unmarshal([]byte(e), &ee); err == nil {
			em.Publish(ee)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		publishFn(scanner.Text())
	}

	return scanner.Err()
}
