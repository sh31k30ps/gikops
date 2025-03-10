package component

func SetComponentDefaults(c *Component) {
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
		c.Files = NewComponentFiles()
	}
	SetComponentFilesDefaults(c.Files)
}

func SetComponentFilesDefaults(f *ComponentFiles) {
	if f.CRDs == "" {
		f.CRDs = DefaultCRDsFileName
	}
}

func SetComponentExecDefaults(e *ComponentExec) {
}

func SetHelmConfigDefaults(h *HelmConfig) {
}
