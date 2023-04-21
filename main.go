package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func handleIndex(conn net.Conn, uri *url.URL, count int) error {
	conn.Write([]byte("<h1>Hello, World!</h1>"))
	conn.Write([]byte(fmt.Sprintf("<h2>Count: %d</h2>", count)))

	xStr := uri.Query().Get("x")
	yStr := uri.Query().Get("y")

	x, err := strconv.Atoi(xStr)
	if err != nil {
		return err
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return err
	}

	sum := x + y

	conn.Write([]byte(fmt.Sprintf("<h3>%d + %d = %d</h3>", x, y, sum)))

	q, err := json.MarshalIndent(uri.Query(), "", "  ")
	if err != nil {
		return err
	}

	conn.Write([]byte(fmt.Sprintf("<pre>%s</pre>", q)))

	return nil
}

func handleAbout(conn net.Conn, uri *url.URL) error {
	conn.Write([]byte("<div>"))
	conn.Write([]byte("<h1>About us</h1>"))
	conn.Write([]byte(`<img src="https://fileasy.ru/keyboard.jpg" alt="image should be here" width="600px">`))
	conn.Write([]byte("</div>"))

	return nil
}

func handleHello(conn net.Conn, uri *url.URL) error {
	conn.Write([]byte("<h1>Hello, World!</h1>"))

	return nil
}

func main() {
	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	count := 0

	for {
		conn, _ := ln.Accept()
		count += 1

		r := bufio.NewReader(conn)
		str, _ := r.ReadString('\n')

		s := strings.Split(str, " ")
		if len(s) < 2 {
			continue
		}

		uri, _ := url.Parse(s[1])

		conn.Write([]byte("HTTP/3 200 OK\r\n"))
		conn.Write([]byte("Content-Type: text/html\r\n"))

		conn.Write([]byte("\r\n"))

		switch uri.Path {
		case "/":
			handleIndex(conn, uri, count)
			break

		case "/about":
			handleAbout(conn, uri)
			break

		case "/hello":
			handleHello(conn, uri)
			break

		default:
			conn.Write([]byte("<h1>404 Not Found</h1>"))
			break
		}

		conn.Close()
	}
}
