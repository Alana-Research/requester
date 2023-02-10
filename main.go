package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"

	tlog "github.com/Alana-Research/terminal-app-log"
	"gopkg.in/yaml.v3"
)

var mutex sync.Mutex

const MAX_REDIRECTS int8 = 10

type testEndpoint struct {
	HttpUrl             string   `yaml:"httpUrl"`
	HttpHeader          []string `yaml:"httpHeaders"`
	ExpectedStatusCodes []string `yaml:"expectedStatusCodes"`
}

func (obj *testEndpoint) requestHTTP() ([]string, error) {
	url := obj.HttpUrl
	respCodes := []string{}
	var counter int8 = 0

	for counter <= MAX_REDIRECTS {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}

		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}

		respCodes = append(respCodes, strconv.Itoa(resp.StatusCode))

		if resp.StatusCode < 300 || resp.StatusCode >= 400 { //not a redirect http code
			return respCodes, nil
		} else {
			url = resp.Header.Get("Location")
			counter += 1
		}
	}
	return respCodes, nil
}

func verifyStatus(statusCodesReceived []string, testName string, expectedStatusCodes []string) bool {
	fmt.Println("here")

	if len(statusCodesReceived) != len(expectedStatusCodes) {
		printFail(fmt.Sprintf("Test %s FAILED. Expected %s Received %s.",
			testName,
			expectedStatusCodes,
			statusCodesReceived,
		))
		return false
	}

	for i, code := range statusCodesReceived {
		if code != expectedStatusCodes[i] {
			printFail(fmt.Sprintf("Test %s FAILED. Expected %s Received %s.",
				testName,
				expectedStatusCodes,
				statusCodesReceived,
			))
			return false
		}
	}

	printSuccess(fmt.Sprintf("Test %s SUCCESS. Expected %s Received %s.",
		testName,
		expectedStatusCodes,
		statusCodesReceived,
	))
	return true
}

func printSuccess(text string) {
	mutex.Lock()
	defer mutex.Unlock()
	tlog.Success(text)
}

func printFail(text string) {
	mutex.Lock()
	defer mutex.Unlock()
	tlog.Fail(text)
}

func getTestsFromYaml(testsData *map[string]testEndpoint) {
	testsFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		tlog.BigError("File config.yml could not be opened:", err.Error())
		os.Exit(1)
	}

	errUnmarshal := yaml.Unmarshal(testsFile, &testsData)
	if errUnmarshal != nil {
		tlog.BigError("YAML file has incorrect fields or syntax:", errUnmarshal)
		os.Exit(1)
	}
}

func main() {
	testsData := make(map[string]testEndpoint)

	getTestsFromYaml(&testsData)

	var wg sync.WaitGroup
	wg.Add(len(testsData))

	passedTests := 0
	for testName, test := range testsData {
		go func(test testEndpoint, testName string) {
			statusCodesReceived, err := test.requestHTTP()
			if err != nil {
				tlog.BigError("Cannot perform the request:", err.Error())
				wg.Done()
				return
			}
			if !verifyStatus(statusCodesReceived, testName, test.ExpectedStatusCodes) {
				wg.Done()
				return
			}
			passedTests += 1
			wg.Done()
		}(test, testName)
	}
	wg.Wait()
	fmt.Println()
	tlog.Info(fmt.Sprintf("TOTAL tests passed: %d/%d", passedTests, len(testsData)))
}
