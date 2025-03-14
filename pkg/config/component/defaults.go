package component

func SetComponentDefaults(c *Component) {
	if c.Files == nil {
		c.Files = NewFilesConfig()
	}
	SetFilesConfigDefaults(c.Files)
}

func SetFilesConfigDefaults(f *FilesConfig) {
	if f.CRDs == "" {
		f.CRDs = DefaultCRDsFileName
	}
}
