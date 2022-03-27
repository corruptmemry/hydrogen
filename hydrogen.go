package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/fhs/gompd/mpd"
	"github.com/hugolgst/rich-go/client"
	"os"
	"time"
)

type Config struct {
	Port    string
	Details string
	AppID   string
}

var conf Config

func readConfig() {
	var home, _ = os.UserHomeDir()
	if _, err := toml.DecodeFile(home+"/.config/hydrogenrpc/config.toml", &conf); err != nil {
		os.Mkdir(home+"/.config/hydrogenrpc", 0755)
		os.WriteFile(home+"/.config/hydrogenrpc/config.toml", []byte(`Details = "Hydrogen RPC"
AppID = "857258957587087380"
Port = "6670"`), 0755)
		readConfig()
	}
}

func login() {
	err := client.Login(conf.AppID)
	if err != nil {
		time.Sleep(time.Millisecond * 10000)
		fmt.Println("No Discord instance was found, retrying...")
		login()
	}
	conn, err2 := mpd.Dial("tcp", "localhost:"+conf.Port)
	if err2 != nil {
		fmt.Println(err2)
	}
	defer conn.Close()
	for {
		status, err3 := conn.Status()
		if err3 != nil {
			fmt.Println(err3)
		}
		song, err := conn.CurrentSong()
		if err != nil {
			fmt.Println(err3)
		}
		if status["state"] == "pause" {
			err = client.SetActivity(client.Activity{
				State:   "🎤  " + song["Artist"],
				Details: "🟨  " + song["Title"],
			})
		}
		if status["state"] == "play" {
			err = client.SetActivity(client.Activity{
				State:   "🎤  " + song["Artist"],
				Details: "🟩  " + song["Title"],
			})
			if err != nil {
				panic(err)
			}
		} else {
			time.Sleep(1e9)
		}
	}
}

func main() {
	readConfig()
	fmt.Println("Details: " + conf.Details)
	fmt.Println("Port: " + conf.Port)
	fmt.Println("AppID: " + conf.AppID)
	login()
}
