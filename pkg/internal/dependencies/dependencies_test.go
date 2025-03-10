package dependencies

import (
	"errors"
	"testing"

	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/stretchr/testify/assert"
)

// MockComponentGetter est utilisé pour simuler pkg.GetComponent
type MockComponentGetter struct {
	components map[string]*component.Component
}

func (m *MockComponentGetter) GetComponent(name string) (*component.Component, error) {
	if comp, exists := m.components[name]; exists {
		return comp, nil
	}
	return nil, errors.New("component not found")
}

func TestDependencyGraph_Resolve(t *testing.T) {
	tests := []struct {
		name           string
		components     []string
		mockComponents map[string]*component.Component
		expectedOrder  []string
		expectedErrors int
	}{
		{
			name:       "sans dépendances",
			components: []string{"core/comp1", "core/comp2"},
			mockComponents: map[string]*component.Component{
				"core/comp1": {DependsOn: []string{}},
				"core/comp2": {DependsOn: []string{}},
			},
			expectedOrder:  []string{"core/comp1", "core/comp2"},
			expectedErrors: 0,
		},
		{
			name:       "avec dépendances simples",
			components: []string{"core/comp2", "core/comp1"},
			mockComponents: map[string]*component.Component{
				"core/comp1": {DependsOn: []string{}},
				"core/comp2": {DependsOn: []string{"comp1"}},
			},
			expectedOrder:  []string{"core/comp1", "core/comp2"},
			expectedErrors: 0,
		},
		{
			name:       "avec dépendances complexes",
			components: []string{"core/comp1", "core/comp3", "core/comp2"},
			mockComponents: map[string]*component.Component{
				"core/comp3": {DependsOn: []string{"comp1", "comp2"}},
				"core/comp1": {DependsOn: []string{}},
				"core/comp2": {DependsOn: []string{"comp1"}},
			},
			expectedOrder:  []string{"core/comp1", "core/comp2", "core/comp3"},
			expectedErrors: 0,
		},
		{
			name:       "avec dépendance manquante",
			components: []string{"core/comp1"},
			mockComponents: map[string]*component.Component{
				"core/comp1": {DependsOn: []string{"missing"}},
			},
			expectedOrder:  []string{"core/comp1"},
			expectedErrors: 1,
		},
		{
			name:       "avec chemin complet dans les dépendances",
			components: []string{"core/comp1", "other/comp2"},
			mockComponents: map[string]*component.Component{
				"core/comp1":  {DependsOn: []string{"other/comp2"}},
				"other/comp2": {DependsOn: []string{}},
			},
			expectedOrder:  []string{"other/comp2", "core/comp1"},
			expectedErrors: 0,
		},
		{
			name:       "avec dépendances circulaires",
			components: []string{"core/comp1", "other/comp2"},
			mockComponents: map[string]*component.Component{
				"core/comp1":  {DependsOn: []string{"other/comp2"}},
				"other/comp2": {DependsOn: []string{"core/comp1"}},
			},
			expectedOrder:  nil,
			expectedErrors: 2,
		},
		{
			name:       "avec dépendance non demandée",
			components: []string{"other/comp2"},
			mockComponents: map[string]*component.Component{
				"other/comp2": {DependsOn: []string{"core/comp1"}},
				"core/comp1":  {DependsOn: []string{"comp2"}},
				"core/comp2":  {DependsOn: []string{}},
			},
			expectedOrder:  []string{"core/comp2", "core/comp1", "other/comp2"},
			expectedErrors: 0,
		},
		{
			name:       "Avec dépendance disabled",
			components: []string{"core/comp1", "core/comp2"},
			mockComponents: map[string]*component.Component{
				"core/comp1": {DependsOn: []string{}},
				"core/comp2": {DependsOn: []string{}, Disabled: true},
			},
			expectedOrder:  []string{"core/comp1"},
			expectedErrors: 0,
		},
		{
			name:       "Avec dépendance disabled dans les dépendances",
			components: []string{"core/comp1"},
			mockComponents: map[string]*component.Component{
				"core/comp1": {DependsOn: []string{"core/comp2"}},
				"core/comp2": {DependsOn: []string{}, Disabled: true},
			},
			expectedOrder:  nil,
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Remplacer temporairement pkg.GetComponent par notre mock
			mockGetComponent := func(name string) (*component.Component, error) {
				if comp, exists := tt.mockComponents[name]; exists {
					return comp, nil
				}
				return nil, errors.New("component not found")
			}

			// Créer et résoudre le graphe de dépendances
			dg := NewDependencyGraph()
			order, errs := dg.Resolve(tt.components, mockGetComponent)

			// Vérifier les résultats
			assert.Len(t, errs, tt.expectedErrors, "nombre d'erreurs incorrect")
			assert.Equal(t, tt.expectedOrder, order, "ordre de résolution incorrect")
		})
	}
}

func TestContainsSlash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "avec slash",
			input:    "core/comp1",
			expected: true,
		},
		{
			name:     "sans slash",
			input:    "comp1",
			expected: false,
		},
		{
			name:     "chaîne vide",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsSlash(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDependencyGraph_AddDependency(t *testing.T) {
	dg := NewDependencyGraph()

	tests := []struct {
		name       string
		component  string
		dependsOn  []string
		checkAfter func(*testing.T, *DependencyGraph)
	}{
		{
			name:      "ajouter une dépendance simple",
			component: "comp1",
			dependsOn: []string{"comp2"},
			checkAfter: func(t *testing.T, dg *DependencyGraph) {
				deps, exists := dg.dependencies["comp1"]
				assert.True(t, exists)
				assert.Equal(t, []string{"comp2"}, deps)
			},
		},
		{
			name:      "ajouter plusieurs dépendances",
			component: "comp3",
			dependsOn: []string{"comp1", "comp2"},
			checkAfter: func(t *testing.T, dg *DependencyGraph) {
				deps, exists := dg.dependencies["comp3"]
				assert.True(t, exists)
				assert.Equal(t, []string{"comp1", "comp2"}, deps)
			},
		},
		{
			name:      "sans dépendances",
			component: "comp4",
			dependsOn: []string{},
			checkAfter: func(t *testing.T, dg *DependencyGraph) {
				deps, exists := dg.dependencies["comp4"]
				assert.True(t, exists)
				assert.Empty(t, deps)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dg.addDependency(tt.component, tt.dependsOn)
			if tt.checkAfter != nil {
				tt.checkAfter(t, dg)
			}
		})
	}
}

func TestDependencyGraph_AddError(t *testing.T) {
	dg := NewDependencyGraph()

	// Test l'ajout d'une erreur
	err1 := errors.New("erreur 1")
	dg.addError(err1)
	assert.Len(t, dg.errors, 1)
	assert.Equal(t, err1, dg.errors[0])

	// Test l'ajout d'une deuxième erreur
	err2 := errors.New("erreur 2")
	dg.addError(err2)
	assert.Len(t, dg.errors, 2)
	assert.Equal(t, err2, dg.errors[1])
}
