package errors_test

import (
	"errors"
	"fmt"

	"github.com/entiqon/entiqon/common"
	enterrors "github.com/entiqon/entiqon/common/errors"
)

func ExampleNewCausableError() {
	err := enterrors.NewCausableError("Database", "Connection failed")

	fmt.Println(err.Cause())
	fmt.Println(err.Reason())
	fmt.Println(err.Error())
	// Output:
	// Database
	// Connection failed
	// Connection failed
}

func ExampleNewProcessStageError() {
	stage := "Init"
	cause := "Loader"
	reason := "Failed to load resource"

	err := enterrors.NewProcessStageError(common.ProcessStage(stage), cause, reason, errors.New("file not found"))

	fmt.Println(err.Error())
	// Output:
	// [Loader] at stage Init: Failed to load resource: file not found
}
