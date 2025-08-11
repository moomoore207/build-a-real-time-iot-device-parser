package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type IoTDeviceData struct {
	DeviceID string `json:"device_id"`
	Temperature float64 `json:"temperature"`
	Humidity float64 `json:"humidity"`
	Timestamp int64 `json:"timestamp"`
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("IoT device parser listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		var data IoTDeviceData
		err = json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("Received data from device %s: temperature=%.2f°C, humidity=%.2f%%, timestamp=%d\n",
			data.DeviceID, data.Temperature, data.Humidity, data.Timestamp)

		// Process the data here (e.g., store it in a database, send it to a cloud service, etc.)
		processData(data)
	}
}

func processData(data IoTDeviceData) {
	// Example: store the data in a log file
	logFile, err := os.OpenFile("iot_data.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer logFile.Close()

	fmt.Fprintf(logFile, "Device %s: temperature=%.2f°C, humidity=%.2f%%, timestamp=%d\n",
		data.DeviceID, data.Temperature, data.Humidity, data.Timestamp)
}