package testutil_test

import (
	_ "embed"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/machadovilaca/operator-observability/pkg/testutil"
)

//go:embed testdata/metrics.txt
var metricsEndpointResponse string

var _ = Describe("FetchMetrics", func() {
	var server *ghttp.Server

	BeforeEach(func() {
		server = ghttp.NewServer()

		path := "/metrics"
		statusCode := 200

		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", path),
				ghttp.RespondWithPtr(&statusCode, &metricsEndpointResponse),
			))
	})

	AfterEach(func() {
		server.Close()
	})

	It("Should fetch all metrics with a given name", func() {
		mr, err := testutil.FetchMetric(server.URL()+"/metrics", "kubevirt_vm_count_total")
		Expect(err).ToNot(HaveOccurred())

		Expect(mr).To(HaveLen(2))
		Expect(mr[0].Name).To(Equal("kubevirt_vm_count_total"))
		Expect(mr[1].Name).To(Equal("kubevirt_vm_count_total"))
	})

	It("Should fetch all metrics with a given name and a single label pair", func() {
		mr, err := testutil.FetchMetric(server.URL()+"/metrics", "kubevirt_migration_count", "status", "failed")
		Expect(err).ToNot(HaveOccurred())

		Expect(mr).To(HaveLen(2))

		Expect(mr[0].Name).To(Equal("kubevirt_migration_count"))
		Expect(mr[0].Labels).To(HaveKeyWithValue("node", "node01"))
		Expect(mr[0].Labels).To(HaveKeyWithValue("status", "failed"))
		Expect(mr[0].Value).To(Equal(2.0))

		Expect(mr[1].Name).To(Equal("kubevirt_migration_count"))
		Expect(mr[1].Labels).To(HaveKeyWithValue("node", "node02"))
		Expect(mr[1].Labels).To(HaveKeyWithValue("status", "failed"))
		Expect(mr[1].Value).To(Equal(1.0))
	})

	It("Should fetch all metrics with a given name and multiple label pairs", func() {
		mr, err := testutil.FetchMetric(server.URL()+"/metrics", "kubevirt_vm_memory_usage_bytes", "namespace", "default", "vm_name", "vm1")
		Expect(err).ToNot(HaveOccurred())

		Expect(mr).To(HaveLen(1))
		Expect(mr[0].Name).To(Equal("kubevirt_vm_memory_usage_bytes"))
		Expect(mr[0].Labels).To(HaveKeyWithValue("namespace", "default"))
		Expect(mr[0].Labels).To(HaveKeyWithValue("vm_name", "vm1"))
		Expect(mr[0].Value).To(Equal(204800.0))
	})
})
