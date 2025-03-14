package component

type FilesConfig struct {
	CRDs     string
	SkipCRDs bool
	Keep     []string
}

func NewFilesConfig() *FilesConfig {
	return &FilesConfig{
		CRDs:     "",
		SkipCRDs: false,
		Keep:     []string{},
	}
}
