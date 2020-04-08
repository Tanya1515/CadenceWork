package main

import (
	"github.com/uber-go/tally"
	_ "github.com/uber-go/tally"
	_ "go.uber.org/cadence/.gen/go/cadence"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	_ "go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"
	_ "go.uber.org/cadence/worker"
	"go.uber.org/yarpc"
	_ "go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/transport/tchannel"
	_ "go.uber.org/zap"
	_ "go.uber.org/zap/zapcore"
)
const (
	TaskListName = "tasklist"
	Domain = "domain"
	ClientName = "worker"
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
		MetricsScope: tally.NewTestScope(TaskListName, map[string]string{}),
	}

	worker := worker.New(
		service,
		Domain,
		TaskListName,
		workerOptions)
	err := worker.Start()
	if err != nil {
		panic("Failed to start worker")
	}

}

func main() {
	startWorker(CreateClient())
}