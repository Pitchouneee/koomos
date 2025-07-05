package resolver

import (
	"github.com/Pitchouneee/koomos/internal/model"
)

type Edge struct {
	From model.Resource
	To   model.Resource
	Type string
}

// Match checks if all keys in A exist with same value in B
func selectorMatches(selector, labels map[string]string) bool {
	for k, v := range selector {
		if labels[k] != v {
			return false
		}
	}
	return true
}

func ResolveRelations(resources []model.Resource) []Edge {
	var edges []Edge

	for _, from := range resources {
		switch from.Kind {
		case "Service":
			for _, to := range resources {
				if to.Kind == "Deployment" &&
					from.Namespace == to.Namespace &&
					selectorMatches(from.Selector, to.Labels) {
					edges = append(edges, Edge{
						From: from,
						To:   to,
						Type: "routes-to",
					})
				}
			}

		case "Application":
			for _, to := range resources {
				if to.Kind == "Namespace" &&
					from.Namespace == to.Name {
					edges = append(edges, Edge{
						From: from,
						To:   to,
						Type: "targets",
					})
				}
			}
		}
	}

	return edges
}
