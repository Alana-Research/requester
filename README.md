# requester

A simple tool (<200 lines of code) to test concurrently HTTP requests status codes following its redirects.

- Tests your test suite endpoints are returning the desired status code (200, 202, 401...).
- If there should be some http redirects before hitting the endpoint you can test the request if doing them all with the desired status codes.
- No fancy and overcomplicated features.
- HEAD requests for faster response.
- Every test runs concurrently.
- You can create test suites using a YAML file.
- Add http headers to requests.

## Installation

```sh
brew...
```

## Usage

1. Create a config yml file as the example one config_example.yml and add your tests there following that schema.

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


