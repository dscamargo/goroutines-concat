package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type GithubResponse struct {
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Error     error  `json:"error"`
}

type OrgsResponse struct {
	Login string `json:"login"`
}

func request(url string, response interface{}) (*http.Response, error) {
	start := time.Now()
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Error] Request - URL: %v - Error: %v", url, err.Error()))
	}
	//req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN_KEY"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Error] Request - URL: %v - Error: %v", url, err.Error()))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Error] ioutil.ReadAll - URL: %v - Error: %v", url, err.Error()))
	}
	if resp.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("[Error] Request - URL: %v - Status Code: %v - Body: %v", url, resp.StatusCode, string(body)))
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Error] json.Unmarshal - URL: %v - Error: %v", url, err.Error()))
	}
	final := time.Since(start).Milliseconds()
	fmt.Printf("Request: URL - %v - Time: %v ms\n", url, final)
	return resp, err
}

func getUsers(org string) []string {
	var response []OrgsResponse
	url := fmt.Sprintf("https://api.github.com/orgs/%s/members", org)
	resp, err := request(url, &response)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var users []string
	for _, item := range response {
		users = append(users, item.Login)
	}
	return users
}

func main() {
	start := time.Now()
	gitUsers := getUsers("microsoft")
	var wg sync.WaitGroup
	ch := make(chan GithubResponse)

	log.Printf("Searching %v users...", len(gitUsers))
	for _, user := range gitUsers {
		wg.Add(1)
		go worker(ch, &wg, user)
	}

	myArray := make([]GithubResponse, len(gitUsers))
	for i := range myArray {
		myArray[i] = <-ch
	}

	wg.Wait()
	close(ch)
	final := time.Since(start).Milliseconds()
	log.Printf("Result: %v", myArray)
	log.Printf("Result Length: %v", len(myArray))
	log.Printf("Time: %v ms\n", final)
}

func worker(ch chan GithubResponse, wg *sync.WaitGroup, username string) {
	defer wg.Done()
	response := GithubResponse{}
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	resp, err := request(url, &response)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ch <- response
}
