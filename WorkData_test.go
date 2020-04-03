package main

import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/cadence/testsuite"
	"testing"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) Test_InfWorkflowTest (){
	s.env = s.NewTestWorkflowEnvironment()

	s.env.ExecuteWorkflow(InfWorkFlow, "02.02.2020")

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}