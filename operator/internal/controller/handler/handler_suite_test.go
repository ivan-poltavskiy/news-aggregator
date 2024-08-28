package handler_test

import (
	v1 "com.teamdev/news-aggregator/api/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}

var _ = BeforeSuite(func() {
	_ = v1.AddToScheme(scheme.Scheme)
})
