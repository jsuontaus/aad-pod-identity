// +build e2e

package helm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Azure/aad-pod-identity/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	chartName = "aad-pod-identity"
)

// InstallInput is the input for Install.
type InstallInput struct {
	Config                *framework.Config
	ManagedMode           bool
	NamespacedMode        bool
	BlockInstanceMetadata bool
}

// Install installs aad-pod-identity via Helm 3.
func Install(input InstallInput) {
	Expect(input.Config).NotTo(BeNil(), "input.Config is required for Helm.Install")

	cwd, err := os.Getwd()
	Expect(err).To(BeNil())

	// Change current working directory to repo root
	// Before installing aad-pod identity through Helm
	os.Chdir("../..")
	defer os.Chdir(cwd)

	args := append([]string{
		"install",
		chartName,
		"charts/aad-pod-identity",
		"--wait",
		fmt.Sprintf("--set=image.repository=%s", input.Config.Registry),
		fmt.Sprintf("--set=mic.tag=%s", input.Config.MICVersion),
		fmt.Sprintf("--set=nmi.tag=%s", input.Config.NMIVersion),
	})

	if input.Config.ImmutableUserMSIs != "" {
		args = append(args, fmt.Sprintf("--set=mic.immutableUserMSIs=%s", input.Config.ImmutableUserMSIs))
	}

	if input.ManagedMode {
		args = append(args, fmt.Sprintf("--set=operationMode=%s", "managed"))
	}

	if input.BlockInstanceMetadata {
		args = append(args, fmt.Sprintf("--set=nmi.blockInstanceMetadata=%t", input.BlockInstanceMetadata))
	}

	helm(args)
}

// Uninstall uninstalls aad-pod-identity via Helm 3.
func Uninstall() {
	args := []string{
		"uninstall",
		chartName,
	}

	helm(args)
}

// UpgradeInput is the input for Upgrade.
type UpgradeInput struct {
	Config                *framework.Config
	BlockInstanceMetadata bool
}

// Upgrade upgrades aad-pod-identity via Helm 3.
func Upgrade(input UpgradeInput) {
	Expect(input.Config).NotTo(BeNil(), "input.Config is required for Helm.Upgrade")

	cwd, err := os.Getwd()
	Expect(err).To(BeNil())

	// Change current working directory to repo root
	// Before installing aad-pod identity through Helm
	os.Chdir("../..")
	defer os.Chdir(cwd)

	args := append([]string{
		"upgrade",
		chartName,
		"charts/aad-pod-identity",
		"--wait",
		fmt.Sprintf("--set=image.repository=%s", input.Config.Registry),
		fmt.Sprintf("--set=mic.tag=%s", input.Config.MICVersion),
		fmt.Sprintf("--set=nmi.tag=%s", input.Config.NMIVersion),
	})

	if input.Config.ImmutableUserMSIs != "" {
		args = append(args, fmt.Sprintf("--set=mic.immutableUserMSIs=%s", input.Config.ImmutableUserMSIs))
	}

	if input.BlockInstanceMetadata {
		args = append(args, fmt.Sprintf("--set=nmi.blockInstanceMetadata=%t", input.BlockInstanceMetadata))
	}

	helm(args)
}

func helm(args []string) {
	By(fmt.Sprintf("helm %s", strings.Join(args, " ")))

	cmd := exec.Command("helm", args...)
	stdoutStderr, err := cmd.CombinedOutput()
	fmt.Printf("%s", stdoutStderr)

	Expect(err).To(BeNil())
}
