package main

import (
	"crypto/tls"
	"flag"
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
var versionTag string

const MAX_REDIRECTS int8 = 10
const REDIRECT_CODE_START int16 = 300
const REDIRECT_CODE_END int16 = 400

type testEndpoint struct {
	HttpUrl             string              `yaml:"httpUrl"`
	HttpHeaders         []map[string]string `yaml:"httpHeaders"`
	ExpectedStatusCodes []string            `yaml:"expectedStatusCodes"`
	IgnoreTLSError      bool                `yaml:"ignoreTLSError"`
}

func (obj *testEndpoint) requestHTTP() ([]string, error) {
	url := obj.HttpUrl
	respCodes := []string{}
	var counter int8 = 0

	for counter <= MAX_REDIRECTS {
		req, err := http.NewRequest("GET", url, nil)

		transport := &http.Transport{}
		if obj.IgnoreTLSError {
			transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: transport,
		}

		for _, headerMap := range obj.HttpHeaders {
			for header, value := range headerMap {
				req.Header.Set(header, value)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		respCodes = append(respCodes, strconv.Itoa(resp.StatusCode))

		if int16(resp.StatusCode) < REDIRECT_CODE_START || int16(resp.StatusCode) >= REDIRECT_CODE_END { //not a redirect http code
			return respCodes, nil
		} else {
			url = resp.Header.Get("Location")
			counter += 1
		}
	}
	return respCodes, nil
}

func verifyStatus(statusCodesReceived []string, testName string, expectedStatusCodes []string) bool {
	if len(statusCodesReceived) != len(expectedStatusCodes) {
		printFail(fmt.Sprintf("Test %s FAILED. Expected %s Received %s.", testName, expectedStatusCodes, statusCodesReceived))
		return false
	}

	for i, code := range statusCodesReceived {
		if code != expectedStatusCodes[i] {
			printFail(fmt.Sprintf("Test %s FAILED. Expected %s Received %s.", testName, expectedStatusCodes, statusCodesReceived))
			return false
		}
	}
	printSuccess(fmt.Sprintf("Test %s SUCCESS. Expected %s Received %s.", testName, expectedStatusCodes, statusCodesReceived))
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

func getTestsFromYaml(testsData *map[string]testEndpoint, configFilePath *string) {
	testsFile, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		tlog.BigError("Config file could not be opened:", err.Error())
		os.Exit(1)
	}

	errUnmarshal := yaml.Unmarshal(testsFile, &testsData)
	if errUnmarshal != nil {
		tlog.BigError("YAML file has incorrect fields or syntax:", errUnmarshal)
		os.Exit(1)
	}
}

func main() {
	configFilePath := flag.String("config", "", "Config file path: --config=./path/configXX.yml")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "requester %s\nUsage:\n", versionTag)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *configFilePath == "" {
		tlog.BigError("Config file not found")
		tlog.Info("Run -help")
		os.Exit(1)
	}

	testsData := make(map[string]testEndpoint)
	getTestsFromYaml(&testsData, configFilePath)

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
