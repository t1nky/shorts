# Shorts

Shorts is a REST API written in [Go](https://golang.org/ "Go") for creating short links

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

 [Make](https://en.wikipedia.org/wiki/Make_(software) "Make"), [Go](https://golang.org/ "Go") and [Docker](https://www.docker.com/ "Docker")


### Installing

Run commands in the root directory of the project

Initialize database:

`make build_db`

This will create a database container named `shorts` with 2 databases (1 for testning and 1 for deployment)

To build docs, run

`make docs`

To run the project:

`go run main.go` or `make run`

## Running the tests

Run `make test` or `go test` in the root directory of the project

## License

MIT License

Copyright (c) 2020 Aleksandr Bondarenko

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
