package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/esgungor/temporal-microservice/app"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.temporal.io/api/filter/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type Workflow struct {
	WorkflowID string `json:"workflows_id"`
	RunID      string `json:"run_id"`
}
type WorkflowList []Workflow

var c client.Client

func updateLaunch(w http.ResponseWriter, r *http.Request) {
	ids := mux.Vars(r)

	err := c.SignalWorkflow(context.Background(), ids["id"], ids["runId"], "CHANGE_LAUNCH", "UPDATE")
	if err != nil {
		panic(err)
	}
	w.Write([]byte("Update operation completed!\nWorkflowID: " + ids["id"] + "\nRunID:" + ids["runId"]))
}
func deleteLaunch(w http.ResponseWriter, r *http.Request) {
	ids := mux.Vars(r)

	c, err := client.NewClient(client.Options{
		HostPort: "23.88.62.179:31313",
	})
	if err != nil {
		panic(err)
	}

	err = c.SignalWorkflow(context.Background(), ids["id"], ids["runId"], "CHANGE_LAUNCH", "DELETE")
	if err != nil {
		panic(err)
	}
	w.Write([]byte("Deletion completed!\nWorkflowID: " + ids["id"] + "\nRunID:" + ids["runId"]))
}

func getLaunch(w http.ResponseWriter, r *http.Request) {
	ids := mux.Vars(r)
	fmt.Println(ids)

	resp, err := c.QueryWorkflow(context.Background(), ids["id"], ids["runId"], "getStatus")
	if err != nil {
		panic(err)
	}
	var result app.LaunchStatus
	if err := resp.Get(&result); err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(&result)
}

func listLaunches(w http.ResponseWriter, r *http.Request) {

	resp, err := c.ListOpenWorkflow(context.Background(), &workflowservice.ListOpenWorkflowExecutionsRequest{
		Namespace: "default",
		Filters: &workflowservice.ListOpenWorkflowExecutionsRequest_TypeFilter{
			TypeFilter: &filter.WorkflowTypeFilter{
				Name: "LaunchWorkflow",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	list := printExecutions(resp)
	err = json.NewEncoder(w).Encode(&list)
	if err != nil {
		panic(err)
	}
}

func createLaunch(w http.ResponseWriter, r *http.Request) {
	var body app.LaunchRequest
	json.NewDecoder(r.Body).Decode(&body)
	options := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: app.LaunchQueue,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, app.LaunchWorkflow, body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Launch operation started. WorkflowID:%v RunID:%v\n", we.GetID(), we.GetRunID())
	json.NewEncoder(w).Encode(&Workflow{
		WorkflowID: we.GetID(),
		RunID:      we.GetRunID(),
	})
}

func main() {
	var err error
	c, err = client.NewClient(client.Options{
		HostPort: os.Getenv("TEMPORAL_SERVER_IP"),
	})
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/launch/{id}/{runId}", getLaunch).Methods("GET")
	r.HandleFunc("/launch/{id}/{runId}", deleteLaunch).Methods("DELETE")
	r.HandleFunc("/launch/{id}/{runId}", updateLaunch).Methods("PATCH")

	r.HandleFunc("/", listLaunches).Methods("GET")
	r.HandleFunc("/", createLaunch).Methods("POST")

	fmt.Println("Connection created.")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":5000", nil), "Server listening...")
}

func printExecutions(r *workflowservice.ListOpenWorkflowExecutionsResponse) WorkflowList {

	workflowList := WorkflowList{}
	for _, workflows := range r.Executions {
		workflowList = append(workflowList, Workflow{
			WorkflowID: workflows.Execution.GetWorkflowId(),
			RunID:      workflows.Execution.GetRunId(),
		})

	}
	return workflowList
}
