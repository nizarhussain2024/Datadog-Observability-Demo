package main

import (
	"context"
)

type Tags map[string]string

func NewTags() Tags {
	return make(Tags)
}

func (t Tags) WithEnvironment(env string) Tags {
	t["environment"] = env
	return t
}

func (t Tags) WithService(service string) Tags {
	t["service"] = service
	return t
}

func (t Tags) WithVersion(version string) Tags {
	t["version"] = version
	return t
}

func (t Tags) WithCustom(key, value string) Tags {
	t[key] = value
	return t
}

func (t Tags) ToContext(ctx context.Context) context.Context {
	// In production, attach tags to context for Datadog
	return ctx
}



