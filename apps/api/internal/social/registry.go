package social

import "fmt"

// Registry maps platform names to their Publisher implementations.
type Registry struct {
	publishers map[string]Publisher
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{publishers: make(map[string]Publisher)}
}

// Register adds a Publisher to the registry.
// It panics if a publisher for that platform is registered twice.
func (r *Registry) Register(p Publisher) {
	key := p.Platform()
	if _, exists := r.publishers[key]; exists {
		panic(fmt.Sprintf("social: publisher for platform %q already registered", key))
	}
	r.publishers[key] = p
}

// Get returns the Publisher for the given platform.
// Returns (nil, error) when no publisher is registered for that platform.
func (r *Registry) Get(platform string) (Publisher, error) {
	p, ok := r.publishers[platform]
	if !ok {
		return nil, fmt.Errorf("social: no publisher registered for platform %q", platform)
	}
	return p, nil
}

// Platforms returns the list of registered platform names.
func (r *Registry) Platforms() []string {
	keys := make([]string, 0, len(r.publishers))
	for k := range r.publishers {
		keys = append(keys, k)
	}
	return keys
}
