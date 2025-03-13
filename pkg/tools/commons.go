package tools

import (
	"fmt"
	"strings"
)

// compareVersions compare deux chaînes de version
// Retourne: -1 si v1 < v2, 0 si v1 == v2, 1 si v1 > v2
func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(strings.TrimPrefix(v1, "v"), ".")
	v2Parts := strings.Split(strings.TrimPrefix(v2, "v"), ".")

	for i := 0; i < len(v1Parts) && i < len(v2Parts); i++ {
		v1Num := parseInt(v1Parts[i])
		v2Num := parseInt(v2Parts[i])

		if v1Num < v2Num {
			return -1
		}
		if v1Num > v2Num {
			return 1
		}
	}

	if len(v1Parts) < len(v2Parts) {
		return -1
	}
	if len(v1Parts) > len(v2Parts) {
		return 1
	}

	return 0
}

func parseInt(s string) int {
	// Supprimer tout préfixe/suffixe non numérique
	numStr := strings.TrimFunc(s, func(r rune) bool {
		return r < '0' || r > '9'
	})

	// Convertir en int, par défaut à 0 si vide ou invalide
	num := 0
	fmt.Sscanf(numStr, "%d", &num)
	return num
}
