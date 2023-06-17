package main

import (
	"github.com/rrebeiz/quicktasks/internal/data"
	"os"
	"testing"
)

var testConfig config
var testApp application

func TestMain(m *testing.M) {
	testConfig.environment = "prod"
	testConfig.port = 4000
	testApp.config = testConfig
	testApp.models = newTestsModel()
	os.Exit(m.Run())
}

func newTestsModel() data.Models {
	return data.Models{
		Tasks: data.NewTaskMockModel(),
	}
}
