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

func contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func ResolveRelations(resources []model.Resource) []Edge {
	var edges []Edge

	for _, from := range resources {
		switch from.Kind {
		case "Service":
			// Service → Deployment (label selector)
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
			// Application → Namespace (target)
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

		case "Deployment":
			// Deployment → ConfigMap / Secret (mounted or referenced)
			for _, to := range resources {
				if (to.Kind == "ConfigMap" || to.Kind == "Secret") &&
					from.Namespace == to.Namespace &&
					contains(from.ReferencedResources, to.Name) {
					edges = append(edges, Edge{
						From: from,
						To:   to,
						Type: "uses",
					})
				}
			}

		case "Ingress":
			// Ingress → Service (ref by backend service name)
			for _, to := range resources {
				if to.Kind == "Service" &&
					from.Namespace == to.Namespace &&
					contains(from.ReferencedResources, to.Name) {
					edges = append(edges, Edge{
						From: from,
						To:   to,
						Type: "routes-to",
					})
				}
			}
		}
	}

	return edges
}
