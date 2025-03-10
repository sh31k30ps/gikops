package v1alpha1

func SetComponentDefaults(c *Component) {
	if c == nil {
		return
	}
	c.TypeMeta = TypeMeta{
		APIVersion: Version,
		Kind:       ComponentKind,
	}
	if len(c.EnvironmentAvailability) == 0 {
		c.EnvironmentAvailability = []string{"local"}
	}
	if c.Helm != nil {
		SetHelmConfigDefaults(c.Helm)
	}
	if c.Exec != nil {
		SetComponentExecDefaults(c.Exec)
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

func SetComponentExecDefaults(e *ComponentExec) {
}

func SetComponentFilesDefaults(f *ComponentFiles) {
	if f.CRDs == "" {
		f.CRDs = "crds.yaml"
	}
}
