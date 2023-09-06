package model

type OpType string

const (
	OpTypeSet    OpType = "SET"
	OpTypeGet    OpType = "GET"
	OpTypeDelete OpType = "DELETE"
)
