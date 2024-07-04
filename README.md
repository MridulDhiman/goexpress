# Express.js implementation in golang (ongoing)

## Features: 
1. Custom router implementation (for golang).
2. Route Groups 

## Key Takeways: 
1. Whenever we make a HTTP request, browser also requests for "favicon.ico".
2. For custom router implementation, we also have to reimplement all the functions in the `http.Handler` interface.

## Issues(for now): 
1. Does not check for duplicate route entries across multiple route groups