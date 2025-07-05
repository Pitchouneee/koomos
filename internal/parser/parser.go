package parser

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Pitchouneee/koomos/internal/model"
	"gopkg.in/yaml.v3"
)

func ParseDirectory(root string) ([]model.Resource, error) {
	var resources []model.Resource

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		if d.IsDir() || (ext != ".yaml" && ext != ".yml") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		dec := yaml.NewDecoder(bytes.NewReader(content))
		for {
			var raw map[string]interface{}
			if err := dec.Decode(&raw); err != nil {
				break
			}

			kind, _ := raw["kind"].(string)
			meta, _ := raw["metadata"].(map[string]interface{})
			name, _ := meta["name"].(string)
			namespace, _ := meta["namespace"].(string)
			spec, _ := raw["spec"].(map[string]interface{})

			labels := map[string]string{}
			if metaLabels, ok := meta["labels"].(map[string]interface{}); ok {
				for k, v := range metaLabels {
					labels[k] = fmt.Sprintf("%v", v)
				}
			}

			selector := map[string]string{}
			if kind == "Service" {
				if sel, ok := spec["selector"].(map[string]interface{}); ok {
					for k, v := range sel {
						selector[k] = fmt.Sprintf("%v", v)
					}
				}
			}

			// For Deployment: labels are inside spec.template.metadata.labels
			if kind == "Deployment" {
				if template, ok := spec["template"].(map[string]interface{}); ok {
					if metaTmpl, ok := template["metadata"].(map[string]interface{}); ok {
						if lbls, ok := metaTmpl["labels"].(map[string]interface{}); ok {
							for k, v := range lbls {
								labels[k] = fmt.Sprintf("%v", v)
							}
						}
					}
				}
			}

			resources = append(resources, model.Resource{
				Kind:      kind,
				Name:      name,
				Namespace: namespace,
				Path:      path,
				Labels:    labels,
				Selector:  selector,
			})
		}

		return nil
	})

	return resources, err
}
