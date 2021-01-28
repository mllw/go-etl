package plugin

import (
	"github.com/Breeze0806/go-etl/element"
)

type RecordSender interface {
	CreateRecord() (element.Record, error)
	SendWriter(record element.Record) error
	Flush() error
	Terminate() error
	Shutdown() error
}