package v1alpha1

// Upload specifies a file to be uploaded
type Upload struct {
	// Name is the target name for the uploaded file
	Name string `json:"name" yaml:"name"`

	// URL is the source URL to download from
	URL string `json:"url" yaml:"url"`
}

// Rename specifies a file renaming operation
type Rename struct {
	// Original is the original file name
	Original string `json:"original" yaml:"original"`

	// Renamed is the new file name
	Renamed string `json:"renamed" yaml:"renamed"`
}

// Concat specifies a file concatenation operation
type Concat struct {
	// Folder is the directory containing source files
	Folder string `json:"folder" yaml:"folder,omitempty"`

	// Includes is a pattern matching files to include
	Includes []string `json:"includes,omitempty" yaml:"includes,omitempty"`

	// Output is the target concatenated file
	Output string `json:"output" yaml:"output"`
}
