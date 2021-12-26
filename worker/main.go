package main

import (
	"log"
	"os"

	"github.com/esgungor/temporal-microservice/app"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	c, err := client.NewClient(client.Options{
		HostPort: os.Getenv("TEMPORAL_SERVER_IP"),
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	w := worker.New(c, app.LaunchQueue, worker.Options{})
	w.RegisterWorkflow(app.LaunchWorkflow)
	w.RegisterActivity(app.CreateKubernetesDeployment)
	w.RegisterActivity(app.DeleteKubernetesDeployment)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Worker cannot start")
	}
}
