/*
Copyright 2021 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flaker

import (
	"context"
	"embed"

	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/manifest"
)

//go:embed *.yaml
var yaml embed.FS

// Install
func Install(name, sink string) feature.StepFn {
	cfg := map[string]interface{}{
		"name": name,
		"sink": sink,
	}

	return func(ctx context.Context, t feature.T) {
		if err := registerImage(ctx); err != nil {
			t.Fatal(err)
		}
		manifest.PodSecurityCfgFn(ctx, t)(cfg)
		if _, err := manifest.InstallYamlFS(ctx, yaml, cfg); err != nil {
			t.Fatal(err)
		}
	}
}

// AsRef returns a KRef for a Service without namespace.
func AsRef(name string) *duckv1.KReference {
	return &duckv1.KReference{
		Kind:       "Service",
		APIVersion: "v1",
		Name:       name,
	}
}

func registerImage(ctx context.Context) error {
	im := manifest.ImagesFromFS(ctx, yaml)
	reg := environment.RegisterPackage(im...)
	_, err := reg(ctx, environment.FromContext(ctx))
	return err
}
