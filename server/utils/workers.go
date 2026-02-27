package utils

import (
	"context"
	"sync"
)

type WorkGroup struct {
	Ctx 	context.Context
	Cancel 	context.CancelFunc
	WG 		sync.WaitGroup
}

func NewWorkGroup(count ...int) *WorkGroup {
	initial_cout := 0;
	if len(count) > 0 {
		initial_cout = count[0];
	}
	ctx, cancel := context.WithCancel(context.Background());

	group := &WorkGroup{
		Ctx: ctx,
		Cancel: cancel,
		WG: sync.WaitGroup{},
	};
	group.WG.Add(initial_cout);

	return group;
}
