package util

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
)

func TestUtil(t *testing.T) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Util Suite")
}

var _ = Describe("utils", func() {
	Context("NewLNodeWrapper", func() {
		// Cant test the error condition for GenerateNodeKey
		It("should work", func() {
			_, err := NewLNodeWrapper("")
			Expect(err).Should(BeNil())
		})
		It("should error for bad db path", func() {
			_, err := NewLNodeWrapper("/not/a/real/path")
			Expect(err).Should(HaveOccurred())
		})
	})
	Context("GenerateRandomPort", func() {
		// Create 100 table entries
		entries := func() []TableEntry {
			ent := make([]TableEntry, 100)
			for i := range ent {
				ent[i] = Entry(fmt.Sprintf("test %d", i))
			}
			return ent
		}()
		DescribeTable("test 100 times",
			func() {
				rn := GenerateRandomPort()
				Expect(rn).Should(BeNumerically("<=", 65535))
				Expect(rn).Should(BeNumerically(">=", 32767))
			},
			entries...,
		)
	})
	Context("NewUDPConn", func() {
		It("should work", func() {
			conn, err := NewUDPConn("127.0.0.1", GenerateRandomPort())
			Expect(err).Should(BeNil())

			err = conn.Close()
			Expect(err).Should(BeNil())
		})
		It("fail for bad IP addr", func() {
			_, err := NewUDPConn("123", GenerateRandomPort())
			Expect(err).Should(HaveOccurred())
		})
		It("fail for already allocated port", func() {
			port := GenerateRandomPort()
			conn, err := NewUDPConn("127.0.0.1", port)
			Expect(err).Should(BeNil())

			_, err = NewUDPConn("127.0.0.1", port)
			Expect(err).Should(HaveOccurred())
			conn.Close()
		})
	})
})
