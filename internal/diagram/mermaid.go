package diagram

import (
	"fmt"
	"strings"

	"github.com/Pitchouneee/koomos/internal/model"
	"github.com/Pitchouneee/koomos/internal/resolver"
)

// func GenerateMermaid(resources []model.Resource) string {
// 	diagram := "```mermaid\ngraph TD\n"
// 	for _, res := range resources {
// 		node := fmt.Sprintf("%s[%s\\n(%s)]", res.Name, res.Kind, res.Namespace)
// 		diagram += "  " + node + "\n"
// 	}
// 	diagram += "```"
// 	return diagram
// }

func sanitizeID(res model.Resource) string {
	id := fmt.Sprintf("%s_%s", res.Name, res.Kind)
	id = strings.ToLower(id)
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, ".", "_")
	id = strings.ReplaceAll(id, "/", "_")
	return id
}

func GenerateMermaid(resources []model.Resource, edges []resolver.Edge) string {
	diagram := "```mermaid\ngraph TD\n"
	seen := make(map[string]bool)

	for _, res := range resources {
		if res.Name == "" || res.Kind == "" {
			continue // skip invalid
		}

		id := sanitizeID(res)
		label := fmt.Sprintf("%s<br>(%s)", res.Name, res.Kind)
		node := fmt.Sprintf("  %s[\"%s\"]\n", id, label)

		if !seen[id] {
			diagram += node
			seen[id] = true
		}
	}

	for _, edge := range edges {
		from := sanitizeID(edge.From)
		to := sanitizeID(edge.To)
		diagram += fmt.Sprintf("  %s --> %s\n", from, to)
	}

	diagram += "```"
	return diagram
}
