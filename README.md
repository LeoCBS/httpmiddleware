![main badge](https://github.com/LeoCBS/httpmiddleware/actions/workflows/makefile.yml/badge.svg?branch=main)

## Why create this project?

It is very common create microservices that use HTTP protocol on communication layer, this way we have a lot of projects which one with a different HTTP server implementantion and a lot of duplication code between projects. 


## What is this project?

This project is a middleware writting in Golang to simplify and avoid code duplication on handling HTTP requests between microservices.

What this project try abstract:

 * HTTP routes declaration
 * URL parameter validations
 * Error handler
 * Write responses
 * Gracefull shutdown

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

## How use this middleware?

TODO put here go get and how create one middleware
 

## More Examples

Access unit tests to see how to use this middleware.

## Testing

To run unit tests:

    make check
