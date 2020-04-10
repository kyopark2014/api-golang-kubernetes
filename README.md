# restapi-golang-sample

REST API using golang

This source is an example to explain how to deploy RESTful API using golang.
I think it is useful for who wants to learn go language and REST API. 

I used the basic source from Brad's youtube but it is modified to use in a real environment. 

Reference:
Source: “Golang REST API With Mux”, Brad Traversy
  https://www.youtube.com/watch?v=SonwZ6MF5BE


## RUN

You may run it in command shall but you can do it in docker environment easily.

```c
go run main.go
```


## Initiate and StartService
The main is initiating and starting the service.

```go
err := Initialize()
if err != nil {
	log.E("Failed to initialize service: %v", err)
	os.Exit(1)
}

err = StartService()
if err != nil {
	log.E("Failed to start service: %v", err)
	os.Exit(1)
}
log.E("Exiting service ...")
```

The main function for the restful api is starting from server.InitServer() in "main.go".
  
```go
func StartService() error {
	log.D("start the service...")

	var err error
	if err = server.InitServer(); err != nil {
	log.E("Failed to start the HTTP(s) server: err:[%v]", err)
}
```  

## Config
#### config.go

Define AppConfig in order to load the configuration

```go
var config *AppConfig
  
func GetInstance() *AppConfig {
	if config == nil {
		config = &AppConfig{}
	}
	return config
}
  
Logging struct {
	Enable bool   `json:"Enable"`
	Level  string `json:"Level"`
} 
```

#### main.go

The configureation is loaded from "config.json".

```go
var configFileName string = "configs/config.json"

conf := config.GetInstance()
if !conf.Load(configFileName) {
	log.E("Failed to load config file: %s", configFileName)
	os.Exit(1)
}
```

## Logging

#### main.go

The log level is setting using SetupLogger.

```go
import ("restapi-golang-sample/pkg/log")

log.SetupLogger(conf.Logging.Enable, conf.Logging.Level)
```

Then, it is used as bellow.

```go
log.I("Starting service ...")
log.E("Failed to load config file: %s", configFileName)
```

#### log.go

The log level is appliying as bellow.

```go
func SetupLogger(isEnabled bool, level string) {
	loggingEnable = isEnabled
	backend := logging.NewLogBackend(os.Stdout, "", 0)

	backendFormatter := logging.NewBackendFormatter(backend, format)

	var lvl logging.Level
	switch level {
	case "ERROR":
		lvl = logging.ERROR
	case "WARNING":
		lvl = logging.WARNING
	case "INFO":
		lvl = logging.INFO
	case "DEBUG":
		lvl = logging.DEBUG
	default:
		lvl = logging.INFO
	}

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(lvl, "")

	logging.SetBackend(backendLeveled)
}

// D writes debug level log
func D(format string, v ...interface{}) {
	if loggingEnable {
		log.Debugf(format, v...)
	}
}

// E writes error level log
func E(format string, v ...interface{}) {
	if loggingEnable {
		log.Errorf(format, v...)
	}
}
```

## server.go

In this source, the information of books is managed by an array since it is an simple example to summarize the operation of RESTful API.

```go
var books []Book

books = append(books, Book{ID: "1", Isbn: "448743", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Park"}})
books = append(books, Book{ID: "2", Isbn: "312234", Title: "Book Two", Author: &Author{Firstname: "Steave", Lastname: "Smith"}})

```

The http server is based on mux as bellow.
```go
func InitServer() error {
	// Init Router
	r := mux.NewRouter()

	// Route Handler / Endpoints
	r.HandleFunc("/api/books", CreatBook).Methods("POST")
	r.HandleFunc("/api/books", GetBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", GetBook).Methods("GET")
	r.HandleFunc("/api/books/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", DeleteBook).Methods("DELETE")

	// log.Fatal(http.ListenAndServe(":8000", r))
	var err error
	err = http.ListenAndServe(":8000", r)

	return err
}
```

In order to create an item of book, json format was used as bellow.

```go
func CreatBook(res http.ResponseWriter, req *http.Request) {
	log.D("Write the book information")

	res.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(req.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID
	books = append(books, book)
	json.NewEncoder(res).Encode(book)
}
```

Json is printed for debugging
```go
jsondata, err := json.Marshal(book)
	if err != nil {
		log.E("Cannot encode to Json", err)
	}
	log.D("%v", string(jsondata))
}
```

Let me show an example of an book information.

```text
CREAT (example)
POST api/books
{
    "isbn": "123812",
    "title": "Book Three",
    "author": {
            "firstname": "Jain",
            "lastname": "Lee"
    }
}

UPDATE (example)
PUT api/books/{id}
{
    "isbn": "123812",
    "title": "Book Three - Update",
    "author": {
            "firstname": "Jain",
            "lastname": "Choi"
    }
}

```
