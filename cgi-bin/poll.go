package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {

	fmt.Print("Content-type: text/html\n\n")

	var cacheFile string = "cache.txt"
	if info, err := os.Stat(cacheFile); err == nil {
		var age float64 = time.Now().Sub(info.ModTime()).Seconds()
		if age < 59.0 {
			if cacheData, err := ioutil.ReadFile(cacheFile); err == nil {
				fmt.Print(string(cacheData))
				return
			}
		}
	}

	var hostFile string = "hosts.txt"
	file, err := os.Open(hostFile)
	if err != nil {
		log.Fatal("Error accessing host file: ", hostFile)
	}

	hostInfos := make(map[string]map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host := scanner.Text()
		hostInfos[host] = make(map[string]string)
	}

	rwho := exec.Command("rwho", "-a")
	rwhoOutput, err := rwho.Output()
	if err != nil {
		log.Fatal("Error running rwho command.")
	} else {

		lines := strings.Split(string(rwhoOutput), "\n")
		for _, line := range lines[:len(lines)-1] {

			fields := strings.Fields(line)
			session := strings.SplitN(fields[1], ":", 2)
			var idle string
			if len(fields) > 5 {
				idle = fields[5]
			}

			var user string = fields[0]
			var host string = session[0]
			var terminal string = session[1]
			var idleHours, idleMinutes int
			var idling bool
			if idle != "" {
				idleFields := strings.SplitN(idle, ":", 2)
				idleHours, _ = strconv.Atoi(idleFields[0])
				idleMinutes, _ = strconv.Atoi(idleFields[1])
				if idleHours > 2 && idleMinutes > 0 {
					idling = true
				}
			}

			var terminalType string
			if terminal == ":0" {
				terminalType = "local"
			} else if strings.HasPrefix(terminal, "pts") {
				terminalType = "remote"
			}

			if hostInfo, ok := hostInfos[host]; ok && terminalType == "local" {
				hostInfo["localUser"] = user
				hostInfo["idling"] = strconv.FormatBool(idling)
			}

		}

	}

	ruptime := exec.Command("ruptime")
	ruptimeOutput, err := ruptime.Output()
	if err != nil {
		log.Fatal("Error running ruptime command.")
	} else {

		lines := strings.Split(string(ruptimeOutput), "\n")
		for _, line := range lines[:len(lines)-1] {

			fields := strings.Fields(line)
			var host string = fields[0]
			var online bool = (fields[1] == "up")

			if hostInfo, ok := hostInfos[host]; ok {
				hostInfo["online"] = strconv.FormatBool(online)
				if online {
					var numUsers string = fields[3]
					var load string = fields[7]
					hostInfo["numUsers"] = numUsers
					hostInfo["load"] = load
				}
			}

		}

	}

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.Encode(hostInfos)
	fmt.Print(buffer)

	if ioutil.WriteFile(cacheFile, buffer.Bytes(), 0644) != nil {
		log.Print("Error writing cache file.");
	}

}
