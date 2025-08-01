package dependencies

import (
	"fmt"
	"sort"
	"strings"

	componentConfig "github.com/sh31k30ps/gikops/pkg/config/component"
)

type DependencyGraph struct {
	dependencies        map[string][]string
	reverseDependencies map[string][]string
	errors              []error
}

func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		dependencies:        make(map[string][]string),
		reverseDependencies: make(map[string][]string),
		errors:              []error{},
	}
}

type ComponentGetter func(name string) (*componentConfig.Component, error)

func (dg *DependencyGraph) ReverseResolve(component string, existingComponents []string, getter ComponentGetter) ([]string, []error) {
	if _, errs := dg.Resolve(existingComponents, getter); len(errs) > 0 {
		return nil, errs
	}
	return dg.reverseDependencies[component], nil
}

func (dg *DependencyGraph) Resolve(components []string, getter ComponentGetter) ([]string, []error) {
	var result []string
	visited := make(map[string]bool)
	visiting := make(map[string]bool) // Pour détecter les dépendances cycliques
	errorComponents := []string{}

	// Clean desabled components
	components = cleanDisabledComponents(components, getter)

	var resolve func(string) error
	resolve = func(component string) error {
		if visited[component] {
			return nil
		}
		if visiting[component] {
			circularDependency := fmt.Errorf("circular dependency detected: %s", component)
			dg.addError(circularDependency)
			return circularDependency
		}
		visiting[component] = true
		defer func() { visiting[component] = false }()

		config, err := getter(component)
		if err != nil {
			dg.addError(fmt.Errorf("error while fetching %s: %w", component, err))
			return err
		}

		if config.Disabled {
			cfgErr := fmt.Errorf("component %s is disabled", component)
			dg.addError(cfgErr)
			errorComponents = append(errorComponents, component)
			return cfgErr
		}

		var dependencies []string
		for _, dep := range config.DependsOn {
			if !containsSlash(dep) {
				dep = componentConfig.GetComponentPrefix(component) + "/" + dep
			}

			if _, err := getter(dep); err != nil {
				dg.addError(fmt.Errorf("dependency not found: %s", dep))
				errorComponents = append(errorComponents, dep)
				continue
			}

			dependencies = append(dependencies, dep)
			if err := resolve(dep); err != nil {
				return err
			}
		}
		dg.addDependency(component, dependencies)
		dg.addReverseDependencies(component, dependencies)
		visited[component] = true
		result = append(result, component)
		return nil
	}

	for _, component := range components {
		resolve(component)
	}

	sort.Slice(result, func(i, j int) bool {
		// Les composants avec le moins de dépendances sont plus critiques
		if len(dg.dependencies[result[i]]) != len(dg.dependencies[result[j]]) {
			return len(dg.dependencies[result[i]]) < len(dg.dependencies[result[j]])
		}
		// À nombre égal de dépendances, on regarde si l'un dépend de l'autre
		for _, dep := range dg.dependencies[result[i]] {
			if dep == result[j] {
				return false // i dépend de j donc j est plus critique
			}
		}
		for _, dep := range dg.dependencies[result[j]] {
			if dep == result[i] {
				return true // j dépend de i donc i est plus critique
			}
		}
		return false
	})

	return result, dg.errors
}

func (dg *DependencyGraph) addDependency(component string, dependsOn []string) {
	dg.dependencies[component] = dependsOn
}

func (dg *DependencyGraph) addReverseDependencies(component string, dependsOn []string) {
	for _, dep := range dependsOn {
		if _, ok := dg.reverseDependencies[dep]; !ok {
			dg.reverseDependencies[dep] = []string{}
		}
		dg.reverseDependencies[dep] = append(dg.reverseDependencies[dep], component)
	}
}

func (dg *DependencyGraph) addError(err error) {
	dg.errors = append(dg.errors, err)
}

func cleanDisabledComponents(components []string, getter ComponentGetter) []string {
	cleaned := []string{}
	for _, component := range components {
		if config, err := getter(component); err == nil && !config.Disabled {
			cleaned = append(cleaned, component)
		}
	}
	return cleaned
}

func containsSlash(s string) bool {
	return strings.Contains(s, "/")
}
