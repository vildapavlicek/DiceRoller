package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"

	httpclient "github.com/vildapavlicek/GoLang/DiceGetter/httpClient"

	"golang.org/x/net/html"
)

//DiceRolls main struct
type DiceRolls struct {
	NumOfRolls int                             `json:"-"`
	Client     *httpclient.RandomOrgDiceRoller `json:"-"`
	Data       []int                           `json:"dice"`
}

//New returns new *DiceRolls
func New(hc *httpclient.RandomOrgDiceRoller) *DiceRolls {
	return &DiceRolls{
		Client: hc,
		Data:   []int{},
	}
}

//ParseHTML parser for results from page https://www.random.org Closes response.Body
func (d *DiceRolls) ParseHTML(response *http.Response) error {
	defer response.Body.Close()
	doc, err := html.Parse(response.Body)
	if err != nil {
		return err
	}

	d.parseAltAtribute(doc)

	return nil
}

func (d *DiceRolls) parseAltAtribute(n *html.Node) {
	if n.Data == "img" && n.Type == html.ElementNode {
		d.parseRollNumber(n)
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		d.parseAltAtribute(child)
	}
}

func (d *DiceRolls) parseRollNumber(n *html.Node) {
	var i int
	var err error
	for _, att := range n.Attr {
		if att.Key == "alt" {
			if i, err = strconv.Atoi(att.Val); err != nil {
				log.Fatalf("Failed to convert string value to int, value: %s", att.Val)
			}
			d.Data = append(d.Data, i)
		}
	}
}

//BucketResults buckets results to map as k: rolled number v: how manytimes
func (d *DiceRolls) BucketResults(output io.Writer) {
	var prev int
	var count = 1

	for i, n := range d.Data {

		switch {
		case prev == 0:
		case prev != n:
			fmt.Fprintf(output, "%v -> %v\n", prev, count)
			count = 1
		case n == prev:
			count++
		}

		if i == len(d.Data)-1 {
			fmt.Fprintf(output, "%v -> %v\n", n, count)
		}
		prev = n
	}
}

//OrderResults orders parsed rolls as ascending
func (d *DiceRolls) OrderResults(output io.Writer) {
	sort.Ints(d.Data)
	fmt.Fprintf(output, "%v\n", d.Data)
}

//DoPost posts data to specified URL
func (d *DiceRolls) DoPost() error {
	b := make([]byte, 0, 1)
	buffer := bytes.NewBuffer(b)
	json.NewEncoder(buffer).Encode(d)
	d.Client.PostRequest(buffer)
	return nil
}
