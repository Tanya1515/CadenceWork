package wf

import (
	//"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"math/rand"
	"time"
)

type InfoData struct {
	Data string
	Status string
}

func RandomStatus() (string){
	var x int
	x = rand.Intn(100)
	if (x >= 50){
		return "Ok"
}
	return "Fail"

}

func InfoStruct (Data string) (*InfoData, error){
	info := &InfoData{Data: Data, Status: RandomStatus()}
	return info, nil
}

func InfoStructResult (info *InfoData) (string, error){
	if info.Status == "Ok" {
		return "Today is " + info.Data, nil
	}
	return "Fail", nil
}


func InfWorkFlow (ctx workflow.Context, Data string) (error) {

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: 3*time.Second,
		StartToCloseTimeout: 3*time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var Info *InfoData
	err := workflow.ExecuteActivity(ctx, InfoStruct, Data).Get(ctx, &Info)
	if err != nil{
		return err
	}

	var result string
	err = workflow.ExecuteActivity(ctx, InfoStructResult, Info).Get(ctx, &result)
	if err != nil{
		return err
	}

	workflow.GetLogger(ctx).Info("Result is " + result)
	return nil
}

