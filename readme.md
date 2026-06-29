# Trug
A lightweight private local web app for keeping track of your finances.

### Definition
* trug /trŭg/ - noun: A shallow, usually oval gardening basket used to hold harvests

# Dependencies
    * Git
    * Go
    * Docker
    * GNU Make

# Docker Setup
    * clone/download the repo
    ```bash
    git clone https://github.com/thedevscott/trug.git
    ```
    * install & run [Docker](https://www.docker.com/products/docker-desktop/)
    * enter the cloned repo & run the command
    ```bash
    make compose-up
    ```
    * open your browser and connect at localhost:4000
    [open trug](localhost:4000)
    * create an account, login and enjoy

# Flags
run 'make help' to learn about all available flags
  -addr string
        HTTP network address (default ":4000")
  -debug
        Enables detailed view of errors and stack traces in the browser
  -dsn string
        MySQL data source name (default "web:trugpass@/trug?parseTime=true")
