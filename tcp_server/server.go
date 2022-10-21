package main

import (
	"fmt"
	"net"
	"time"
)

func log(conn net.Conn, message string) {
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "-", conn.RemoteAddr().String(), "-",  message)
}

func main() {
    addr := ":5004"
    timeoutSeconds := 10 * time.Second

    // Listen for incoming connections
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        panic(err)
    }

    // Close the listener when the application closes
    defer listener.Close()

    // Get the connection
    host, port, err := net.SplitHostPort(listener.Addr().String())
    if err != nil {
        panic(err)
    }
    fmt.Println("Listening on host: " + host +", port: " + port)

    for {
        // Listen for an incoming connection
        conn, err := listener.Accept()
        if err != nil {
            panic(err)
        }
        fmt.Println("")
        log(conn, "connection received")

        var timer *time.Timer
        var canceled = make(chan struct{})

        // handler
        go func () {
            log(conn, "TCP session open")

            // Close the connection when goroutine completes
            defer conn.Close()

            for {
                // Read the incoming connection into the buffer
                buffer := make([]byte, 1024)
                data, err := conn.Read(buffer)
                if err != nil {
                    log(conn, "connection closed")
                    return
                }

                // Data received from client
                packetString := string(buffer[:data])
                if len(packetString) > 0 {
                    log(conn, "packet: " + packetString)
                    timer.Reset(timeoutSeconds)
                }

                // Send a response back to client
                _, err = conn.Write([]byte("OK"))
                if err != nil {
                    log(conn, "Error writing to connection: " + err.Error())
                    return
                }
            }
        }()

        // Timer to close the connection after timeoutSeconds
        go func() {
            timer = time.NewTimer(timeoutSeconds)
            select {
            case <-timer.C:
                // do something for timeout
                defer conn.Close()
            case <-canceled:
                // timer aborted
                return
            }
        }()
    }
}