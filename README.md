# Express.js implementation in golang 

## Features: 
1. Custom router implementation.
2. Route Groups 

## Key Takeways: 
1. Whenever we make a HTTP request, browser also requests for "favicon.ico".
2. For custom router implementation, we also have to reimplement all the functions in the `http.Handler` interface.
3. You can add multiple (key,value) pairs to request context, without affecting it's original structure.

## Issues(for now): 
1. Does not check for duplicate route entries across multiple route groups