# requester

A simple (<500 lines of code) and concurrent tool to test HTTP requests status codes following its redirects.

- Create test suites using a YAML file.
- Every test runs concurrently.
- Test that all redirects a request does are the expected.
- No fancy and useless features, it only does one job.
- Option to add headers to the requests.

## Installation

```sh
brew...
```

## Usage

1. Create a config.yml file as the example one config_example.yml and add your tests there following that schema.

2. Run xxxxxx 

## Roadmap

- Support all other HTTP methods (currently only GET)
- Add more validations features as body response, headers, domains redirects...
