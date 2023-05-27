package main

import (
	"cloud.google.com/go/run/apiv2/runpb"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/run/apiv2"
	"github.com/urfave/cli/v2"
)

func matchImageName(a, b string) bool {
	return strings.Split(a, ":")[0] == strings.Split(b, ":")[1]
}

func main() {
	app := &cli.App{
		Name:  "gcloud-run-deploy-multi",
		Usage: "update images on Cloud Run service with multiple containers",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "name of the service",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "image",
				Usage:    "image used in the service to update",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := run.NewServicesClient(ctx.Context)
			if err != nil {
				return err
			}

			fmt.Println("Fetching the service")

			service, err := client.GetService(ctx.Context, &runpb.GetServiceRequest{
				Name: ctx.String("name"),
			})
			if err != nil {
				return err
			}

			service.GetTemplate().Revision = ""

			for _, container := range service.GetTemplate().GetContainers() {
				for _, image := range ctx.StringSlice("image") {
					if matchImageName(image, container.GetImage()) {
						container.Image = image
					}
				}
			}

			fmt.Println("Modifying the service")

			future, err := client.UpdateService(ctx.Context, &runpb.UpdateServiceRequest{
				Service: service,
			})
			if err != nil {
				return err
			}

			newService, err := future.Wait(ctx.Context)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", newService)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
