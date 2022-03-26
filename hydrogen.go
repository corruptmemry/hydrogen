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
		os.WriteFile(home+"/.config/hydrogenrpc/config.toml", []byte(`AppID = "857258957587087380"
Port = "6670"`), 0755)
		readConfig()
	}
}

func main() {
	readConfig()
	fmt.Println("Port: " + conf.Port)
	fmt.Println("AppID: " + conf.AppID)
	err := client.Login(conf.AppID)
	if err != nil {
		panic(err)
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
                                State:   song["Title"],
                                Details: "Paused",
                        })
                }
		if status["state"] == "play" {
			err = client.SetActivity(client.Activity{
				State:   song["Title"],
				Details: "Playing",
			})
			if err != nil {
				panic(err)
			}
		} else {
			time.Sleep(1e9)
		}
	}
}
