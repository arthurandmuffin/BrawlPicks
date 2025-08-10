package scheduler

import "context"

type Scheduler interface {
	Start(context.Context)
}
