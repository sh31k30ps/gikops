package component

import "errors"

func ValidateComponent(c Component) []error {
	var errs []error
	if c.Name == "" {
		errs = append(errs, errors.New("component name is required"))
	}
	if c.Files != nil {
		if c.Files.CRDs == "" {
			errs = append(errs, errors.New("component files CRDs is required"))
		}
	}
	if c.Helm != nil {
		errs = append(errs, ValidateHelmConfig(*c.Helm)...)
	}
	return errs
}

func ValidateHelmConfig(h HelmConfig) []error {
	var errs []error
	if h.Repo == "" {
		errs = append(errs, errors.New("helm repo is required"))
	}
	if h.URL == "" {
		errs = append(errs, errors.New("helm url is required"))
	}
	if h.Chart == nil {
		errs = append(errs, errors.New("helm chart is required"))
		return errs
	}
	errs = append(errs, ValidateHelmChart(*h.Chart)...)
	if h.CRDsChart != nil {
		errs = append(errs, ValidateHelmChart(*h.CRDsChart)...)
	}
	if h.Before != nil {
		errs = append(errs, ValidateHelmInitHooks(*h.Before)...)
	}
	if h.After != nil {
		errs = append(errs, ValidateHelmInitHooks(*h.After)...)
	}

	return errs
}

func ValidateHelmChart(h HelmChart) []error {
	var errs []error
	if h.Chart == "" {
		errs = append(errs, errors.New("helm chart is required"))
	}
	return errs
}

func ValidateHelmInitHooks(h HelmInitHooks) []error {
	var errs []error
	if len(h.Uploads) > 0 {
		for _, upload := range h.Uploads {
			if upload.Name == "" {
				errs = append(errs, errors.New("helm init hooks upload name is required"))
			}
			if upload.URL == "" {
				errs = append(errs, errors.New("helm init hooks upload url is required"))
			}
		}
	}
	if len(h.Resolves) > 0 {
		for _, resolve := range h.Resolves {
			if resolve == "" {
				errs = append(errs, errors.New("helm init hooks resolve is required"))
			}
		}
	}
	if len(h.Renames) > 0 {
		for _, rename := range h.Renames {
			if rename.Original == "" {
				errs = append(errs, errors.New("helm init hooks rename original is required"))
			}
			if rename.Renamed == "" {
				errs = append(errs, errors.New("helm init hooks renamed is required"))
			}
		}
	}
	if len(h.Concats) > 0 {
		for _, concat := range h.Concats {
			if concat.Folder == "" {
				errs = append(errs, errors.New("helm init hooks concat folder is required"))
			}
			if concat.Output == "" {
				errs = append(errs, errors.New("helm init hooks concat output is required"))
			}
		}
	}
	return errs
}
