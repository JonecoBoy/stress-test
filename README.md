# Stress Test

This is a command-line tool for performing stress tests on a specified URL.

## Usage

The main command for the tool is `test`. Here's how you can use it:

```bash
stresser test --url http://google.com --requests 100 --concurrency 10
```
or
```bash
docker build -t stress-test .
docker run -it --rm stress-test test --url=http://google.com --requests=1000 --concurrency=10 --quiet -e txt -f stdout
```


This command will start a stress test on the URL http://localhost:8080/api/v1/test with 100 total requests and 10 concurrent requests.  
Flags
The test command supports the following flags:  
-u, --url: Set the target URL (Required)
-r, --requests: Set the number of total requests (Default: 100)
-c, --concurrency: Set the number of concurrent requests (Default: 10)
-o, --output: Set the output file for the report
-t, --timestamp: Add a timestamp to the output file name (Default: true)
-e, --encode: Set the output format for the report (Options: txt, json, toml, yaml,csv)
-f, --format: Set the report output format (Options: stdout, txt, html)
-q, --quiet: Set quiet mode, which won't print all response statuses (Default: false)


# Installation
To install the tool, you need to have Go installed on your machine. Then, you can clone this repository and build the tool using the Go compiler.  
License

# Releases
You can find the latest releases of the tool in the Releases section of this repository.

# Example
[View HTML Report](./report.html)
[View TXT Report](./report.txt)