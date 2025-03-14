package v1alpha1

func SetComponentDefaults(c *Component) {
	if c == nil {
		return
	}
	c.TypeMeta = TypeMeta{
		APIVersion: Version,
		Kind:       ComponentKind,
	}
	if c.Helm != nil {
		SetHelmConfigDefaults(c.Helm)
	}
	if c.Kustomize != nil {
		SetKustomizeConfigDefaults(c.Kustomize)
	}
	if c.Files == nil {
		c.Files = &ComponentFiles{}
	}
	SetComponentFilesDefaults(c.Files)
}

func SetHelmConfigDefaults(h *HelmConfig) {
	if h.Before == nil {
		h.Before = &HelmBeforeInitConfig{}
	}
	if h.After == nil {
		h.After = &HelmAfterInitConfig{}
	}
}

func SetComponentFilesDefaults(f *ComponentFiles) {
	if f.CRDs == "" {
		f.CRDs = "crds.yaml"
	}
}

func SetKustomizeConfigDefaults(k *KustomizeConfig) {
	if k == nil {
		return
	}
	if k.URLs == nil {
		k.URLs = []string{}
	}
}
