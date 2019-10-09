package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	parallel *int
	ignore   *bool
)

func checkSingle(client *http.Client, lnk string) error {
	u, err := url.Parse(lnk)
	if err != nil {
		return fmt.Errorf("parse %q failed with err %w", lnk, err)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("create request for %q failed with err %w", lnk, err)
	}

	rsp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch %q failed with err %w", lnk, err)
	}

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch %q failed with status code %d", lnk, rsp.StatusCode)
	}

	return nil
}

func checkLink(wg *sync.WaitGroup, in chan string, out chan error) {
	defer wg.Done()
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	for lnk := range in {
		if err := checkSingle(client, lnk); err != nil {
			out <- err
		}
	}
}

func validateLinks(_ *command, fl *File) error {
	if *parallel < 1 {
		*parallel = 3
	}

	in := make(chan string, *parallel)
	out := make(chan error, *parallel)
	wg := &sync.WaitGroup{}
	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go checkLink(wg, in, out)
	}

	go func() {
		for i := range fl.Events {
			if len(fl.Events[i].Sources) == 0 {
				log.Printf("The event with key %s has no source\n", fl.Events[i].PartialKey)
				continue
			}

			for _, lnk := range fl.Events[i].Sources {
				in <- lnk
			}
		}

		close(in)
		wg.Wait()
		close(out)
	}()

	var failed bool
	for err := range out {
		if err != nil {
			failed = true
			log.Printf("Link validation failed with err: %q", err)
		}
	}

	if failed && !*ignore {
		return errors.New("validation failed, check the log")
	}
	return nil
}

func init() {
	cmd := newCommand("validate-links", "validate links for the events", validateLinks)
	parallel = cmd.Flags.Int("parallel", 3, "how many parallel check")
	ignore = cmd.Flags.Bool("ignore", false, "just print log and return ok even with failure")
	registerCommand(cmd)
}
