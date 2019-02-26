# CS3031 Web Proxy Assignment
James Tait - 16321184

I used the Go programming language to build this web proxy, using only the built-in packages. The main function simply reads the console for input to blacklist URLs in one 'goroutine' (thread), and creates a server to listen for requests on port 8080. When a request is received, a new thread is created and subsequently runs the HTTP handler.  

For HTTP requests, the request is first checked against the list of blacklisted URLs. If the requested URL is found in the blacklist, then a `HTTP/1.0 403 Forbidden` error response is returned to the client. The method of the request is then checked and if it is a `CONNECT` request, the HTTPS handler is called. Otherwise, the request is forwarded on to the server and then the response is just relayed to the client.

For HTTPS requests, 
