package main

import (
	"fmt"
	"io/ioutil"
	"log"

	selenium "sourcegraph.com/sourcegraph/go-selenium"
)

func integrationTests(l Logger) error {
	l = newLogger(timed)
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	selenium.Log = log.New(ioutil.Discard, "", 0)
	d, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		return err
	}
	defer d.Close()

	err = d.Get("http://127.0.0.1:3000")
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
