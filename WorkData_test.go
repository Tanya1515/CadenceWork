package main

import (
	"awesomeProject/wf"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/cadence/testsuite"
	"go.uber.org/cadence/workflow"
	"testing"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) Test_Activity_InfoStruct_err(){
	s.env.OnActivity(wf.InfoStruct, mock.Anything ).Return(nil, errors.New("InfoStructFailure"))
	s.env.ExecuteWorkflow(wf.InfWorkFlow, "01.04.2020")

	s.True(s.env.IsWorkflowCompleted())

	s.NotNil(s.env.GetWorkflowError())
	_, ok := s.env.GetWorkflowError().(*workflow.GenericError)
	s.True(ok)
	s.Equal("InfoStructFailure",s.env.GetWorkflowError().Error())
}

func (s *UnitTestSuite) Test_Activity_InfoStructResult_err(){ //проверяет возвращение параметров при выполнении InfoStructResult
	s.env.OnActivity(wf.InfoStructResult, mock.Anything ).Return("", errors.New("InfoStructResultFailure"))
	s.env.ExecuteWorkflow(wf.InfWorkFlow, "01.04.2020")

	s.True(s.env.IsWorkflowCompleted())

	s.NotNil(s.env.GetWorkflowError())
	_, ok := s.env.GetWorkflowError().(*workflow.GenericError)
	s.True(ok)
	s.Equal("InfoStructResultFailure",s.env.GetWorkflowError().Error())
}

func (s *UnitTestSuite) Test_Activity_InfoStruct_param(){
	s.env.OnActivity(wf.InfoStruct, mock.Anything).Return(
			func (Data string) (*wf.InfoData, error){
				s.Equal("02.02.2020", Data)
				info := &wf.InfoData{Data: Data, Status: wf.RandomStatus()}
				return info, nil
			})
	s.env.ExecuteWorkflow(wf.InfWorkFlow, "02.02.2020")

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_Activity_InfoStructResult_param(){
	s.env.OnActivity(wf.InfoStructResult, mock.Anything).Return(
		func (info *wf.InfoData) (string, error){
			s.Equal("02.02.2020", info.Data)
			if info.Status == "Ok" {
				return "Today is " + info.Data, nil
			}
			return "Fail", nil
		})
	s.env.ExecuteWorkflow(wf.InfWorkFlow, "02.02.2020")

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) Test_InfWorkflowTest (){

	s.env.ExecuteWorkflow(wf.InfWorkFlow, "02.02.2020")

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}