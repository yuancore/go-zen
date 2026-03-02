package zen

import "fmt"

// topoSort performs Kahn's algorithm on registered modules.
func topoSort(modules map[string]Module) ([]Module, error) {
	inDeg := make(map[string]int, len(modules))
	downstream := make(map[string][]string)
	for name := range modules {
		inDeg[name] = 0
	}
	for name, mod := range modules {
		for _, dep := range mod.Depends() {
			if _, ok := modules[dep]; !ok {
				return nil, fmt.Errorf("module %q depends on unregistered module %q", name, dep)
			}
			inDeg[name]++
			downstream[dep] = append(downstream[dep], name)
		}
	}
	var queue []string
	for name, d := range inDeg {
		if d == 0 {
			queue = append(queue, name)
		}
	}
	sorted := make([]Module, 0, len(modules))
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		sorted = append(sorted, modules[cur])
		for _, next := range downstream[cur] {
			inDeg[next]--
			if inDeg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	if len(sorted) != len(modules) {
		return nil, fmt.Errorf("circular dependency detected among modules")
	}
	return sorted, nil
}
