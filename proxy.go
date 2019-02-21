package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
)

var hasPort = regexp.MustCompile(`:\d+$`)

func removeProxyHeaders(r *http.Request) {
	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")
	r.Header.Del("Proxy-Authenticate")
	r.Header.Del("Proxy-Authorization")
	// Connection, Authenticate and Authorization are single hop Header:
	// http://www.w3.org/Protocols/rfc2616/rfc2616.txt
	// 14.10 Connection
	//   The Connection general-header field allows the sender to specify
	//   options that are desired for that particular connection and MUST NOT
	//   be communicated by proxies over further connections.
	r.Header.Del("Connection")
}

func copyAndClose(dst, src *net.TCPConn, host string) {
	if _, err := io.Copy(dst, src); err != nil {
		fmt.Printf("Error copying to client: %s", err)
	}
	dst.CloseWrite()
	src.CloseRead()
	if host != "" {
		fmt.Printf("Connection to %s Closed\n", host)
	}
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	client := &http.Client{}
	fmt.Printf("Rquest received: %s %s\n", req.Method, req.URL)
	removeProxyHeaders(req)

	if req.Method == "CONNECT" {
		handleHTTPS(w, req)
	} else {
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Response received: %s\n\n", res.Status)
		res.Write(w)

		req.Body.Close()
		res.Body.Close()
	}
}

func handleHTTPS(w http.ResponseWriter, req *http.Request) {

	req.URL.Scheme = "https"

	hij, ok := w.(http.Hijacker)
	if !ok {
		panic("httpserver does not support hijacking")
	}

	proxyClient, _, err := hij.Hijack()
	if err != nil {
		panic("Cannot hijack connection " + err.Error())
	}

	host := req.URL.Host
	if !hasPort.MatchString(host) {
		host += ":80"
	}

	targetSiteCon, err := net.Dial("tcp", host)
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Printf("Accepting CONNECT to %s\n", host)

	proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	targetTCP, targetOK := targetSiteCon.(*net.TCPConn)
	proxyClientTCP, clientOK := proxyClient.(*net.TCPConn)
	if targetOK && clientOK {
		go copyAndClose(targetTCP, proxyClientTCP, "")
		go copyAndClose(proxyClientTCP, targetTCP, host)
	}
}

func main() {
	httpHandler := http.HandlerFunc(handleHTTP)
	fmt.Printf("Proxy activated!\n\n")
	http.ListenAndServe(":8080", httpHandler)

}
