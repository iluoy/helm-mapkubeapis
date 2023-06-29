/*
Copyright

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

package main

import (
	"github.com/spf13/pflag"
)

// EnvSettings defined settings
type EnvSettings struct {
	DryRun                      bool
	KubeConfigFile              string
	KubeContext                 string
	MapFile                     string
	Namespaces                  []string
	AllNamespaces               bool
	ReleasesAndNamespaces       []string
	ExceptNamespaces            []string
	ExceptReleasesAndNamespaces []string
}

// New returns default env settings
func New() *EnvSettings {
	envSettings := EnvSettings{}
	return &envSettings
}

// AddBaseFlags binds base flags to the given flagset.
func (s *EnvSettings) AddBaseFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.DryRun, "dry-run", false, "simulate a command")
}

// AddFlags binds flags to the given flagset.
func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	s.AddBaseFlags(fs)
	fs.StringVar(&s.KubeConfigFile, "kubeconfig", "", "path to the kubeconfig file")
	fs.StringVar(&s.KubeContext, "kube-context", s.KubeContext, "name of the kubeconfig context to use")
	fs.StringVar(&s.MapFile, "mapfile", s.MapFile, "path to the API mapping file")
	fs.BoolVarP(&s.AllNamespaces, "all-namespaces", "A", false, "map kube api of all releases across all namespaces")
	fs.StringSliceVar(&s.ReleasesAndNamespaces, "releases-namespaces", []string{}, "multiple releases, for example: --releases-namespaces Release1.NS1 Release2.NS2")
	fs.StringSliceVar(&s.Namespaces, "namespaces", []string{}, "multiple namespaces, for example: --namespaces NS1 NS2")
	fs.StringSliceVar(&s.ExceptNamespaces, "except-namespaces", []string{}, "except multiple namespaces, for example: --except-namespaces NS1 NS2")
	fs.StringSliceVar(&s.ExceptReleasesAndNamespaces, "except-releases-namespaces", []string{}, "except multiple releases namespaces, for example: --except-releases-namespaces Release1.NS1 Release2.NS2")
}
