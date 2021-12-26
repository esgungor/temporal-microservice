package app

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

func CreateKubernetesDeployment(l LaunchRequest) (LaunchKubernetesResponse, error) {
	fmt.Println("Service creation request created")
	// Call kubeapps from here
	return LaunchKubernetesResponse{
		Name:       l.Name,
		Namespace:  l.Namespace,
		LaunchType: l.LaunchType,
		TheiaPort:  30102,
		RpcPort:    30302,
	}, nil

}

var LaunchQueue = "LAUNCH_QUEUE"

type LaunchRequest struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	LaunchType string `json:"launchType"`
	Operation  string `json:"operation"`
}

type LaunchStatus struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	LaunchType  string `json:"launchType"`
	LaunchState string `json:"launchStatus"`
}

type LaunchKubernetesResponse struct {
	Name       string
	Namespace  string
	LaunchType string
	TheiaPort  int
	RpcPort    int
}

func LaunchWorkflow(ctx workflow.Context, req LaunchRequest) error {
	//Workflow Congfiurations
	launchState := LaunchStatus{
		Name:        req.Name,
		Namespace:   req.Namespace,
		LaunchType:  req.LaunchType,
		LaunchState: "CREATING",
	}
	result := LaunchKubernetesResponse{}
	options := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	//Query Handler
	err := workflow.SetQueryHandler(ctx, "getStatus", func() (LaunchStatus, error) {
		return launchState, nil
	})
	workflow.Sleep(ctx, 1*time.Minute)

	if err != nil {
		panic(err)
	}
	//Execute Creation
	err = workflow.ExecuteActivity(ctx, CreateKubernetesDeployment, req).Get(ctx, &result)
	if err != nil {
		panic(err)
	}
	launchState.LaunchState = "RUNNING"

	// Signaling
	signalVal := LaunchRequest{}
	signalName := "CHANGE_LAUNCH"
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	s := workflow.NewSelector(ctx)
	s.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signalVal)
		workflow.GetLogger(ctx).Info("Received signal!", "Signal", signalName, "value", signalVal)
		if signalVal.Operation == "DELETE" {
			workflow.ExecuteActivity(ctx, DeleteKubernetesDeployment)
			launchState.LaunchState = "DELETED"
		}
		if signalVal.Operation == "UPDATE" {
			launchState.Name = signalVal.Name
			launchState.Namespace = signalVal.Namespace
			launchState.LaunchType = signalVal.LaunchType
		}
	})
	for {
		s.Select(ctx)
		if signalVal.Operation == "DELETE" {
			return nil
		}

	}

}
