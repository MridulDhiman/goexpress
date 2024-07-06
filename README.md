# Express.js implementation in golang 

## Features: 
1. Custom router implementation.
2. Route Groups 
3. Local and global middleware support

## Key Takeways: 
1. Whenever we make a HTTP request, browser also requests for `favicon.ico`.
2. For custom router implementation, we also have to reimplement all the functions in the `http.Handler` interface.
3. You can add multiple (key,value) pairs to request context, without affecting it's original structure.
4. You can use variadic parameters (e.g. `func (x ...int)`) for providing any no. of parameters to go functions 
5. You cannot have struct(or pointer to struct) as map key.

## Issues(for now): 
1. Does not check for duplicate route entries across multiple route groups