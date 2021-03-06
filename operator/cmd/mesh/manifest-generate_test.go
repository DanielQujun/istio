// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mesh

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	klabels "k8s.io/apimachinery/pkg/labels"

	"istio.io/istio/operator/pkg/compare"
	"istio.io/istio/operator/pkg/helm"
	"istio.io/istio/operator/pkg/object"
	"istio.io/istio/operator/pkg/tpath"
	"istio.io/istio/operator/pkg/util"
	"istio.io/istio/operator/pkg/util/httpserver"
	"istio.io/istio/operator/pkg/util/tgz"
	"istio.io/istio/pkg/test"
	"istio.io/pkg/version"
)

const (
	istioTestVersion = "istio-1.5.0"
	testTGZFilename  = istioTestVersion + "-linux.tar.gz"
	testDataSubdir   = "cmd/mesh/testdata/manifest-generate"
)

// chartSourceType defines where charts used in the test come from.
type chartSourceType int

const (
	// Snapshot charts are in testdata/manifest-generate/data-snapshot
	snapshotCharts chartSourceType = iota
	// Compiled in charts come from assets.gen.go
	compiledInCharts
	// Live charts come from manifests/
	liveCharts
)

type testGroup []struct {
	desc string
	// Small changes to the input profile produce large changes to the golden output
	// files. This makes it difficult to spot meaningful changes in pull requests.
	// By default we hide these changes to make developers life's a bit easier. However,
	// it is still useful to sometimes override this behavior and show the full diff.
	// When this flag is true, use an alternative file suffix that is not hidden by
	// default github in pull requests.
	showOutputFileInPullRequest bool
	flags                       string
	noInput                     bool
	outputDir                   string
	diffSelect                  string
	diffIgnore                  string
	chartSource                 chartSourceType
}

// TestMain is required to create a local release package in /tmp from manifests and operator/data in the format that
// istioctl expects.
func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	operatorRootDir = filepath.Join(wd, "../..")
	manifestsDir = filepath.Join(operatorRootDir, "manifests")
	liveReleaseDir, err = createLocalReleaseCharts()
	defer os.RemoveAll(liveReleaseDir)
	if err != nil {
		panic(err)
	}
	liveInstallPackageDir = filepath.Join(liveReleaseDir, istioTestVersion, helm.OperatorSubdirFilePath)
	snapshotInstallPackageDir = filepath.Join(operatorRootDir, testDataSubdir, "data-snapshot")

	flag.Parse()
	code := m.Run()
	os.Exit(code)
}

func TestManifestGenerateFlags(t *testing.T) {
	flagOutputDir := createTempDirOrFail(t, "flag-output")
	flagOutputValuesDir := createTempDirOrFail(t, "flag-output-values")
	runTestGroup(t, testGroup{
		{
			desc: "all_off",
		},
		{
			desc:                        "all_on",
			diffIgnore:                  "ConfigMap:*:istio",
			showOutputFileInPullRequest: true,
		},
		{
			desc:       "prometheus",
			diffIgnore: "ConfigMap:*:istio",
		},
		{
			desc:       "gateways",
			diffIgnore: "ConfigMap:*:istio",
		},
		{
			desc:       "gateways_override_default",
			diffIgnore: "ConfigMap:*:istio",
		},
		{
			desc:       "component_hub_tag",
			diffSelect: "Deployment:*:*",
		},
		{
			desc:       "flag_set_values",
			diffSelect: "Deployment:*:istio-ingressgateway,ConfigMap:*:istio-sidecar-injector",
			flags:      "-s values.global.proxy.image=myproxy --set values.global.proxy.includeIPRanges=172.30.0.0/16,172.21.0.0/16",
			noInput:    true,
		},
		{
			desc:       "flag_values_enable_egressgateway",
			diffSelect: "Service:*:istio-egressgateway",
			flags:      "--set values.gateways.istio-egressgateway.enabled=true",
			noInput:    true,
		},
		{
			desc:       "flag_override_values",
			diffSelect: "Deployment:*:istiod",
			flags:      "-s tag=my-tag",
		},
		{
			desc:       "flag_output",
			flags:      "-o " + flagOutputDir,
			diffSelect: "Deployment:*:istiod",
			outputDir:  flagOutputDir,
		},
		{
			desc:       "flag_output_set_values",
			diffSelect: "Deployment:*:istio-ingressgateway",
			flags:      "-s values.global.proxy.image=mynewproxy -o " + flagOutputValuesDir,
			outputDir:  flagOutputValuesDir,
			noInput:    true,
		},
		{
			desc:       "flag_force",
			diffSelect: "no:resources:selected",
			flags:      "--force",
		},
		{
			desc:       "flag_output_set_profile",
			diffIgnore: "ConfigMap:*:istio",
			flags:      "-s profile=minimal",
			noInput:    true,
		},
	})
	removeDirOrFail(t, flagOutputDir)
	removeDirOrFail(t, flagOutputValuesDir)
}

func TestManifestGeneratePilot(t *testing.T) {
	runTestGroup(t, testGroup{
		{
			desc:       "pilot_default",
			diffIgnore: "CustomResourceDefinition:*:*,ConfigMap:*:istio",
		},
		{
			desc:       "pilot_k8s_settings",
			diffSelect: "Deployment:*:istiod,HorizontalPodAutoscaler:*:istiod",
		},
		{
			desc:       "pilot_override_values",
			diffSelect: "Deployment:*:istiod,HorizontalPodAutoscaler:*:istiod",
		},
		{
			desc:       "pilot_override_kubernetes",
			diffSelect: "Deployment:*:istiod, Service:*:istiod",
		},
		// TODO https://github.com/istio/istio/issues/22347 this is broken for overriding things to default value
		// This can be seen from REGISTRY_ONLY not applying
		{
			desc:       "pilot_merge_meshconfig",
			diffSelect: "ConfigMap:*:istio$",
		},
	})
}

func TestManifestGenerateTelemetry(t *testing.T) {
	runTestGroup(t, testGroup{
		{
			desc: "all_off",
		},
		{
			desc:       "telemetry_default",
			diffIgnore: "",
		},
		{
			desc:       "telemetry_k8s_settings",
			diffSelect: "Deployment:*:istio-telemetry, HorizontalPodAutoscaler:*:istio-telemetry",
		},
		{
			desc:       "telemetry_override_values",
			diffSelect: "handler:*:prometheus",
		},
		{
			desc:       "telemetry_override_kubernetes",
			diffSelect: "Deployment:*:istio-telemetry, handler:*:prometheus",
		},
	})
}

func TestManifestGenerateGateway(t *testing.T) {
	runTestGroup(t, testGroup{
		{
			desc:       "ingressgateway_k8s_settings",
			diffSelect: "Deployment:*:istio-ingressgateway, Service:*:istio-ingressgateway",
		},
	})
}

func TestManifestGenerateAddonK8SOverride(t *testing.T) {
	runTestGroup(t, testGroup{
		{
			desc:       "addon_k8s_override",
			diffSelect: "Service:*:prometheus, Deployment:*:prometheus, Service:*:kiali",
		},
	})
}

// TestManifestGenerateHelmValues tests whether enabling components through the values passthrough interface works as
// expected i.e. without requiring enablement also in IstioOperator API.
func TestManifestGenerateHelmValues(t *testing.T) {
	runTestGroup(t, testGroup{
		{
			desc: "helm_values_enablement",
			diffSelect: "Deployment:*:istio-egressgateway, Service:*:istio-egressgateway," +
				" Deployment:*:kiali, Service:*:kiali, Deployment:*:prometheus, Service:*:prometheus",
		},
	})
}

func TestManifestGenerateOrdered(t *testing.T) {
	testDataDir = filepath.Join(operatorRootDir, "cmd/mesh/testdata/manifest-generate")
	// Since this is testing the special case of stable YAML output order, it
	// does not use the established test group pattern
	inPath := filepath.Join(testDataDir, "input/all_on.yaml")
	got1, err := runManifestGenerate([]string{inPath}, "", snapshotCharts)
	if err != nil {
		t.Fatal(err)
	}
	got2, err := runManifestGenerate([]string{inPath}, "", snapshotCharts)
	if err != nil {
		t.Fatal(err)
	}

	if got1 != got2 {
		fmt.Printf("%s", util.YAMLDiff(got1, got2))
		t.Errorf("stable_manifest: Manifest generation is not producing stable text output.")
	}
}

func TestMultiICPSFiles(t *testing.T) {
	testDataDir = filepath.Join(operatorRootDir, "cmd/mesh/testdata/manifest-generate")
	inPathBase := filepath.Join(testDataDir, "input/all_off.yaml")
	inPathOverride := filepath.Join(testDataDir, "input/telemetry_override_only.yaml")
	got, err := runManifestGenerate([]string{inPathBase, inPathOverride}, "", snapshotCharts)
	if err != nil {
		t.Fatal(err)
	}
	outPath := filepath.Join(testDataDir, "output/telemetry_override_values"+goldenFileSuffixHideChangesInReview)

	want, err := readFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	diffSelect := "handler:*:prometheus"
	got, err = compare.SelectAndIgnoreFromOutput(got, diffSelect, "")
	if err != nil {
		t.Errorf("error selecting from output manifest: %v", err)
	}
	diff := compare.YAMLCmp(got, want)
	if diff != "" {
		t.Errorf("`manifest generate` diff = %s", diff)
	}
}

func TestBareSpec(t *testing.T) {
	testDataDir = filepath.Join(operatorRootDir, "cmd/mesh/testdata/manifest-generate")
	inPathBase := filepath.Join(testDataDir, "input/bare_spec.yaml")
	_, err := runManifestGenerate([]string{inPathBase}, "", liveCharts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstallPackagePath(t *testing.T) {
	testDataDir = filepath.Join(operatorRootDir, "cmd/mesh/testdata/manifest-generate")
	serverDir, err := ioutil.TempDir(os.TempDir(), "istio-test-server-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(serverDir)
	if err := tgz.Create(liveReleaseDir, filepath.Join(serverDir, testTGZFilename)); err != nil {
		t.Fatal(err)
	}
	srv := httpserver.NewServer(serverDir)
	runTestGroup(t, testGroup{
		{
			// Use some arbitrary small test input (pilot only) since we are testing the local filesystem code here, not
			// manifest generation.
			desc:       "install_package_path",
			diffSelect: "Deployment:*:istiod",
			flags:      "--set installPackagePath=" + liveInstallPackageDir,
		},
		{
			// Specify both charts and profile from local filesystem.
			desc:       "install_package_path",
			diffSelect: "Deployment:*:istiod",
			flags:      fmt.Sprintf("--set installPackagePath=%s --set profile=%s/profiles/default.yaml", liveInstallPackageDir, liveInstallPackageDir),
		},
		{
			// --force is needed for version mismatch.
			desc:       "install_package_path",
			diffSelect: "Deployment:*:istiod",
			flags:      "--force --set installPackagePath=" + srv.URL() + "/" + testTGZFilename,
		},
	})

}

// This test enforces that objects that reference other objects do so properly, such as Service selecting deployment
func TestConfigSelectors(t *testing.T) {
	got, err := runManifestGenerate([]string{}, "", liveCharts)
	if err != nil {
		t.Fatal(err)
	}
	objs, err := object.ParseK8sObjectsFromYAMLManifest(got)
	if err != nil {
		t.Fatal(err)
	}
	gotRev, e := runManifestGenerate([]string{}, "--set revision=canary", liveCharts)
	if e != nil {
		t.Fatal(e)
	}
	objsRev, err := object.ParseK8sObjectsFromYAMLManifest(gotRev)
	if err != nil {
		t.Fatal(err)
	}

	// First we fetch all the objects for our default install
	name := "istiod"
	deployment := mustFindObject(t, objs, name, "Deployment")
	service := mustFindObject(t, objs, name, "Service")
	pdb := mustFindObject(t, objs, name, "PodDisruptionBudget")
	hpa := mustFindObject(t, objs, name, "HorizontalPodAutoscaler")
	podLabels := mustGetLabels(t, deployment, "spec.template.metadata.labels")
	// Check all selectors align
	mustSelect(t, mustGetLabels(t, pdb, "spec.selector.matchLabels"), podLabels)
	mustSelect(t, mustGetLabels(t, service, "spec.selector"), podLabels)
	mustSelect(t, mustGetLabels(t, deployment, "spec.selector.matchLabels"), podLabels)
	if hpaName := mustGetPath(t, hpa, "spec.scaleTargetRef.name"); name != hpaName {
		t.Fatalf("HPA does not match deployment: %v != %v", name, hpaName)
	}

	// Next we fetch all the objects for a revision install
	nameRev := "istiod-canary"
	deploymentRev := mustFindObject(t, objsRev, nameRev, "Deployment")
	serviceRev := mustFindObject(t, objsRev, nameRev, "Service")
	pdbRev := mustFindObject(t, objsRev, nameRev, "PodDisruptionBudget")
	hpaRev := mustFindObject(t, objsRev, nameRev, "HorizontalPodAutoscaler")
	podLabelsRev := mustGetLabels(t, deploymentRev, "spec.template.metadata.labels")
	// Check all selectors align for revision
	mustSelect(t, mustGetLabels(t, pdbRev, "spec.selector.matchLabels"), podLabelsRev)
	mustSelect(t, mustGetLabels(t, serviceRev, "spec.selector"), podLabelsRev)
	mustSelect(t, mustGetLabels(t, deploymentRev, "spec.selector.matchLabels"), podLabelsRev)
	if hpaName := mustGetPath(t, hpaRev, "spec.scaleTargetRef.name"); nameRev != hpaName {
		t.Fatalf("HPA does not match deployment: %v != %v", nameRev, hpaName)
	}

	// Make sure default and revisions do not cross
	mustNotSelect(t, mustGetLabels(t, serviceRev, "spec.selector"), podLabels)
	mustNotSelect(t, mustGetLabels(t, service, "spec.selector"), podLabelsRev)
	mustNotSelect(t, mustGetLabels(t, pdbRev, "spec.selector.matchLabels"), podLabels)
	mustNotSelect(t, mustGetLabels(t, pdb, "spec.selector.matchLabels"), podLabelsRev)

	// Check selection of previous versions . This only matters for in place upgrade (non revision)
	podLabels15 := map[string]string{
		"app":   "istiod",
		"istio": "pilot",
	}
	mustSelect(t, mustGetLabels(t, service, "spec.selector"), podLabels15)
	mustNotSelect(t, mustGetLabels(t, serviceRev, "spec.selector"), podLabels15)
	mustSelect(t, mustGetLabels(t, pdb, "spec.selector.matchLabels"), podLabels15)
	mustNotSelect(t, mustGetLabels(t, pdbRev, "spec.selector.matchLabels"), podLabels15)

	// Check we aren't changing immutable fields. This only matters for in place upgrade (non revision)
	// This one is not a selector, it must be an exact match
	deploymentSelector15 := map[string]string{
		"istio": "pilot",
	}
	if sel := mustGetLabels(t, deployment, "spec.selector.matchLabels"); !reflect.DeepEqual(deploymentSelector15, sel) {
		t.Fatalf("Depployment selectors are immutable, but changed since 1.5. Was %v, now is %v", deploymentSelector15, sel)
	}
}

func mustSelect(t test.Failer, selector map[string]string, labels map[string]string) {
	t.Helper()
	kselector := klabels.Set(selector).AsSelectorPreValidated()
	if !kselector.Matches(klabels.Set(labels)) {
		t.Fatalf("%v does not select %v", selector, labels)
	}
}

func mustNotSelect(t test.Failer, selector map[string]string, labels map[string]string) {
	t.Helper()
	kselector := klabels.Set(selector).AsSelectorPreValidated()
	if kselector.Matches(klabels.Set(labels)) {
		t.Fatalf("%v selects %v when it should not", selector, labels)
	}
}

func mustGetLabels(t test.Failer, obj object.K8sObject, path string) map[string]string {
	t.Helper()
	got := mustGetPath(t, obj, path)
	conv, ok := got.(map[string]interface{})
	if !ok {
		t.Fatalf("could not convert %v", got)
	}
	ret := map[string]string{}
	for k, v := range conv {
		sv, ok := v.(string)
		if !ok {
			t.Fatalf("could not convert to string %v", v)
		}
		ret[k] = sv
	}
	return ret
}

func mustGetPath(t test.Failer, obj object.K8sObject, path string) interface{} {
	t.Helper()
	got, f, err := tpath.GetFromTreePath(obj.UnstructuredObject().UnstructuredContent(), util.PathFromString(path))
	if err != nil {
		t.Fatal(err)
	}
	if !f {
		t.Fatalf("couldn't find path %v", path)
	}
	return got
}

func mustFindObject(t test.Failer, objs object.K8sObjects, name, kind string) object.K8sObject {
	t.Helper()
	o := findObject(objs, name, kind)
	if o == nil {
		t.Fatalf("expected %v/%v", name, kind)
		return object.K8sObject{}
	}
	return *o
}

func findObject(objs object.K8sObjects, name, kind string) *object.K8sObject {
	for _, o := range objs {
		if o.Kind == kind && o.Name == name {
			return o
		}
	}
	return nil
}

// TestLDFlags checks whether building mesh command with
// -ldflags "-X istio.io/pkg/version.buildHub=myhub -X istio.io/pkg/version.buildVersion=mytag"
// results in these values showing up in a generated manifest.
func TestLDFlags(t *testing.T) {
	testDataDir = filepath.Join(operatorRootDir, "cmd/mesh/testdata/manifest-generate")
	tmpHub, tmpTag := version.DockerInfo.Hub, version.DockerInfo.Tag
	defer func() {
		version.DockerInfo.Hub, version.DockerInfo.Tag = tmpHub, tmpTag
	}()
	version.DockerInfo.Hub = "testHub"
	version.DockerInfo.Tag = "testTag"
	l := NewLogger(true, os.Stdout, os.Stderr)
	ysf, err := yamlFromSetFlags([]string{"installPackagePath=" + liveInstallPackageDir}, false, l)
	if err != nil {
		t.Fatal(err)
	}
	_, iops, err := GenerateConfig(nil, ysf, true, nil, l)
	if err != nil {
		t.Fatal(err)
	}
	if iops.Hub != version.DockerInfo.Hub || iops.Tag != version.DockerInfo.Tag {
		t.Fatalf("DockerInfoHub, DockerInfoTag got: %s,%s, want: %s, %s", iops.Hub, iops.Tag, version.DockerInfo.Hub, version.DockerInfo.Tag)
	}
}

func runTestGroup(t *testing.T, tests testGroup) {
	testDataDir = filepath.Join(operatorRootDir, testDataSubdir)
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			inPath := filepath.Join(testDataDir, "input", tt.desc+".yaml")
			outputSuffix := goldenFileSuffixHideChangesInReview
			if tt.showOutputFileInPullRequest {
				outputSuffix = goldenFileSuffixShowChangesInReview
			}
			outPath := filepath.Join(testDataDir, "output", tt.desc+outputSuffix)

			var filenames []string
			if !tt.noInput {
				filenames = []string{inPath}
			}

			got, err := runManifestGenerate(filenames, tt.flags, tt.chartSource)
			if err != nil {
				t.Fatal(err)
			}

			if tt.outputDir != "" {
				got, err = util.ReadFilesWithFilter(tt.outputDir, func(fileName string) bool {
					return strings.HasSuffix(fileName, ".yaml")
				})
				if err != nil {
					t.Fatal(err)
				}
			}

			diffSelect := "*:*:*"
			if tt.diffSelect != "" {
				diffSelect = tt.diffSelect
				got, err = compare.SelectAndIgnoreFromOutput(got, diffSelect, "")
				if err != nil {
					t.Errorf("error selecting from output manifest: %v", err)
				}
			}

			if refreshGoldenFiles() {
				t.Logf("Refreshing golden file for %s", outPath)
				if err := ioutil.WriteFile(outPath, []byte(got), 0644); err != nil {
					t.Error(err)
				}
			}

			want, err := readFile(outPath)
			if err != nil {
				t.Fatal(err)
			}

			for _, v := range []bool{true, false} {
				diff, err := compare.ManifestDiffWithRenameSelectIgnore(got, want,
					"", diffSelect, tt.diffIgnore, v)
				if err != nil {
					t.Fatal(err)
				}
				if diff != "" {
					t.Errorf("%s: got:\n%s\nwant:\n%s\n(-got, +want)\n%s\n", tt.desc, "", "", diff)
				}
			}

		})
	}
}

// runManifestGenerate runs the manifest generate command. If filenames is set, passes the given filenames as -f flag,
// flags is passed to the command verbatim. If you set both flags and path, make sure to not use -f in flags.
func runManifestGenerate(filenames []string, flags string, chartSource chartSourceType) (string, error) {
	args := "manifest generate"
	for _, f := range filenames {
		args += " -f " + f
	}
	if flags != "" {
		args += " " + flags
	}
	switch chartSource {
	case snapshotCharts:
		args += " --set installPackagePath=" + filepath.Join(testDataDir, "data-snapshot")
	case liveCharts:
		args += " --set installPackagePath=" + liveInstallPackageDir
	case compiledInCharts:
	default:
	}
	return runCommand(args)
}

func createTempDirOrFail(t *testing.T, prefix string) string {
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func removeDirOrFail(t *testing.T, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatal(err)
	}
}

func createLocalReleaseCharts() (string, error) {
	releaseDir, err := ioutil.TempDir(os.TempDir(), "istio-test-release-*")
	if err != nil {
		return "", err
	}
	releaseSubDir := filepath.Join(releaseDir, istioTestVersion, helm.OperatorSubdirFilePath)
	cmd := exec.Command("../../release/create_release_charts.sh", "-o", releaseSubDir)
	if stdo, err := cmd.Output(); err != nil {
		return "", fmt.Errorf("%s: \n%s", err, string(stdo))
	}
	return releaseDir, nil
}
