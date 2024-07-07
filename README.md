# Stress Test

This is a command-line tool for performing stress tests on a specified URL.

## Usage

The main command for the tool is `test`. Here's how you can use it:

First you need to clone the project and [Build](##build) the binary:

```bash
./stresser test -u https://google.com -q -r 100 -c 10
```
or you can clone this repository and run directly from the docker image
```bash
docker run -it --rm joneco/stress-test:latest test -u https://google.com -q -r 100 -c 10
```

you can also build the docker image instead of downloading from docker huband run the command from the image
```bash
docker build -t stress-test .
docker run -it stress-test test -u https://google.com -q -r 100 -c 10
````

### Flags

This command will start a stress test on the URL http://localhost:8080/api/v1/test with 100 total requests and 10 concurrent requests.  
Flags
The test command supports the following flags:  
-u, --url: Set the target URL (Required)
-r, --requests: Set the number of total requests (Default: 100)
-c, --concurrency: Set the number of concurrent requests (Default: 10)
-o, --output: Set the output file for the report
-t, --timestamp: Add a timestamp to the output file name (Default: true)
-e, --encode: Set the output format for the report (Options: txt, json, toml, yaml,csv)
-q, --quiet: Set quiet mode, which won't print all response statuses (Default: false)


## Build
To install the tool, you need to have Go installed on your machine. Then, you can clone this repository and build the tool using the Go compiler.  
```bash
docker build -t stress-test .
```

## Releases
You can find the latest releases of the tool in the Releases section of this repository.

## Example
[View HTML Report](https://rawcdn.githack.com/JonecoBoy/stress-test/c8f4db1b6f685c1dd02e2028e4d546ce38bdbc9d/report.html)

[View TXT Report](./report.txt)

## FAQ
- If you are using with docker run, it won't have direct access to the file system, so it wont save the html and file reports. So you need to build a volume, so run:
    ```bash
    docker run -it --rm -v $(pwd):/home joneco/stress-test:latest test -u https://google.com -q -r 100 -c 10 -e json -o /home/report.json 
  ```
