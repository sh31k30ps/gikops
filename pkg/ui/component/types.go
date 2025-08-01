package component

import (
	"github.com/sh31k30ps/gikops/pkg/ui/helm"
	"github.com/sh31k30ps/gikops/pkg/ui/kustomize"
)

type UIComponentResults struct {
	Name      string
	Namespace string
	Disabled  bool
	DependsOn []string
	Helm      *helm.UIHelmResults
	Kustomize *kustomize.UIKustomizeResults
	Clusters  []string
}
