// file: common/process_stage.go

package common

// ProcessStage identifies a named stage or step in a multi-stage process.
//
// It is a string alias type used to represent the various phases or checkpoints
// in a workflow or processing pipeline where errors or specific events may occur.
//
// This type is designed to be flexible and extensible; developers can define
// additional custom stages by declaring new constants of this type as needed.
//
// Example:
//
//	const (
//	    StageInit     ProcessStage = "Init"
//	    StageParse    ProcessStage = "Parse"
//	    StageValidate ProcessStage = "Validate"
//	    StageExecute  ProcessStage = "Execute"
//	    StageFinalize ProcessStage = "Finalize"
//	)
//
// Custom stages can be added in user code:
//
//	const StageCustom ProcessStage = "CustomStage"
type ProcessStage string

// Predefined common stages in the processing lifecycle.
const (
	// StageInit represents the initialization stage.
	StageInit ProcessStage = "Init"

	// StageParse represents the parsing stage.
	StageParse ProcessStage = "Parse"

	// StageValidate represents the validation stage.
	StageValidate ProcessStage = "Validate"

	// StageExecute represents the execution stage.
	StageExecute ProcessStage = "Execute"

	// StageFinalize represents the finalization stage.
	StageFinalize ProcessStage = "Finalize"
)
