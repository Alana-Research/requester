package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	tlog "github.com/Alana-Research/terminal-app-log"
	"gopkg.in/yaml.v3"
)

var mutex sync.Mutex

type testEndpoint struct {
	HttpUrl            string `yaml:"httpUrl"`
	HttpBody           string `yaml:"httpBody"`
	ExpectedStatusCode string `yaml:"expectedStatusCode"`
}

func (obj *testEndpoint) requestHTTP() (string, error) {

	//if body do the request with the body

	return "200", nil
}

func verifyStatus(statusCodeReceived string, testName string, expectedStatusCode string) bool {
	if statusCodeReceived != expectedStatusCode {
		printFail(fmt.Sprintf("Test %s FAILED. Expected %s Received %s.",
			testName,
			expectedStatusCode,
			statusCodeReceived,
		))
		return false
	} else {
		printSuccess(fmt.Sprintf("Test %s SUCCESS. Expected %s Received %s.",
			testName,
			expectedStatusCode,
			statusCodeReceived,
		))
		return true
	}
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

	failedTests := 0
	for testName, test := range testsData {
		go func(test testEndpoint, testName string) {
			statusCodeReceived, err := test.requestHTTP()
			if err != nil {
				tlog.BigError("Cannot perform the request:", err.Error())
				wg.Done()
				return
			}
			if !verifyStatus(statusCodeReceived, testName, test.ExpectedStatusCode) {
				failedTests += 1
				wg.Done()
				return
			}
			wg.Done()
		}(test, testName)
	}
	wg.Wait()
	fmt.Println()
	tlog.Info(fmt.Sprintf("TOTAL tests failed: %d/%d", failedTests, len(testsData)))
}
