package httpclient

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

//RandomOrgDiceRoller holds *http.Client and URI to make requests to constructed in method New
type RandomOrgDiceRoller struct {
	client  *http.Client
	genURI  string
	postURI string
}

//New sets up *http.Client with timeout and constructs URI from https://www.random.org/dice/?num= + nOfRolls if nOfRolls == "" then sets defaul value of 10
func New(timeout time.Duration, nOfRolls int, postURL string) *RandomOrgDiceRoller {
	s := strconv.Itoa(nOfRolls)
	if s == "" {
		log.Println("number of rolls doesn't contain any value, setting defaul value of 10")
		s = "10"
	}
	completeURI := "https://www.random.org/dice/?num=" + s
	return &RandomOrgDiceRoller{
		client: &http.Client{
			Timeout: timeout,
		},
		genURI:  completeURI,
		postURI: postURL,
	}
}

func getRequest(method, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return request, err
}

// GetResponse gets response
func (c *RandomOrgDiceRoller) GetResponse(method string, body io.Reader) (*http.Response, error) {
	request, err := getRequest(method, c.genURI, body)

	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		s := "Failed to get reponse 200 OK"
		return nil, errors.New(s)
	}
	return response, nil
}

//PostRequest send a POST request to specified url in RandomOrgDiceRoller and sets body
func (c *RandomOrgDiceRoller) PostRequest(body io.Reader) error {
	request, err := getRequest("POST", c.postURI, body)
	if err != nil {
		return err
	}
	response, err := c.client.Do(request)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	return err
}
