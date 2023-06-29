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

package v3

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"

	common "github.com/helm/helm-mapkubeapis/pkg/common"
)

type NamespaceRelease struct {
	Namespace string
	Release   string
}

type NamespaceReleases struct {
	AllNamespaces    bool
	Namespaces       []string
	ExceptNamespaces []string
	Releases         []NamespaceRelease
	ExceptReleases   []NamespaceRelease
}

func NewNamespaceReleases() *NamespaceReleases {
	return &NamespaceReleases{}
}

func getReleases(mapOptions common.MapOptions) *NamespaceReleases {
	releases := NewNamespaceReleases()
	if !mapOptions.AllNamespaces {
		if len(mapOptions.ReleasesAndNamespaces) > 0 {
			for _, v := range mapOptions.ReleasesAndNamespaces {
				releaseName := strings.Split(v, ".")[0]
				namespace := strings.Split(v, ".")[1]
				if len(mapOptions.Namespaces) > 0 {
					if common.Contains(releases.Namespaces, namespace) >= 0 {
						continue
					}
					if common.Contains(mapOptions.Namespaces, namespace) >= 0 {
						releases.Namespaces = append(releases.Namespaces, namespace)
						continue
					}
					releases.Releases = append(releases.Releases, NamespaceRelease{
						Release:   releaseName,
						Namespace: namespace,
					})
				} else {
					releases.Releases = append(releases.Releases, NamespaceRelease{
						Release:   releaseName,
						Namespace: namespace,
					})
				}
			}
		}

		if len(mapOptions.Namespaces) > 0 {
			for _, v := range mapOptions.Namespaces {
				if common.Contains(releases.Namespaces, v) >= 0 {
					continue
				}
				releases.Namespaces = append(releases.Namespaces, v)
			}
		}

		if len(mapOptions.ExceptReleasesAndNamespaces) > 0 {
			for _, v := range mapOptions.ExceptReleasesAndNamespaces {
				releaseName := strings.Split(v, ".")[0]
				namespace := strings.Split(v, ".")[1]
				if len(mapOptions.ExceptNamespaces) > 0 {
					if common.Contains(releases.ExceptNamespaces, namespace) >= 0 {
						continue
					}
					if common.Contains(mapOptions.ExceptNamespaces, namespace) >= 0 {
						releases.ExceptNamespaces = append(releases.ExceptNamespaces, namespace)
						continue
					}
					releases.ExceptReleases = append(releases.ExceptReleases, NamespaceRelease{
						Release:   releaseName,
						Namespace: namespace,
					})
				} else {
					releases.ExceptReleases = append(releases.ExceptReleases, NamespaceRelease{
						Release:   releaseName,
						Namespace: namespace,
					})
				}
			}
		}
		if len(mapOptions.ExceptNamespaces) > 0 {
			for _, v := range mapOptions.ExceptNamespaces {
				if common.Contains(releases.ExceptNamespaces, v) >= 0 {
					continue
				}
				releases.ExceptNamespaces = append(releases.ExceptNamespaces, v)
			}
		}
	} else {
		releases.AllNamespaces = mapOptions.AllNamespaces
		if len(mapOptions.ExceptReleasesAndNamespaces) > 0 {
			for _, v := range mapOptions.ExceptReleasesAndNamespaces {
				releaseName := strings.Split(v, ".")[0]
				namespace := strings.Split(v, ".")[1]
				if len(mapOptions.ExceptNamespaces) > 0 {
					if common.Contains(releases.ExceptNamespaces, namespace) >= 0 {
						continue
					}
					if common.Contains(mapOptions.ExceptNamespaces, namespace) >= 0 {
						releases.ExceptNamespaces = append(releases.ExceptNamespaces, namespace)
						continue
					}
					releases.ExceptReleases = append(releases.ExceptReleases, NamespaceRelease{
						Release:   releaseName,
						Namespace: namespace,
					})
				} else {
					releases.ExceptReleases = append(releases.ExceptReleases, NamespaceRelease{
						Release:   releaseName,
						Namespace: namespace,
					})
				}
			}
		}
		if len(mapOptions.ExceptNamespaces) > 0 {
			for _, v := range mapOptions.ExceptNamespaces {
				if common.Contains(releases.ExceptNamespaces, v) >= 0 {
					continue
				}
				releases.ExceptNamespaces = append(releases.ExceptNamespaces, v)
			}
		}
	}

	return releases
}

func filterReleases(results []*release.Release, releases *NamespaceReleases) (filtered_results []*release.Release) {
FILTER:
	for _, res := range results {
		releaseName := res.Name
		namespace := res.Namespace

		if len(releases.ExceptNamespaces) > 0 {
			if common.Contains(releases.ExceptNamespaces, namespace) >= 0 {
				continue FILTER
			}
		}

		for _, v := range releases.ExceptReleases {
			if v.Release == releaseName && v.Namespace == namespace {
				continue FILTER
			}
		}

		if releases.AllNamespaces {
			filtered_results = append(filtered_results, res)
			continue FILTER
		}

		if len(releases.Namespaces) > 0 {
			if common.Contains(releases.Namespaces, namespace) >= 0 {
				filtered_results = append(filtered_results, res)
			}
		}

		for _, v := range releases.Releases {
			if v.Release == releaseName && v.Namespace == namespace {
				filtered_results = append(filtered_results, res)
			}
		}

	}

	return
}

// MapReleaseWithUnSupportedAPIs checks the latest release version for any deprecated or removed APIs in its metadata
// If it finds any, it will create a new release version with the APIs mapped to the supported versions
func MapReleaseWithUnSupportedAPIs(mapOptions common.MapOptions) error {
	logger := mapOptions.Logger
	cfg, err := GetActionConfig(mapOptions.KubeConfig)
	if err != nil {
		return errors.Wrap(err, "failed to get Helm action configuration")
	}

	client := action.NewList(cfg)
	results, err := client.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to list all releases")
	}
	releases := getReleases(mapOptions)
	results = filterReleases(results, releases)

	for _, res := range results {
		releaseName := res.Name
		namespace := res.Namespace
		log_with_fields := logger.WithFields(logrus.Fields{
			"namespace":   namespace,
			"releaseName": releaseName,
		})

		fmt.Println()
		log_with_fields.Infof("Check release %s.%s for deprecated or removed APIs...\n", releaseName, namespace)
		origManifest := res.Manifest

		modifiedManifest, err := common.ReplaceManifestUnSupportedAPIs(origManifest, mapOptions.MapFile, mapOptions.KubeConfig, logger)
		if err != nil {
			continue
			//return err
		}
		log_with_fields.Infof("Finished checking release %s.%s for deprecated or removed APIs.\n", releaseName, namespace)
		if modifiedManifest == origManifest {
			log_with_fields.Infof("Release %s.%s has no deprecated or removed APIs.\n", releaseName, namespace)
			continue
			//return nil
		}

		if mapOptions.DryRun {
			log_with_fields.Infof("Deprecated or removed APIs exist, for release: %s.%s.\n", releaseName, namespace)
		} else {
			log_with_fields.Infof("Deprecated or removed APIs exist, updating release: %s.%s.\n", releaseName, namespace)
			if err := updateRelease(res, modifiedManifest, cfg, logger); err != nil {
				continue
				//return errors.Wrapf(err, "failed to update release '%s'.'%s'", releaseName, namespace)
			}
			log_with_fields.Infof("Release '%s'.'%s' with deprecated or removed APIs updated successfully to new version.\n", releaseName, namespace)
		}
	}

	return nil
}

func updateRelease(origRelease *release.Release, modifiedManifest string, cfg *action.Configuration, logger *logrus.Logger) error {
	// Update current release version to be superseded
	logger.Infof("Set status of release version '%s' to 'superseded'.\n", getReleaseVersionName(origRelease))
	origRelease.Info.Status = release.StatusSuperseded
	if err := cfg.Releases.Update(origRelease); err != nil {
		return errors.Wrapf(err, "failed to update release version '%s'", getReleaseVersionName(origRelease))
	}
	logger.Infof("Release version '%s' updated successfully.\n", getReleaseVersionName(origRelease))

	// Using a shallow copy of current release version to update the object with the modification
	// and then store this new version
	var newRelease = origRelease
	newRelease.Manifest = modifiedManifest
	newRelease.Info.Description = common.UpgradeDescription
	newRelease.Info.LastDeployed = cfg.Now()
	newRelease.Version = origRelease.Version + 1
	newRelease.Info.Status = release.StatusDeployed
	logger.Infof("Add release version '%s' with updated supported APIs.\n", getReleaseVersionName(origRelease))
	if err := cfg.Releases.Create(newRelease); err != nil {
		return errors.Wrapf(err, "failed to create new release version '%s'", getReleaseVersionName(origRelease))
	}
	logger.Infof("Release version '%s' added successfully.\n", getReleaseVersionName(origRelease))
	return nil
}

func getLatestRelease(releaseName string, cfg *action.Configuration) (*release.Release, error) {
	return cfg.Releases.Last(releaseName)
}

func getReleaseVersionName(rel *release.Release) string {
	return fmt.Sprintf("%s.v%d", rel.Name, rel.Version)
}
