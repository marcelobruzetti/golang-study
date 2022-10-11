package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const softwareVersion = 0.1
const monitoringTimes = 3
const delay = 5

func main(){
	showIntroduction()

	for {
		showMenu()
		command := readCommand();
		
		switch command {
			case 1:
				startMonitoring()
			case 2:
				showLog()
			case 3:
				eraseLog()
			case 0: 
				fmt.Println("Exiting...")
				os.Exit(0)
			default: 
				fmt.Println("Invalid command")
				os.Exit(-1)
		}
	}
}

func showIntroduction() {
	fmt.Println("--------------------------------------------")
	fmt.Printf("|          Website Monitor v%.1f            |\n", softwareVersion)
	fmt.Println("--------------------------------------------")
}

func showMenu() {
	fmt.Println("")
	fmt.Println("Choose an option:")
	fmt.Println("1 - Start monitoring")
	fmt.Println("2 - Show logs")
	fmt.Println("3 - Erase logs")
	fmt.Println("0 - Exit")
}

func readCommand() int {
	var command int
	fmt.Scan(&command)
	return command
}

func startMonitoring() {
	// fmt.Println("Monitoring...")
	sites := readSitesFromFile()

	for i := 0; i < monitoringTimes; i++ {
		for _, site := range sites {
			testSite(site)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
}

func testSite(site string)  {
	// fmt.Println("Testing site...")
	response, error := http.Get(site)

	if error != nil {
		fmt.Println("Error:", error)
	}

	if response.StatusCode == 200 {
		fmt.Println("Online -", site)
		log(site, true)
	} else {
		fmt.Println("Offline", response.StatusCode,  "-", site)
		log(site, false)
	}
}

func readSitesFromFile() []string {
	var sites []string
	
	file, error := os.Open("websites.txt")
	if error != nil {
		fmt.Println("Error:", error)
	}

	reader := bufio.NewReader(file)

	for {
		line, error := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		sites = append(sites, line)

		if error == io.EOF {
			break;
		}
	}

	file.Close()

	return sites
}

func log(site string, status bool) {
	file, error := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if error != nil {
		fmt.Println("Error:", error)
	}

	var siteStatus string
	if status {
		siteStatus = "Online"
	} else {
		siteStatus = "Offline"
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05")+" "+site+" - "+siteStatus+"\n")

	file.Close()
}

func showLog() {
	file, error := ioutil.ReadFile("log.txt")
	if error != nil {
		fmt.Println("Error:", error)
	}

	fmt.Println(string(file))
}

func eraseLog() {
	error := os.Remove("log.txt")
	if error != nil {
		fmt.Println("Error:", error)
	}
	
	file, error := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE, 0666)
	if error != nil {
		fmt.Println("Error:", error)
	}
	file.Close()
	fmt.Println("Log erased")
}