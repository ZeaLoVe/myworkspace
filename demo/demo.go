package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strconv"
)

func sayHello() {
	fmt.Println("hello")
}
func sayHelloTimes(c *cli.Context) {
	//fmt.Println(c.Args()[0])
	times, _ := strconv.Atoi(c.Args()[0])
	for i := 0; i < times; i++ {
		sayHello()
	}
}

func main() {
	fmt.Println("System start:")
	app := cli.NewApp()
	app.Name = "test"
	app.Usage = "make a test"
	app.Commands = []cli.Command{
		{
			Name:  "say",
			Usage: "say onece",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "say", Value: "", Usage: "say one time"},
			},
			Action: func(c *cli.Context) {
				sayHello()
			},
		},
		{
			Name:  "sayTimes",
			Usage: "say times",
			Flags: []cli.Flag{
				cli.IntFlag{Name: "saytimes", Value: 0, Usage: "say times"},
			},
			Action: func(c *cli.Context) {
				sayHelloTimes(c)
			},
		},
	}
	app.Run(os.Args)

}
