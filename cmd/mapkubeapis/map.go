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
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/helm/helm-mapkubeapis/pkg/common"
	"github.com/helm/helm-mapkubeapis/pkg/log"
	v3 "github.com/helm/helm-mapkubeapis/pkg/v3"
)

var (
	settings *EnvSettings
)

func newMapCmd(out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "mapkubeapis [flags] ",
		Short:        "Map release deprecated or removed Kubernetes APIs in-place",
		Long:         "Map release deprecated or removed Kubernetes APIs in-place",
		SilenceUsage: true,
		/*
			Args: func(cmd *cobra.Command, args []string) error {
				if len(args) < 0 {
					return errors.New("only one release name may be passed at a time")
					cmd.Help()
					os.Exit(1)
				}
				return nil
			},
		*/

		RunE: runMap,
	}

	flags := cmd.PersistentFlags()
	flags.Parse(args)
	settings = new(EnvSettings)

	// Get the default mapping file
	if ctx := os.Getenv("HELM_PLUGIN_DIR"); ctx != "" {
		settings.MapFile = filepath.Join(ctx, "config", "Map.yaml")
	} else {
		settings.MapFile = filepath.Join("config", "Map.yaml")
	}

	// When run with the Helm plugin framework, Helm plugins are not passed the
	// plugin flags that correspond to Helm global flags e.g. helm mapkubeapis v3map --kube-context ...
	// The flag values are set to corresponding environment variables instead.
	// The flags are passed as expected when run directly using the binary.
	// The below allows to use Helm's --kube-context global flag.
	if ctx := os.Getenv("HELM_KUBECONTEXT"); ctx != "" {
		settings.KubeContext = ctx
	}

	// Note that the plugin's --kubeconfig flag is set by the Helm plugin framework to
	// the KUBECONFIG environment variable instead of being passed into the plugin.

	settings.AddFlags(flags)

	return cmd
}

func runMap(cmd *cobra.Command, args []string) error {

	logger := log.NewLogger()
	kubeConfig := common.KubeConfig{
		Context: settings.KubeContext,
		File:    settings.KubeConfigFile,
	}

	return Map(settings, logger, kubeConfig)
}

// Map checks for Kubernetes deprecated or removed APIs in the manifest of the last deployed release version
// and maps those API versions to supported versions. It then adds a new release version with
// the updated APIs and supersedes the version with the unsupported APIs.
func Map(settings *EnvSettings, logger *logrus.Logger, kubeConfig common.KubeConfig) error {
	if settings.DryRun {
		logger.Info("NOTE: This is in dry-run mode, the following actions will not be executed.")
		logger.Info("Run without --dry-run to take the actions described below:")
	}

	options := common.MapOptions{
		Logger:                      logger,
		DryRun:                      settings.DryRun,
		KubeConfig:                  kubeConfig,
		MapFile:                     settings.MapFile,
		Namespaces:                  settings.Namespaces,
		ExceptNamespaces:            settings.ExceptNamespaces,
		AllNamespaces:               settings.AllNamespaces,
		ReleasesAndNamespaces:       settings.ReleasesAndNamespaces,
		ExceptReleasesAndNamespaces: settings.ExceptReleasesAndNamespaces,
	}

	if err := v3.MapReleaseWithUnSupportedAPIs(options); err != nil {
		return err
	}

	return nil
}
