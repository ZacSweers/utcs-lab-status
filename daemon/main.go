package main

import (

	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docopt/docopt.go"

)

const Name = "Lab Status"
const Version = "0.0.1"
const Identifier = Name + " " + Version
const RepeaterPort = "8000"

const (
	Reachable = iota
	Unreachable
)

func main() {
	usage := Identifier + `

Usage:
  labstatus [--reachable] <hostfile> [<repeater>]
  labstatus -h | --help
  labstatus --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  --reachable   Check and only add hosts that are reachable.`

	args, err := docopt.Parse(usage, nil, true, Identifier, false)
	fmt.Println(Identifier)

	var hostfile string
	hostfile = args["<hostfile>"].(string)
	fmt.Println("Using", hostfile, "as hostfile.")

	file, err := os.Open(hostfile)
	if err != nil {
		log.Fatal("Error opening hostfile: ", hostfile)
	}

	var hosts []string
	reachability := args["--reachable"].(bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host := scanner.Text()
		if fingerHostIsReachable(host) || !reachability {
			hosts = append(hosts, host)
		} else {
			log.Print("Timeout when checking host: ", host)
		}
	}

	var repeater string
	if value := args["<repeater>"]; value != nil {
		repeater = value.(string)
		if _, err := net.DialTimeout("tcp", repeater+":"+RepeaterPort,
			1*time.Second); err != nil {
			log.Fatal("Unable to connect to repeater server.")
		}
	}

	fmt.Println("Added", len(hosts), "host(s).")
	fmt.Println("Going into daemon mode.")
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				poll(hosts, logfile, repeater)
			}
		}
	}()

	select {}

}

func poll(hosts []string, logfile string, repeater string) {

	fmt.Println("Polling hosts ... ")

	pollTime := time.Now()
	var availableHosts, unavailableHosts, unreachableHosts []string
	for _, host := range hosts {

		// Check for connectivity.
		reachable := fingerHostIsReachable(host)

		// Check user status using finger.
		var localUser string
		var remoteUsers map[string]int = make(map[string]int)
		if reachable {

			cmd := exec.Command("finger", "@"+host)
			output, err := cmd.Output()
			if err != nil {
				log.Print("Error running finger command.")
				break
			} else {
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {

					// User is logged in physically.
					if strings.Contains(line, " :0 ") {
						localUser = strings.Fields(line)[0]
					}

					// User is accessing a virtual terminal.
					if strings.Contains(line, " pt ") {
						user := strings.Fields(line)[0]
						remoteUsers[user] += 1
					}

				}
			}

			if localUser == "" {
				availableHosts = append(availableHosts, host)
			} else {
				unavailableHosts = append(unavailableHosts, host)
			}

		} else {

			unreachableHosts = append(unreachableHosts, host)

		}

	}

	if repeater != "" {
		var availability = map[string][]string{
			"Yes":     availableHosts,
			"No":      unavailableHosts,
			"Offline": unreachableHosts,
		}
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.Encode(availability)
		sendStringToRepeater(repeater, buffer.String())
	}

}

func fingerHostIsReachable(host string) bool {

	_, err := net.DialTimeout("tcp", host+":79", 500*time.Millisecond)
	return err == nil

}

func sendStringToRepeater(repeater string, message string) {

	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Print("Error loading certificate keypair.")
		return
	}

	config := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", repeater+":"+RepeaterPort, &config)
	if err != nil {
		log.Print("Unable to connect to repeater server.")
		return
	}

	defer conn.Close()
	if _, err := io.WriteString(conn, message); err != nil {
		log.Print("Unable to send message to repeater.")
	}

}
