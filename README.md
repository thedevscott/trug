# Trug
A lightweight private local web app for keeping track of your monthly finances.

## Definition
* trug /trŭg/ - noun: A shallow, usually oval gardening basket used to hold harvests

# Motivation
Take responsiblity and protect your data from accidental exposure by default,
and turn spending into clear, actionable patterns so you can make
better financial choices faster.

# Quick Start
## Dependencies
* [Docker](https://www.docker.com/products/docker-desktop/)
* [Go](https://go.dev/)
* Git
* GNU Make

### Docker Setup
To run the app via docker do the following:

* clone/download the repo
```bash
git clone https://github.com/thedevscott/trug.git
```
* install & run [Docker](https://www.docker.com/products/docker-desktop/)
* enter the cloned repo & run the command
```bash
make compose-up
```
* open your browser and connect at [localhost:4000](https://localhost:4000)
* create an account, login and enjoy

**Note**: 
      
* the ***makefile*** contains several other commands you may find useful if you
  plan to modify the app.

# Usage
Once the app is up and running, open your browser and connect at
[localhost:4000](https://localhost:4000).
Once connected:
* click signup to make your account    
      * Be sure to write down your login info
* click login and enter your credentials

You are now able to enter transaction information manually and see your monthly
stats:
 * expenses, income and money left over


# Command Line Flags
run 'make help' to learn about all available flags:

-addr string

      HTTP network address (default ":4000")
-debug bool

      Enables detailed view of errors and stack traces in the browser
-dsn string

      MySQL data source name (default "web:trugpass@/trug?parseTime=true")

# Contributing
I am still thinking about what features to add, things like:

* Receipt OCR
* LLM/AI summary (local)
* Querable stats page

If you would like to work on any of those feel free to make a PR.

Thanks in advance.
