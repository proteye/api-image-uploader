package main

import (
	"errors"
	"github.com/codegangsta/cli"
	"github.com/proteye/api-image-uploader/service"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
)

func getConfig(c *cli.Context) (service.Config, error) {
	yamlPath := "config-local.yaml"
	config := service.Config{}

	if _, err := os.Stat(yamlPath); err != nil {
		yamlPath = c.GlobalString("config")
		if _, err := os.Stat(yamlPath); err != nil {
			return config, errors.New("config path not valid")
		}
	}

	ymlData, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(ymlData), &config)
	return config, err
}

func main() {

	app := cli.NewApp()
	app.Name = "Image Uploader RESTful API"
	app.Usage = "upload images to server"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", Value: "config.yaml", Usage: "config file to use", EnvVar: "APP_CONFIG"},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the http server",
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}

				svc := service.ImageUploaderService{}

				if err = svc.Run(cfg); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "migratedb",
			Usage: "Perform database migrations",
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}

				svc := service.ImageUploaderService{}

				if err = svc.Migrate(cfg); err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)

}
