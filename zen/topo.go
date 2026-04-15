package zen

import "fmt"

// topoSort performs Kahn's algorithm on registered components.
// Returns components in dependency order (leaves first).
func topoSort(components map[string]Component) ([]Component, error) {
	inDeg := make(map[string]int, len(components))
	downstream := make(map[string][]string)

	for name := range components {
		inDeg[name] = 0
	}

	for name, comp := range components {
		for _, dep := range comp.Depends() {
			if _, ok := components[dep]; !ok {
				return nil, fmt.Errorf("component %q depends on unregistered component %q", name, dep)
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

	sorted := make([]Component, 0, len(components))
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		sorted = append(sorted, components[cur])
		for _, next := range downstream[cur] {
			inDeg[next]--
			if inDeg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if len(sorted) != len(components) {
		return nil, fmt.Errorf("circular dependency detected among components")
	}
	return sorted, nil
}

// topoLevels groups components by dependency depth.
// Components within the same level have no inter-dependencies and can be initialized/started in parallel.
// Level 0 contains components with no dependencies; level N depends only on components in levels < N.
func topoLevels(components map[string]Component) ([][]Component, error) {
	inDeg := make(map[string]int, len(components))
	downstream := make(map[string][]string, len(components))

	for name := range components {
		inDeg[name] = 0
	}

	for name, comp := range components {
		for _, dep := range comp.Depends() {
			if _, ok := components[dep]; !ok {
				return nil, fmt.Errorf("component %q depends on unregistered component %q", name, dep)
			}
			inDeg[name]++
			downstream[dep] = append(downstream[dep], name)
		}
	}

	// Collect roots (in-degree == 0)
	var current []string
	for name, d := range inDeg {
		if d == 0 {
			current = append(current, name)
		}
	}

	var levels [][]Component
	total := 0
	for len(current) > 0 {
		level := make([]Component, len(current))
		for i, name := range current {
			level[i] = components[name]
		}
		levels = append(levels, level)
		total += len(current)

		var next []string
		for _, name := range current {
			for _, child := range downstream[name] {
				inDeg[child]--
				if inDeg[child] == 0 {
					next = append(next, child)
				}
			}
		}
		current = next
	}

	if total != len(components) {
		return nil, fmt.Errorf("circular dependency detected among components")
	}
	return levels, nil
}

func componentNames(comps []Component) []string {
	out := make([]string, len(comps))
	for i, c := range comps {
		out[i] = c.Name()
	}
	return out
}
