# requester

<p align="center" >
  <img alt="Go report card" src="https://goreportcard.com/badge/github.com/Alana-Research/requester">
  <img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/Alana-Research/requester">
  <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/Alana-Research/requester">
  <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/Alana-Research/requester">
  <img alt="Lines of code" src="https://img.shields.io/tokei/lines/github/Alana-Research/requester">
</p>


A simple tool (<200 lines of code) to test concurrently HTTP requests status codes following its redirects.

🚀 Tests your test suite endpoints are returning the desired status code (200, 202, 401...).

🚀 If there should be some http redirects before hitting the endpoint you can test the request if doing them all with the desired status codes.

🚀 No fancy and overcomplicated features.

🚀 HEAD requests for faster response.

🚀 Every test runs concurrently.

🚀 You can create test suites using a YAML file.

🚀 Add http headers to requests.

## How to install

```sh
brew tap Alana-Research/requester https://github.com/Alana-Research/requester
brew install requester
```

## How to use

1. Create a config yml file as the example one config_example.yml and add your tests there following that schema.

Example:

```yml
test-name-1:
  httpUrl: http://example.com/notfound
  expectedStatusCodes: ["301", "404"]

test-name-2:
  httpUrl: https://example.com/superauth?md5=123&&nothing=none
  expectedStatusCodes: ["301", "403"]
  ignoreTLSError: True 

test-name-3:
  httpUrl: https://example.com/awesomefile
  httpHeaders: 
    - "X-file": "awesome.png"
  expectedStatusCodes: ["202"]

  ...
```


2. Run:

```
requester --config=/path/configXXX.yml 
```


## Roadmap

- [ ] Support all other HTTP methods (currently only GET/HEAD)
- [ ] Add more validations features as body response, headers, hostnames redirects...
- [ ] DNS resolver per test
- [ ] Proxy support

---

### License

Congo is released under the MIT license. See [`LICENSE`](https://github.com/Alana-Research/requester/blob/master/LICENSE) for more details.
