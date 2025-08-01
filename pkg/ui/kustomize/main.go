package kustomize

import "github.com/sh31k30ps/gikops/pkg/ui"

type UIKustomizeRequester struct {
	results *UIKustomizeResults
}

type UIKustomizeResults struct {
	URLs []string
}

func NewRequester() *UIKustomizeRequester {
	return &UIKustomizeRequester{
		results: &UIKustomizeResults{},
	}
}

func (ui *UIKustomizeRequester) Request() (ui.UIRequestResult, error) {
	return ui.results, nil
}
