package mykind

import (
	"context"

	mykindv1alpha1 "github.com/SamYuan1990/OperatorLearning/learning/pkg/apis/mykind/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_mykind")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Mykind Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMykind{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mykind-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Mykind
	err = c.Watch(&source.Kind{Type: &mykindv1alpha1.Mykind{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Mykind
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mykindv1alpha1.Mykind{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileMykind implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMykind{}

// ReconcileMykind reconciles a Mykind object
type ReconcileMykind struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Mykind object and makes changes based on the state read
// and what is in the Mykind.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMykind) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Mykind")

	// Fetch the Mykind instance
	instance := &mykindv1alpha1.Mykind{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	configMap := NewConfigMap(instance)
	reqLogger.Info(configMap.String())
	r.client.Create(context.TODO(), configMap)
	selector := NewConfigMapKeySelector(configMap)
	envSource := NewEnvVarSource(selector)
	service := NewService(instance)
	r.client.Create(context.TODO(), service)
	// Define a new Pod object
	pod := NewPodForCR(instance, envSource)

	// Set Mykind instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

func NewService(cr *mykindv1alpha1.Mykind) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}
	m1 := make(map[string]string)
	m1["app"] = cr.Name
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-service",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: 80,
				},
			},
			Selector: m1,
		},
	}
}

// newConfigMap
func NewConfigMap(cr *mykindv1alpha1.Mykind) *corev1.ConfigMap {
	labels := map[string]string{
		"app": cr.Name,
	}
	m1 := make(map[string]string)
	m1["a"] = cr.Spec.EnvsValue
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-configmap",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Data: m1,
	}
}

func NewConfigMapKeySelector(configmap *corev1.ConfigMap) *corev1.ConfigMapKeySelector {
	return &corev1.ConfigMapKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: configmap.ObjectMeta.Name,
		},
		Key: "a",
	}
}

func NewEnvVarSource(configMapSelector *corev1.ConfigMapKeySelector) *corev1.EnvVarSource {
	return &corev1.EnvVarSource{
		ConfigMapKeyRef: configMapSelector,
	}
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func NewPodForCR(cr *mykindv1alpha1.Mykind, envVarSource *corev1.EnvVarSource) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	image := "busybox"
	if cr.Spec.Image != "" {
		image = cr.Spec.Image
	}

	EnvVar := corev1.EnvVar{
		Name:      "dummy",
		ValueFrom: envVarSource,
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   image,
					Command: []string{"sleep", "3600"},
					Env:     []corev1.EnvVar{EnvVar},
				},
			},
		},
	}
}
