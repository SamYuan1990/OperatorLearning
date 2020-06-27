package mykind_test

import (
	mykindv1alpha1 "github.com/SamYuan1990/OperatorLearning/learning/pkg/apis/mykind/v1alpha1"

	"github.com/SamYuan1990/OperatorLearning/learning/pkg/controller/mykind"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MykindController", func() {

	Context("NewConfigMap", func() {
		It("should return a not nil configmap", func() {
			cr := mykindv1alpha1.MykindSpec{EnvsValue: "value"}
			instance := &mykindv1alpha1.Mykind{Spec: cr}
			configmap := mykind.NewConfigMap(instance)
			Expect(configmap).NotTo(BeNil())
			Expect(configmap.Data["a"]).To(Equal(cr.EnvsValue))
			Expect(configmap.Data["a"]).To(Equal("value"))

		})

		It("should able to filter out", func() {
			cr := mykindv1alpha1.MykindSpec{EnvsValue: "value"}
			instance := &mykindv1alpha1.Mykind{Spec: cr}
			configmap := mykind.NewConfigMap(instance)
			selector := mykind.NewConfigMapKeySelector(configmap)
			Expect(selector).NotTo(BeNil())
			Expect(selector.LocalObjectReference).NotTo(BeNil())
		})

		It("should able to coverent to env", func() {
			cr := mykindv1alpha1.MykindSpec{EnvsValue: "value"}
			instance := &mykindv1alpha1.Mykind{Spec: cr}
			configmap := mykind.NewConfigMap(instance)
			selector := mykind.NewConfigMapKeySelector(configmap)
			envSource := mykind.NewEnvVarSource(selector)
			Expect(envSource).NotTo(BeNil())
		})
	})

	Context("newPodForCR", func() {
		It("should use busybox as image by default", func() {
			cr := mykindv1alpha1.MykindSpec{}
			instance := &mykindv1alpha1.Mykind{Spec: cr}
			Pod := mykind.NewPodForCR(instance, nil)
			Expect("busybox").To(Equal(Pod.Spec.Containers[0].Image))
		})

		It("should use image by given", func() {
			cr := mykindv1alpha1.MykindSpec{Image: "test"}
			instance := &mykindv1alpha1.Mykind{Spec: cr}
			Pod := mykind.NewPodForCR(instance, nil)
			Expect("test").To(Equal(Pod.Spec.Containers[0].Image))
		})

		It("should pass envs if given", func() {
			cr := mykindv1alpha1.MykindSpec{EnvsValue: "value"}
			instance := &mykindv1alpha1.Mykind{Spec: cr}
			configmap := mykind.NewConfigMap(instance)
			selector := mykind.NewConfigMapKeySelector(configmap)
			envSource := mykind.NewEnvVarSource(selector)
			Pod := mykind.NewPodForCR(instance, envSource)
			Expect(envSource).To(Equal(Pod.Spec.Containers[0].Env[0].ValueFrom))
		})
	})
})
