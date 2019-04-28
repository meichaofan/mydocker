package main

import (
	"github.com/urfave/cli"
	"github.com/Sirupsen/logrus"
	"os"
)

const usage = `my docker is a simple container runtime implement`

// ./mydocker run -it /bin/bash
// ./mydocker run -it -m 100m -cpushare 512 -cpuset 1 /bin/sh
func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	//设置日志格式
	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil;
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
