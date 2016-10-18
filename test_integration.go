package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	selenium "sourcegraph.com/sourcegraph/go-selenium"
)

func integrationTests(l Logger) error {
	l = newLogger(timed)
	appURL := "http://127.0.0.1:3000"
	seleniumURL := "http://localhost:4444/wd/hub"
	for _, u := range []string{appURL, seleniumURL} {
		check := func(s string) func() (bool, error) {
			return func() (bool, error) {
				rsp, err := http.Get(s)
				if err != nil {
					return false, nil
				}
				defer rsp.Body.Close()
				if rsp.Status[0] != '2' {
					b, _ := ioutil.ReadAll(rsp.Body)
					return false, nil
				}
				return true, nil
			}
		}(u)
		err := waitFor(1*time.Second, 5*time.Minute, check)
		if err != nil {
			return fmt.Errorf("error pinging %s: %s", u, err)
		}
		l.Printf("%s available", u)
	}
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	selenium.Log = log.New(ioutil.Discard, "", 0)
	d, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		return err
	}
	defer d.Close()

	err = d.Get(appURL)
	if err != nil {
		return err
	}
	l.Printf("got response")

	element, err := d.Q("a")
	if err != nil {
		return err
	} else if element == nil {
		return fmt.Errorf("unable to find link")
	}
	l.Printf("found link")
	err = element.Click()
	if err != nil {
		return err
	}
	element, err = d.Q("div#clicked")
	if err != nil {
		return err
	}
	txt, err := element.Text()
	if err != nil {
		return err
	}
	if v, x := txt, "you clicked!"; x != v {
		return fmt.Errorf(`expected txt to be %#v, was %#v`, x, v)
	}
	l.Printf("found element with txt=%q", txt)
	return nil
}

func waitFor(waitDur, timeoutDur time.Duration, f func() (bool, error)) error {
	tick := time.Tick(waitDur)
	timeout := time.After(timeoutDur)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout")
		case <-tick:
			ok, err := f()
			if err != nil {
				return err
			} else if ok {
				return nil
			}
		}
	}
}
