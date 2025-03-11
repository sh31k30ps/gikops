package v1alpha1

func SetProjectDefaults(p *Project) {
	if p == nil {
		return
	}
	p.TypeMeta = TypeMeta{
		APIVersion: Version,
		Kind:       ProjectKind,
	}
	if len(p.Clusters) == 0 {
		kindConfig := &ClusterKindConfig{}
		SetKindConfigDefaults(kindConfig)
		p.Clusters = []Cluster{
			{
				Name:       "local",
				KindConfig: kindConfig,
			},
		}
	}
}
