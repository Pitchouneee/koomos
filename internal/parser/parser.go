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

			var references []string

			// Handle Deployment mounting ConfigMaps or Secrets
			if kind == "Deployment" {
				// Labels
				if template, ok := spec["template"].(map[string]interface{}); ok {
					if metaTmpl, ok := template["metadata"].(map[string]interface{}); ok {
						if lbls, ok := metaTmpl["labels"].(map[string]interface{}); ok {
							for k, v := range lbls {
								labels[k] = fmt.Sprintf("%v", v)
							}
						}
					}

					// Parse volumes
					if specTmpl, ok := template["spec"].(map[string]interface{}); ok {
						if volumes, ok := specTmpl["volumes"].([]interface{}); ok {
							for _, v := range volumes {
								vol := v.(map[string]interface{})
								if cm, ok := vol["configMap"].(map[string]interface{}); ok {
									if n, ok := cm["name"].(string); ok {
										references = append(references, n)
									}
								}
								if s, ok := vol["secret"].(map[string]interface{}); ok {
									if n, ok := s["secretName"].(string); ok {
										references = append(references, n)
									}
								}
							}
						}

						// Parse env[].valueFrom
						if containers, ok := specTmpl["containers"].([]interface{}); ok {
							for _, c := range containers {
								container := c.(map[string]interface{})
								if envs, ok := container["env"].([]interface{}); ok {
									for _, e := range envs {
										env := e.(map[string]interface{})
										if valFrom, ok := env["valueFrom"].(map[string]interface{}); ok {
											if cmRef, ok := valFrom["configMapKeyRef"].(map[string]interface{}); ok {
												if n, ok := cmRef["name"].(string); ok {
													references = append(references, n)
												}
											}
											if sRef, ok := valFrom["secretKeyRef"].(map[string]interface{}); ok {
												if n, ok := sRef["name"].(string); ok {
													references = append(references, n)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}

			if kind == "Ingress" {
				if rules, ok := spec["rules"].([]interface{}); ok {
					for _, r := range rules {
						if rule, ok := r.(map[string]interface{}); ok {
							if http, ok := rule["http"].(map[string]interface{}); ok {
								if paths, ok := http["paths"].([]interface{}); ok {
									for _, p := range paths {
										if path, ok := p.(map[string]interface{}); ok {
											if backend, ok := path["backend"].(map[string]interface{}); ok {
												if service, ok := backend["service"].(map[string]interface{}); ok {
													if svcName, ok := service["name"].(string); ok {
														references = append(references, svcName)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}

			resources = append(resources, model.Resource{
				Kind:                kind,
				Name:                name,
				Namespace:           namespace,
				Path:                path,
				Labels:              labels,
				Selector:            selector,
				ReferencedResources: references,
			})
		}

		return nil
	})

	return resources, err
}
