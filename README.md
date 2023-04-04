![main badge](https://github.com/LeoCBS/httpmiddleware/actions/workflows/makefile.yml/badge.svg?branch=main)

## Why create this project?

It is very common create microservices that use HTTP protocol on communication layer, this way we have a lot of projects which one with a different HTTP server implementantion and a lot of duplication code between projects. 


## What is this project?

This project is a middleware writting in Golang to simplify and avoid code duplication on handling HTTP requests between microservices.

What this project try abstract:

 * HTTP routes declaration
 * URL parameter validations
 * Write response headers
 * Error handling
 * Write responses
 * Gracefull shutdown (TODO)

### HTTP routes declaration

If you already try create a HTTP server and your routes using builtin `http` package, you must realized how costly it is to address requests,
parse URLs values and make basic validation like check HTTP methods. Having it in mind, we choose to use
the lib [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) as default HTTP router. Below is an example ilustrating how is
simple to define one new route:   

* Declaring a POST

```golang
	//register a simple route POST using key/value URL pattern
	md.POST("/name/:name/age/:age", fnHandlePOST)
```


* Declaring a POST and your handler

```golang
	fnHandlePOST := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
		// here you will add your business logic, call some storage func, etc...
		return httpmiddleware.Response{
			StatusCode: http.StatusOK,
		}
	}
	//register a simple route POST using key/value URL pattern
	md.POST("/name/:name/age/:age", fnHandlePOST)
```

### URL parameter validation

Another validation that is repeatedly (and duplicated) made when we are using key/value URL
pattern is checked if one value key isn't empty. Here we have one example to how
send one empty value:

Here we have one URL with two values/parameters:
    
    md.POST("/name/:name/age/:age", fnHandlePOST)

Supose that your service send one request with this URL:

    /name//age/37
     

Using this lib, you don't need check if your parameter is empty `(GetParamValue("name") != "")`,
your server will automatic reply a Bad Request error informaing which parameter
is wrong, like this:

    401 Bad Request {"error":"your URL must inform name value"}

Take a lot on unit test to check one full example :) 

### Write response headers

To write custom response headers just use atributte `Response.Headers` that
middleware will write values on reponse.

```
	fnHandlePOST := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
	        headerKey := "Location"
        	headerValue := "/whatever/01234"
	        respHeaders := map[string]string{
		    headerKey: headerValue,
        	}
		return httpmiddleware.Response{
			StatusCode: http.StatusOK,
			Headers:    respHeaders,
		}
	}
```

### Error Handling

Another common behavior that you always must to do on a HTTP microservices is
handling error. The lack of a standard for handling errors leads to confusion
and code duplication. Handling error on middleware layer we simplify source code and
focused only on the business logic.

When you use `github.com/LeoCBS/httpmiddleware/errors` inside your core business,
middleware will take care of handling the error properly, writing the correct HTTP Status Code
and right response body for the client.


Write a Bad Request error to the client:

```
fnHandle := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
    return httpmiddleware.Response{
        ror: errors.NewBadRequest("your body must be a valid JSON"),
    }
}
URL := "/clienterrorhandling"
f.md.GET(URL, fnHandle)
```

Returning a `errors.NewBadRequest` middleware will handling error and will
write HTTP Status Code 400 and write on response body your custom message to
the client `{"error":"your body must be a valid JSON"}`

Access middleware_test.go to check more examples to how use all custom errors.

### Write JSON responses

In most cases we want write a JSON in our response body to the client, so we
always write and duplicate code between microservices that encode a struct and check
if some errors happens in this process. This way this middleware already convert your
struct into JSON and write it on response body and you just need return
a `httpmiddleware.Body` like code below.

```
	handleFn := func(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
	    type person struct {
		Name string `json:"name"`
            }
	    record := person{Name: "leo"}
	    return httpmiddleware.Response{
	        StatusCode: http.StatusOK,
		Body:       expectedBody,
	    }
	}
```

## How use this middleware?

Here we have a simple application that start a HTTP Server on port 8080 and
handling a GET and POST on path `/name/:name`. Pay attention how few lines of code
are needed to write the handlers.

```
package main

import (
	"net/http"

	"github.com/LeoCBS/httpmiddleware"
	"github.com/sirupsen/logrus"
)

type user struct {
	Login string `json:"login"`
}

func main() {
	l := logrus.New()
	md := httpmiddleware.New(l)
	md.POST("/name/:name", createUser)
	md.GET("/name/:name", getUser)

	s := &http.Server{
		Addr:    ":8080",
		Handler: md,
	}
	panic(s.ListenAndServe())
}

func createUser(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
	return httpmiddleware.Response{
		StatusCode: http.StatusCreated,
		Body:       user{Login: ps.ByName("name")},
	}
}

func getUser(w http.ResponseWriter, r *http.Request, ps httpmiddleware.Params) httpmiddleware.Response {
	return httpmiddleware.Response{
		StatusCode: http.StatusCreated,
		Body:       user{Login: ps.ByName("name")},
	}
}
```

Running a CURL to check:

    curl -X POST "http://localhost:8080/name/leonardo"
    {"login":"leonardo"}


## More Examples

Access unit tests to see how to use this middleware.

## Testing

To run unit tests:

    make check

## Advanges to don't use default `http.Handler`

To simplify HTTP router definition and get more features, we decide to use the lib
[julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) instead
default http.Handler, we recomend access this project on github to understand
which more feature are avaliable to use like: router multi-domains and
sub-domains, basic auth and others features. 
