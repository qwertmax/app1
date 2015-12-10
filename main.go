package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
)

type App struct {
	Service string `json:"service"`
	Host    string `json:"host"`
	IP      string `json:"ip"`
	Port    string `json:"port"`
}

func GetEndpoint(name string) []App {
	reader := strings.NewReader("")
	request, err := http.NewRequest("GET", "http://52.34.228.148:8123/v1/services/_"+name+"._tcp.marathon.mesos", reader)
	if err != nil {
		panic(err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var apps []App
	json.Unmarshal(body, &apps)

	return apps
}

func Route(apps []App) ([]byte, error) {
	item := rand.Intn(len(apps))
	url := "http://" + apps[item].IP + ":" + apps[item].Port

	reader := strings.NewReader("")
	request, err := http.NewRequest("GET", url, reader)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Main of app1\n")

		addrs, _ := net.InterfaceAddrs()
		for i, addr := range addrs {
			fmt.Fprintf(w, "%d %v\n", i, addr)
		}
	})

	http.HandleFunc("/from2", func(w http.ResponseWriter, req *http.Request) {

		app2 := GetEndpoint("maxapp2")
		resp, err := Route(app2)
		if err != nil {
			fmt.Fprintf(w, "error: %s\n", err.Error())
		}

		fmt.Fprintf(w, "From2 of app1\n %s\n", resp)
	})

	println("ready")
	http.ListenAndServe(":9091", nil)
}
