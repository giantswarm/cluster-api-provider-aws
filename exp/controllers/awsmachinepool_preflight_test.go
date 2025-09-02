package controllers

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilfeature "k8s.io/component-base/featuregate/testing"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// Ported over from the CAPI MachineSet PreflightChecks tests.
func TestMachineSetReconciler_runPreflightChecks(t *testing.T) {
	ns := "ns1"

	controlPlaneWithNoVersion := ControlPlane(ns, "cp1").Build()

	controlPlaneWithInvalidVersion := ControlPlane(ns, "cp1").
		WithVersion("v1.25.6.0").Build()

	controlPlaneProvisioning := ControlPlane(ns, "cp1").
		WithVersion("v1.25.6").Build()

	controlPlaneUpgrading := ControlPlane(ns, "cp1").
		WithVersion("v1.26.2").
		WithStatusFields(map[string]any{
			"status.version": "v1.25.2",
		}).
		Build()

	controlPlaneStable := ControlPlane(ns, "cp1").
		WithVersion("v1.26.2").
		WithStatusFields(map[string]any{
			"status.version": "v1.26.2",
		}).
		Build()

	controlPlaneStable128 := ControlPlane(ns, "cp1").
		WithVersion("v1.28.0").
		WithStatusFields(map[string]any{
			"status.version": "v1.28.0",
		}).
		Build()

	t.Run("should run preflight checks if the feature gate is enabled", func(t *testing.T) {
		defer utilfeature.SetFeatureGateDuringTest(t, feature.Gates, feature.MachinePoolPreflightChecks, true)()

		tests := []struct {
			name             string
			cluster          *clusterv1.Cluster
			controlPlane     *unstructured.Unstructured
			machinePool      *expclusterv1.MachinePool
			infraMachinePool *expinfrav1.AWSMachinePool
			wantPass         bool
			wantErr          bool
		}{
			{
				name:             "should pass if cluster has no control plane",
				cluster:          &clusterv1.Cluster{},
				machinePool:      &expclusterv1.MachinePool{},
				infraMachinePool: &expinfrav1.AWSMachinePool{},
				wantPass:         true,
			},
			{
				name: "should pass if the control plane version is not defined",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneWithNoVersion),
					},
				},
				controlPlane: controlPlaneWithNoVersion,
				machinePool:  &expclusterv1.MachinePool{},
				wantPass:     true,
			},
			{
				name: "should error if the control plane version is invalid",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneWithInvalidVersion),
					},
				},
				controlPlane: controlPlaneWithInvalidVersion,
				machinePool:  &expclusterv1.MachinePool{},
				wantErr:      true,
			},
			// {
			// 	name: "should pass if all preflight checks are skipped",
			// 	cluster: &clusterv1.Cluster{
			// 		ObjectMeta: metav1.ObjectMeta{
			// 			Namespace: ns,
			// 		},
			// 		Spec: clusterv1.ClusterSpec{
			// 			ControlPlaneRef: objToRef(controlPlaneUpgrading),
			// 		},
			// 	},
			// 	controlPlane: controlPlaneUpgrading,
			// 	machinePool: &expclusterv1.MachinePool{
			// 		ObjectMeta: metav1.ObjectMeta{
			// 			Namespace: ns,
			// 			Annotations: map[string]string{
			// 				clusterv1.MachineSetSkipPreflightChecksAnnotation: string(clusterv1.MachineSetPreflightCheckAll),
			// 			},
			// 		},
			// 	},
			// 	wantPass: true,
			// },
			{
				name: "control plane preflight check: should fail if the control plane is provisioning",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneProvisioning),
					},
				},
				controlPlane: controlPlaneProvisioning,
				machinePool:  &expclusterv1.MachinePool{},
				wantPass:     false,
			},
			{
				name: "control plane preflight check: should fail if the control plane is upgrading",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneUpgrading),
					},
				},
				controlPlane: controlPlaneUpgrading,
				machinePool:  &expclusterv1.MachinePool{},
				wantPass:     false,
			},
			// {
			// 	name: "control plane preflight check: should pass if the control plane is upgrading but the preflight check is skipped",
			// 	cluster: &clusterv1.Cluster{
			// 		ObjectMeta: metav1.ObjectMeta{
			// 			Namespace: ns,
			// 		},
			// 		Spec: clusterv1.ClusterSpec{
			// 			ControlPlaneRef: objToRef(controlPlaneUpgrading),
			// 		},
			// 	},
			// 	controlPlane: controlPlaneUpgrading,
			// 	machinePool: &expclusterv1.MachinePool{
			// 		ObjectMeta: metav1.ObjectMeta{
			// 			Namespace: ns,
			// 			Annotations: map[string]string{
			// 				clusterv1.MachineSetSkipPreflightChecksAnnotation: string(clusterv1.MachineSetPreflightCheckControlPlaneIsStable),
			// 			},
			// 		},
			// 		Spec: expclusterv1.MachinePoolSpec{
			// 			Template: clusterv1.MachineTemplateSpec{
			// 				Spec: clusterv1.MachineSpec{
			// 					Version:   ptr.To("v1.26.2"),
			// 					Bootstrap: clusterv1.Bootstrap{ConfigRef: &corev1.ObjectReference{Kind: "KubeadmConfigTemplate"}},
			// 				},
			// 			},
			// 		},
			// 	},
			// 	wantPass: true,
			// },
			{
				name: "control plane preflight check: should pass if the control plane is stable",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool:  &expclusterv1.MachinePool{},
				wantPass:     true,
			},
			{
				name: "should pass if the machine pool version is not defined",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{},
				},
				wantPass: true,
			},
			{
				name: "should error if the machine pool version is invalid",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{

						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.27.0.0"),
							},
						},
					},
				},
				wantErr: true,
			},
			{
				name: "kubernetes version preflight check: should fail if the machine pool minor version is greater than control plane minor version",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.27.0"),
							},
						},
					},
				},
				wantPass: false,
			},
			{
				name: "kubernetes version preflight check: should fail if the machine pool minor version is 4 older than control plane minor version for >= v1.28",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable128),
					},
				},
				controlPlane: controlPlaneStable128,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.24.0"),
							},
						},
					},
				},
				wantPass: false,
			},
			{
				name: "kubernetes version preflight check: should fail if the machine pool minor version is 3 older than control plane minor version for < v1.28",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.23.0"),
							},
						},
					},
				},
				wantPass: false,
			},
			// {
			// 	name: "kubernetes version preflight check: should pass if the machine pool minor version is greater than control plane minor version but the preflight check is skipped",
			// 	cluster: &clusterv1.Cluster{
			// 		ObjectMeta: metav1.ObjectMeta{
			// 			Namespace: ns,
			// 		},
			// 		Spec: clusterv1.ClusterSpec{
			// 			ControlPlaneRef: objToRef(controlPlaneStable),
			// 		},
			// 	},
			// 	controlPlane: controlPlaneStable,
			// 	machinePool: &expclusterv1.MachinePool{
			// 		ObjectMeta: metav1.ObjectMeta{
			// 			Namespace: ns,
			// 			Annotations: map[string]string{
			// 				clusterv1.MachineSetSkipPreflightChecksAnnotation: string(clusterv1.MachineSetPreflightCheckKubernetesVersionSkew),
			// 			},
			// 		},
			// 		Spec: expclusterv1.MachinePoolSpec{
			// 			Template: clusterv1.MachineTemplateSpec{
			// 				Spec: clusterv1.MachineSpec{
			// 					Version: ptr.To("v1.27.0"),
			// 				},
			// 			},
			// 		},
			// 	},
			// 	wantPass: true,
			// },
			{
				name: "kubernetes version preflight check: should pass if the machine pool minor version and control plane version conform to kubernetes version skew policy >= v1.28",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable128),
					},
				},
				controlPlane: controlPlaneStable128,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.25.0"),
							},
						},
					},
				},
				wantPass: true,
			},
			{
				name: "kubernetes version preflight check: should pass if the machine pool minor version and control plane version conform to kubernetes version skew policy < v1.28",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.24.0"),
							},
						},
					},
				},
				wantPass: true,
			},
			{
				name: "kubeadm version preflight check: should fail if the machine pool version is not equal (major+minor) to control plane version when using kubeadm bootstrap provider",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.25.5"),
								Bootstrap: clusterv1.Bootstrap{ConfigRef: &corev1.ObjectReference{
									APIVersion: bootstrapv1.GroupVersion.String(),
									Kind:       "KubeadmConfigTemplate",
								}},
							},
						},
					},
				},
				wantPass: false,
			},
			{
				name: "kubeadm version preflight check: should pass if the machine pool is not using kubeadm bootstrap provider",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.25.0"),
							},
						},
					},
				},
				wantPass: true,
			},
			{
				name: "kubeadm version preflight check: should pass if the machine pool version and control plane version do not conform to kubeadm version skew policy but the preflight check is skipped",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
						Annotations: map[string]string{
							clusterv1.MachineSetSkipPreflightChecksAnnotation: "foobar," + string(clusterv1.MachineSetPreflightCheckKubeadmVersionSkew),
						},
					},
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.25.0"),
								Bootstrap: clusterv1.Bootstrap{ConfigRef: &corev1.ObjectReference{
									APIVersion: bootstrapv1.GroupVersion.String(),
									Kind:       "KubeadmConfigTemplate",
								}},
							},
						},
					},
				},
				wantPass: true,
			},
			{
				name: "kubeadm version preflight check: should pass if the machine pool version and control plane version conform to kubeadm version skew when using kubeadm bootstrap provider",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.26.2"),
								Bootstrap: clusterv1.Bootstrap{ConfigRef: &corev1.ObjectReference{
									APIVersion: bootstrapv1.GroupVersion.String(),
									Kind:       "KubeadmConfigTemplate",
								}},
							},
						},
					},
				},
				wantPass: true,
			},
			{
				name: "kubeadm version preflight check: should error if the bootstrap ref APIVersion is invalid",
				cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: clusterv1.ClusterSpec{
						ControlPlaneRef: objToRef(controlPlaneStable),
					},
				},
				controlPlane: controlPlaneStable,
				machinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns,
					},
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Version: ptr.To("v1.26.2"),
								Bootstrap: clusterv1.Bootstrap{ConfigRef: &corev1.ObjectReference{
									APIVersion: "bootstrap.cluster.x-k8s.io/v1beta1/invalid",
									Kind:       "KubeadmConfigTemplate",
								}},
							},
						},
					},
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				g := NewWithT(t)
				objs := []client.Object{}
				if tt.controlPlane != nil {
					objs = append(objs, tt.controlPlane)
				}
				fakeClient := fake.NewClientBuilder().WithObjects(objs...).Build()
				r := &AWSMachinePoolReconciler{
					Client: fakeClient,
				}
				machinePoolScope := &scope.MachinePoolScope{
					Client:         fakeClient,
					Cluster:        tt.cluster,
					MachinePool:    tt.machinePool,
					AWSMachinePool: tt.infraMachinePool,
				}
				clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
					Client:     fakeClient,
					Cluster:    tt.cluster,
					AWSCluster: &v1beta2.AWSCluster{},
				})
				g.Expect(err).ToNot(HaveOccurred())

				result, err := r.runPreflightChecks(ctx, machinePoolScope, clusterScope)
				if tt.wantErr {
					g.Expect(err).To(HaveOccurred())
				} else {
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(result).To(Equal(tt.wantPass))
				}
			})
		}
	})

	t.Run("should not run the preflight checks if the feature gate is disabled", func(t *testing.T) {
		defer utilfeature.SetFeatureGateDuringTest(t, feature.Gates, feature.MachinePoolPreflightChecks, false)()

		g := NewWithT(t)
		cluster := &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns,
			},
			Spec: clusterv1.ClusterSpec{
				ControlPlaneRef: objToRef(controlPlaneUpgrading),
			},
		}
		controlPlane := controlPlaneUpgrading
		machinePool := &expclusterv1.MachinePool{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns,
			},
			Spec: expclusterv1.MachinePoolSpec{
				Template: clusterv1.MachineTemplateSpec{
					Spec: clusterv1.MachineSpec{
						Version: ptr.To("v1.26.0"),
						Bootstrap: clusterv1.Bootstrap{ConfigRef: &corev1.ObjectReference{
							APIVersion: bootstrapv1.GroupVersion.String(),
							Kind:       "KubeadmConfigTemplate",
						}},
					},
				},
			},
		}
		infraMachinePool := &expinfrav1.AWSMachinePool{}
		fakeClient := fake.NewClientBuilder().WithObjects(controlPlane).Build()
		r := &AWSMachinePoolReconciler{Client: fakeClient}
		machinePoolScope := &scope.MachinePoolScope{
			Client:         fakeClient,
			Cluster:        cluster,
			MachinePool:    machinePool,
			AWSMachinePool: infraMachinePool,
		}
		clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
			Client:     fakeClient,
			Cluster:    cluster,
			AWSCluster: &v1beta2.AWSCluster{},
		})
		g.Expect(err).ToNot(HaveOccurred())

		result, err := r.runPreflightChecks(ctx, machinePoolScope, clusterScope)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(result).To(BeTrue())
	})
}

// Copied from the sigs.k8s.io/cluster-api/util/test/builder package, which was made public in CAPI 1.9.
var (
	// ControlPlaneGroupVersion is group version used for control plane objects.
	ControlPlaneGroupVersion = schema.GroupVersion{Group: "controlplane.cluster.x-k8s.io", Version: "v1beta1"}

	// GenericControlPlaneKind is the Kind for the GenericControlPlane.
	GenericControlPlaneKind = "GenericControlPlane"
)

// ControlPlaneBuilder holds the variables and objects needed to build a generic object for cluster.spec.controlPlaneRef.
type ControlPlaneBuilder struct {
	obj *unstructured.Unstructured
}

// ControlPlane returns a ControlPlaneBuilder with the given name and Namespace.
func ControlPlane(namespace, name string) *ControlPlaneBuilder {
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion(ControlPlaneGroupVersion.String())
	obj.SetKind(GenericControlPlaneKind)
	obj.SetNamespace(namespace)
	obj.SetName(name)
	return &ControlPlaneBuilder{
		obj: obj,
	}
}

// WithReplicas sets the number of replicas for the ControlPlaneBuilder.
func (c *ControlPlaneBuilder) WithReplicas(replicas int64) *ControlPlaneBuilder {
	if err := unstructured.SetNestedField(c.obj.Object, replicas, "spec", "replicas"); err != nil {
		panic(err)
	}
	return c
}

// WithVersion adds the passed version to the ControlPlaneBuilder.
func (c *ControlPlaneBuilder) WithVersion(version string) *ControlPlaneBuilder {
	if err := unstructured.SetNestedField(c.obj.Object, version, "spec", "version"); err != nil {
		panic(err)
	}
	return c
}

// WithSpecFields sets a map of spec fields on the unstructured object. The keys in the map represent the path and the value corresponds
// to the value of the spec field.
//
// Note: all the paths should start with "spec."
//
//	Example map: map[string]any{
//	    "spec.version": "v1.2.3",
//	}.
func (c *ControlPlaneBuilder) WithSpecFields(fields map[string]any) *ControlPlaneBuilder {
	setSpecFields(c.obj, fields)
	return c
}

// WithStatusFields sets a map of status fields on the unstructured object. The keys in the map represent the path and the value corresponds
// to the value of the status field.
//
// Note: all the paths should start with "status."
//
//	Example map: map[string]any{
//	    "status.version": "v1.2.3",
//	}.
func (c *ControlPlaneBuilder) WithStatusFields(fields map[string]any) *ControlPlaneBuilder {
	setStatusFields(c.obj, fields)
	return c
}

// Build generates an Unstructured object from the information passed to the ControlPlaneBuilder.
func (c *ControlPlaneBuilder) Build() *unstructured.Unstructured {
	return c.obj
}

// setSpecFields sets fields in an unstructured object from a map.
func setSpecFields(obj *unstructured.Unstructured, fields map[string]any) {
	for k, v := range fields {
		fieldParts := strings.Split(k, ".")
		if len(fieldParts) == 0 {
			panic(fmt.Errorf("fieldParts invalid"))
		}
		if fieldParts[0] != "spec" {
			panic(fmt.Errorf("can not set fields outside spec"))
		}
		if err := unstructured.SetNestedField(obj.UnstructuredContent(), v, strings.Split(k, ".")...); err != nil {
			panic(err)
		}
	}
}

// setStatusFields sets fields in an unstructured object from a map.
func setStatusFields(obj *unstructured.Unstructured, fields map[string]any) {
	for k, v := range fields {
		fieldParts := strings.Split(k, ".")
		if len(fieldParts) == 0 {
			panic(fmt.Errorf("fieldParts invalid"))
		}
		if fieldParts[0] != "status" {
			panic(fmt.Errorf("can not set fields outside status"))
		}
		if err := unstructured.SetNestedField(obj.UnstructuredContent(), v, strings.Split(k, ".")...); err != nil {
			panic(err)
		}
	}
}

// objToRef returns a reference to the given object.
// Note: This function only operates on Unstructured instead of client.Object
// because it is only safe to assume for Unstructured that the GVK is set.
func objToRef(obj *unstructured.Unstructured) *corev1.ObjectReference {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return &corev1.ObjectReference{
		Kind:       gvk.Kind,
		APIVersion: gvk.GroupVersion().String(),
		Namespace:  obj.GetNamespace(),
		Name:       obj.GetName(),
	}
}
