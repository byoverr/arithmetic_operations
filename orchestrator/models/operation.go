package models

type Operation struct {
	OperationKind         OperationType `json:"operationKind" validate:"required"`
	DurationInMilliSecond int           `json:"durationInMilliSecond" validate:"duration_in_millisec,required"`
}

func IsAllowedOperation(operationType OperationType) bool {
	return operationType == Addition ||
		operationType == Subtraction ||
		operationType == Multiplication ||
		operationType == Division
}

type OperationType string

const (
	Addition       OperationType = "addition"
	Subtraction                  = "subtraction"
	Multiplication               = "multiplication"
	Division                     = "division"
)

type Status string

const (
	Completed Status = "completed"
	InProcess        = "in process"
	Invalid          = "invalid"
)
