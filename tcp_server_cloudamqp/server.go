package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

// use godot package to load/read the .env file and
// return the value of the key
func getEnv(key string) string {
    // load .env file
    err := godotenv.Load(".env")

    if err != nil {
        panic("Error loading .env file")
    }

    return os.Getenv(key)
}

func log(conn net.Conn, message string) {
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "-", conn.RemoteAddr().String(), "-",  message)
}

func main() {
    addr := ":"+getEnv("GATEWAY_PORT")
    timeoutSeconds := 10 * time.Second

    url := getEnv("CLOUDAMP_URL")
    if url == "" {
        panic("CLOUDAMP_URL is not set")
    }
    connection, _ := amqp.Dial(url)
    defer connection.Close()

    channel, _ := connection.Channel()
    q, _ := channel.QueueDeclare(
        "server", // name
        true,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )

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

                    // Send data to cloudamqp
                    msg := amqp.Publishing{
                        DeliveryMode: 1,
                        // Timestamp:    time.Second,
                        ContentType:  "text/plain",
                        Body:         []byte(packetString),
                    }
                    mandatory, immediate := false, false
                    channel.Publish("", q.Name, mandatory, immediate, msg)

                    timer.Reset(timeoutSeconds)
                }

                // Send a response back to client
                _, err = conn.Write([]byte("OK\n"))
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