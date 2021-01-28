package reader

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

type Job interface {
	plugin.Job
	Split(ctx context.Context, number int) ([]*config.JSON, error)
}