package main

import (
	"awesomeProject/wf"
	"context"
	"fmt"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/worker"
	"github.com/uber-go/tally"
	_ "github.com/uber-go/tally"
	"go.uber.org/cadence/client"
	"time"

	_ "go.uber.org/cadence/.gen/go/cadence"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	_ "go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	//"go.uber.org/cadence/worker"
	_ "go.uber.org/cadence/worker"
	"go.uber.org/yarpc"
	_ "go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/transport/tchannel"
	_ "go.uber.org/zap"
	_ "go.uber.org/zap/zapcore"
)
const (
	TaskList = "tasklist"
	Domain = "samples-domain"
	ClientName = "sample-worker"
	HostPort = "127.0.0.1:7933"
	CadenceService = "cadence-frontend"
)

func CreateClient() workflowserviceclient.Interface{
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
	if err != nil{
		panic("Failed to setup tchannel")
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientName,
		Outbounds: yarpc.Outbounds{
			CadenceService : {Unary: ch.NewSingleOutbound(HostPort)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		panic("Failed to start dispatcher")
	}
	return workflowserviceclient.New(dispatcher.ClientConfig(CadenceService))
}

func startWorker(service workflowserviceclient.Interface) {
	workerOptions := worker.Options{
		MetricsScope: tally.NewTestScope(TaskList, map[string]string{}),
	}

	worker := worker.New(
		service,
		Domain,
		TaskList,
		workerOptions)
	err := worker.Start()
	if err != nil {
		panic("Failed to start worker")
	}

}

func StartWork () {
	cadenceClient := CreateClient()
	newclient := client.NewClient(cadenceClient, Domain, &client.Options{})
	workflowOptions := client.StartWorkflowOptions{
		ID:									fmt.Sprintf("workflow %v", uuid.New()),
		TaskList:							TaskList,
		ExecutionStartToCloseTimeout:		time.Minute,
		DecisionTaskStartToCloseTimeout:	time.Second *3,
	}
	work, err := newclient.StartWorkflow( context.Background(), workflowOptions, wf.InfWorkFlow, "02.02.2020")
	if err != nil {
		fmt.Printf("Failed to start workflow")
		return
	}
	fmt.Printf("Workflow started %v", work.ID)
}

func main() {
	startWorker(CreateClient())
	StartWork()
}