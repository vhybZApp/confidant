package main

import (
	"fmt"
	"os"

	"github.com/farhoud/confidant/internal/config"
	"github.com/farhoud/confidant/internal/mind"
)

func main() {
	conf := config.Configuration(config.WithDotEnvConfig)

	args := os.Args[1:]

	fmt.Printf("start app with configs: %+v\n", conf)
	mvi := mind.NewRobotScreenInspector()
	m := mind.NewApp(*conf, mvi)
	err := m.Run(args[0])
	if err != nil {
		panic(err)
	}
}
