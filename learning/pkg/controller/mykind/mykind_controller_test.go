package mykind_test

import (
	mykindv1alpha1 "github.com/SamYuan1990/OperatorLearning/learning/pkg/apis/mykind/v1alpha1"

	"github.com/SamYuan1990/OperatorLearning/learning/pkg/controller/mykind"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MykindController", func() {

	Context("newPodForCR", func() {
		It("should use busybox as image by default", func() {
			cr := mykindv1alpha1.MykindSpec{}
			instance := &mykindv1alpha1.Mykind{Spec: cr}
			Pod := mykind.NewPodForCR(instance)
			Expect("busybox").To(Equal(Pod.Spec.Containers[0].Image))
		})

		It("should use image by given", func() {
			cr := mykindv1alpha1.MykindSpec{Image: "test"}
			instance := &mykindv1alpha1.Mykind{Spec: cr}
			Pod := mykind.NewPodForCR(instance)
			Expect("test").To(Equal(Pod.Spec.Containers[0].Image))
		})

		It("should pass envs if given", func() {
			cr := mykindv1alpha1.MykindSpec{EnvsValue: "value"}
			instance := &mykindv1alpha1.Mykind{Spec: cr}
			Pod := mykind.NewPodForCR(instance)
			Expect("value").To(Equal(Pod.Spec.Containers[0].Env[0].Value))
		})
	})

})
