package helpers

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helper Suite")
}

var _ = Describe("the helpers", func() {
	Describe("GetRegistryAddress", func() {
		It("should return error if passed empty string", func() {
			_, err := GetRegistryAddress("")
			Expect(err).To(HaveOccurred())
		})
		It("should return index.docker.io for image refs with no explicit registry", func() {
			host, err := GetRegistryAddress("watchtower")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("index.docker.io"))

			host, err = GetRegistryAddress("containrrr/watchtower")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("index.docker.io"))
		})
		It("should return index.docker.io for image refs with docker.io domain", func() {
			host, err := GetRegistryAddress("docker.io/watchtower")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("index.docker.io"))

			host, err = GetRegistryAddress("docker.io/containrrr/watchtower")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("index.docker.io"))
		})
		It("should return the host if passed an image name containing a local host", func() {
			host, err := GetRegistryAddress("henk:80/watchtower")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("henk:80"))

			host, err = GetRegistryAddress("localhost/watchtower")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("localhost"))
		})
		It("should return the server address if passed a fully qualified image name", func() {
			host, err := GetRegistryAddress("github.com/containrrr/config")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("github.com"))
		})
	})

	Describe("GetRegistryAddressForRequest", func() {
		It("should keep non-Docker Hub registries unchanged", func() {
			host, err := GetRegistryAddressForRequest("ghcr.io/containrrr/watchtower", "mirror.ccs.tencentyun.com")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("ghcr.io"))
		})
		It("should use the override for Docker Hub registries", func() {
			host, err := GetRegistryAddressForRequest("containrrr/watchtower", "mirror.ccs.tencentyun.com")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("mirror.ccs.tencentyun.com"))
		})
	})

	Describe("NormalizeRegistryHost", func() {
		It("should return a host from a mirror URL", func() {
			host, err := NormalizeRegistryHost("https://mirror.ccs.tencentyun.com")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("mirror.ccs.tencentyun.com"))
		})
		It("should accept host and port without a scheme", func() {
			host, err := NormalizeRegistryHost("mirror.example.com:5000")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("mirror.example.com:5000"))
		})
		It("should reject mirror URLs with paths", func() {
			_, err := NormalizeRegistryHost("https://mirror.example.com/v2")
			Expect(err).To(MatchError(ContainSubstring("must not include a path")))
		})
	})
})
