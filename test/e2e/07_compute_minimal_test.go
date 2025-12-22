/*
Copyright 2025.

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

package e2e

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Arubacloud/arubacloud-resource-operator/test/utils"
)

var _ = Describe("07-ComputeMinimal", Ordered, func() {
	const (
		projectName      = "aruba-test-compute-min"
		vpcName          = "aruba-test-compute-min"
		sgName           = "aruba-test-compute-min"
		ruleName         = "aruba-test-compute-min"
		subnetName       = "aruba-test-compute-min"
		blockStorageName = "aruba-test-compute-min"
		keyPairName      = "aruba-test-compute-min"
		namespace        = "default"
		testTimeout      = 20 * time.Minute
	)

	BeforeAll(func() {
		By("ensuring the namespace exists")
		cmd := exec.Command("kubectl", "get", "ns", namespace)
		_, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "namespace should exist")
	})

	AfterAll(func() {
		By("cleaning up resources in reverse order")

		By("deleting KeyPair")
		cmd := exec.Command("kubectl", "delete", "keypair", keyPairName, "-n", namespace, "--ignore-not-found")
		_, _ = utils.Run(cmd)

		By("deleting BlockStorage")
		cmd = exec.Command("kubectl", "delete", "blockstorage", blockStorageName, "-n", namespace, "--ignore-not-found")
		_, _ = utils.Run(cmd)

		By("deleting Subnet")
		cmd = exec.Command("kubectl", "delete", "subnet", subnetName, "-n", namespace, "--ignore-not-found")
		_, _ = utils.Run(cmd)

		By("deleting SecurityRule")
		cmd = exec.Command("kubectl", "delete", "securityrule", ruleName, "-n", namespace, "--ignore-not-found")
		_, _ = utils.Run(cmd)

		By("deleting SecurityGroup")
		cmd = exec.Command("kubectl", "delete", "securitygroup", sgName, "-n", namespace, "--ignore-not-found")
		_, _ = utils.Run(cmd)

		By("deleting VPC")
		cmd = exec.Command("kubectl", "delete", "vpc", vpcName, "-n", namespace, "--ignore-not-found")
		_, _ = utils.Run(cmd)

		By("deleting Project")
		cmd = exec.Command("kubectl", "delete", "project", projectName, "-n", namespace, "--ignore-not-found")
		_, _ = utils.Run(cmd)
	})

	It("should create Project resource", func(ctx SpecContext) {
		projectYAML, err := utils.LoadSampleManifest("arubacloud.com_v1alpha1_project.yaml", map[string]string{
			"__NAME__":      projectName,
			"__NAMESPACE__": namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		By("creating Project")
		cmd := exec.Command("kubectl", "apply", "-f", "-")
		cmd.Stdin = stringReader(projectYAML)
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())

		By("waiting for Project to be Created")
		Eventually(func() string {
			cmd := exec.Command("kubectl", "get", "project", projectName, "-n", namespace, "-o", "jsonpath={.status.phase}")
			output, _ := utils.Run(cmd)
			return output
		}, testTimeout, 10*time.Second).Should(Equal("Created"))

		By("verifying Project has resourceID")
		cmd = exec.Command("kubectl", "get", "project", projectName, "-n", namespace, "-o", "jsonpath={.status.resourceID}")
		output, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
	}, SpecTimeout(testTimeout))

	It("should create VPC resource", func(ctx SpecContext) {
		vpcYAML, err := utils.LoadSampleManifest("arubacloud.com_v1alpha1_vpc.yaml", map[string]string{
			"__NAME__":      vpcName,
			"__NAMESPACE__": namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		By("creating VPC")
		cmd := exec.Command("kubectl", "apply", "-f", "-")
		cmd.Stdin = stringReader(vpcYAML)
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())

		By("waiting for VPC to be Created")
		Eventually(func() string {
			cmd := exec.Command("kubectl", "get", "vpc", vpcName, "-n", namespace, "-o", "jsonpath={.status.phase}")
			output, _ := utils.Run(cmd)
			return output
		}, testTimeout, 10*time.Second).Should(Equal("Created"))

		By("verifying VPC has resourceID")
		cmd = exec.Command("kubectl", "get", "vpc", vpcName, "-n", namespace, "-o", "jsonpath={.status.resourceID}")
		output, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
	}, SpecTimeout(testTimeout))

	It("should create SecurityGroup resource", func(ctx SpecContext) {
		sgYAML, err := utils.LoadSampleManifest("arubacloud.com_v1alpha1_securitygroup.yaml", map[string]string{
			"__NAME__":      sgName,
			"__NAMESPACE__": namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		By("creating SecurityGroup")
		cmd := exec.Command("kubectl", "apply", "-f", "-")
		cmd.Stdin = stringReader(sgYAML)
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())

		By("waiting for SecurityGroup to be Created")
		Eventually(func() string {
			cmd := exec.Command("kubectl", "get", "securitygroup", sgName, "-n", namespace, "-o", "jsonpath={.status.phase}")
			output, _ := utils.Run(cmd)
			return output
		}, testTimeout, 10*time.Second).Should(Equal("Created"))

		By("verifying SecurityGroup has securityGroupID")
		cmd = exec.Command("kubectl", "get", "securitygroup", sgName, "-n", namespace, "-o", "jsonpath={.status.securityGroupID}")
		output, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
	}, SpecTimeout(testTimeout))

	It("should create SecurityRule resource", func(ctx SpecContext) {
		ruleYAML, err := utils.LoadSampleManifest("arubacloud.com_v1alpha1_securityrule_minimal.yaml", map[string]string{
			"__NAME__":              ruleName,
			"__NAMESPACE__":         namespace,
			"__SG_NAME__":           sgName,
			"__SG_NAMESPACE__":      namespace,
			"__PROJECT_NAME__":      projectName,
			"__PROJECT_NAMESPACE__": namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		By("creating SecurityRule")
		cmd := exec.Command("kubectl", "apply", "-f", "-")
		cmd.Stdin = stringReader(ruleYAML)
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())

		By("waiting for SecurityRule to be Created")
		Eventually(func() string {
			cmd := exec.Command("kubectl", "get", "securityrule", ruleName, "-n", namespace, "-o", "jsonpath={.status.phase}")
			output, _ := utils.Run(cmd)
			return output
		}, testTimeout, 10*time.Second).Should(Equal("Created"))

		By("verifying SecurityRule has securityRuleID")
		cmd = exec.Command("kubectl", "get", "securityrule", ruleName, "-n", namespace, "-o", "jsonpath={.status.securityRuleID}")
		output, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
	}, SpecTimeout(testTimeout))

	It("should create Subnet resource", func(ctx SpecContext) {
		subnetYAML, err := utils.LoadSampleManifest("arubacloud.com_v1alpha1_subnet_minimal.yaml", map[string]string{
			"__NAME__":              subnetName,
			"__NAMESPACE__":         namespace,
			"__VPC_NAME__":          vpcName,
			"__VPC_NAMESPACE__":     namespace,
			"__PROJECT_NAME__":      projectName,
			"__PROJECT_NAMESPACE__": namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		By("creating Subnet")
		cmd := exec.Command("kubectl", "apply", "-f", "-")
		cmd.Stdin = stringReader(subnetYAML)
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())

		By("waiting for Subnet to be Created")
		Eventually(func() string {
			cmd := exec.Command("kubectl", "get", "subnet", subnetName, "-n", namespace, "-o", "jsonpath={.status.phase}")
			output, _ := utils.Run(cmd)
			return output
		}, testTimeout, 10*time.Second).Should(Equal("Created"))

		By("verifying Subnet has subnetID")
		cmd = exec.Command("kubectl", "get", "subnet", subnetName, "-n", namespace, "-o", "jsonpath={.status.subnetID}")
		output, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
	}, SpecTimeout(testTimeout))

	It("should create BlockStorage boot volume", func(ctx SpecContext) {
		blockStorageYAML, err := utils.LoadSampleManifest("arubacloud.com_v1alpha1_blockstorage_minimal.yaml", map[string]string{
			"__NAME__":              blockStorageName,
			"__NAMESPACE__":         namespace,
			"__PROJECT_NAME__":      projectName,
			"__PROJECT_NAMESPACE__": namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		By("creating BlockStorage")
		cmd := exec.Command("kubectl", "apply", "-f", "-")
		cmd.Stdin = stringReader(blockStorageYAML)
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())

		By("waiting for BlockStorage to be Created")
		Eventually(func() string {
			cmd := exec.Command("kubectl", "get", "blockstorage", blockStorageName, "-n", namespace, "-o", "jsonpath={.status.phase}")
			output, _ := utils.Run(cmd)
			return output
		}, testTimeout, 10*time.Second).Should(Equal("Created"))

		By("verifying BlockStorage has resourceID")
		cmd = exec.Command("kubectl", "get", "blockstorage", blockStorageName, "-n", namespace, "-o", "jsonpath={.status.resourceID}")
		output, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
	}, SpecTimeout(testTimeout))

	It("should create KeyPair resource", func(ctx SpecContext) {
		keyPairYAML, err := utils.LoadSampleManifest("arubacloud.com_v1alpha1_keypair.yaml", map[string]string{
			"__NAME__":      keyPairName,
			"__NAMESPACE__": namespace,
		})
		Expect(err).NotTo(HaveOccurred())

		By("creating KeyPair")
		cmd := exec.Command("kubectl", "apply", "-f", "-")
		cmd.Stdin = stringReader(keyPairYAML)
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())

		By("waiting for KeyPair to be Created")
		Eventually(func() string {
			cmd := exec.Command("kubectl", "get", "keypair", keyPairName, "-n", namespace, "-o", "jsonpath={.status.phase}")
			output, _ := utils.Run(cmd)
			return output
		}, testTimeout, 10*time.Second).Should(Equal("Created"))

		By("verifying KeyPair has keyPairID")
		cmd = exec.Command("kubectl", "get", "keypair", keyPairName, "-n", namespace, "-o", "jsonpath={.status.keyPairID}")
		output, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).NotTo(BeEmpty())
	}, SpecTimeout(testTimeout))
})
