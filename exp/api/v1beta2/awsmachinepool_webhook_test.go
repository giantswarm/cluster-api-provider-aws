/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

import (
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	utildefaulting "sigs.k8s.io/cluster-api/util/defaulting"
)

func TestAWSMachinePoolDefault(t *testing.T) {
	m := &AWSMachinePool{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
	t.Run("for AWSCluster", utildefaulting.DefaultValidateTest(m))
	m.Default()
	g := NewWithT(t)
	g.Expect(m.Spec.DefaultCoolDown.Duration).To(BeNumerically(">=", 0))
}

func TestAWSMachinePoolValidateCreate(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name             string
		pool             *AWSMachinePool
		wantErrToContain *string
	}{
		{
			name: "pool with valid tags is accepted",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
				},
			},
			wantErrToContain: nil,
		},
		{
			name: "invalid tags are rejected",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1":                    "value-1",
						"":                         "value-2",
						strings.Repeat("CAPI", 33): "value-3",
						"key-4":                    strings.Repeat("CAPI", 65),
					},
				},
			},
			wantErrToContain: ptr.To[string]("additionalTags"),
		},
		{
			name: "Should fail if additional security groups are provided with both ID and Filters",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AWSLaunchTemplate: AWSLaunchTemplate{AdditionalSecurityGroups: []infrav1.AWSResourceReference{{
						ID: aws.String("sg-1"),
						Filters: []infrav1.Filter{
							{
								Name:   "sg-1",
								Values: []string{"test"},
							},
						},
					}}},
				},
			},
			wantErrToContain: ptr.To[string]("filter"),
		},
		{
			name: "Should fail if both subnet ID and filters passed in AWSMachinePool spec",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
					Subnets: []infrav1.AWSResourceReference{
						{
							ID:      ptr.To[string]("subnet-id"),
							Filters: []infrav1.Filter{{Name: "filter_name", Values: []string{"filter_value"}}},
						},
					},
				},
			},
			wantErrToContain: ptr.To[string]("filter"),
		},
		{
			name: "Should pass if either subnet ID or filters passed in AWSMachinePool spec",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
					Subnets: []infrav1.AWSResourceReference{
						{
							ID: ptr.To[string]("subnet-id"),
						},
					},
				},
			},
			wantErrToContain: nil,
		},
		{
			name: "Ensure root volume with device name works (for clusterctl move)",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AWSLaunchTemplate: AWSLaunchTemplate{
						RootVolume: &infrav1.Volume{
							DeviceName: "name",
							Type:       "gp2",
							Size:       *aws.Int64(8),
						},
					},
				},
			},
			wantErrToContain: nil,
		},
		{
			name: "Should fail if both spot market options or mixed instances policy are set",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					MixedInstancesPolicy: &MixedInstancesPolicy{
						Overrides: []Overrides{{InstanceType: "t3.medium"}},
					},
					AWSLaunchTemplate: AWSLaunchTemplate{
						SpotMarketOptions: &infrav1.SpotMarketOptions{MaxPrice: aws.String("0.1")},
					},
				},
			},
			wantErrToContain: ptr.To[string]("spotMarketOptions"),
		},
		{
			name: "Should fail if MaxHealthyPercentage is set, but MinHealthyPercentage is not set",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					RefreshPreferences: &RefreshPreferences{MaxHealthyPercentage: aws.Int64(100)},
				},
			},
			wantErrToContain: ptr.To[string]("minHealthyPercentage"),
		},
		{
			name: "Should fail if the difference between MaxHealthyPercentage and MinHealthyPercentage is greater than 100",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					RefreshPreferences: &RefreshPreferences{
						MaxHealthyPercentage: aws.Int64(150),
						MinHealthyPercentage: aws.Int64(25),
					},
				},
			},
			wantErrToContain: ptr.To[string]("minHealthyPercentage"),
		},
		{
			name: "Should fail if lifecycle hook only has roleARN, but not notificationTargetARN",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AWSLifecycleHooks: []AWSLifecycleHook{
						{
							Name:                "the-hook",
							LifecycleTransition: LifecycleHookTransitionInstanceTerminating,
							RoleARN:             aws.String("role-arn"),
						},
					},
				},
			},
			wantErrToContain: ptr.To[string]("notificationTargetARN"),
		},
		{
			name: "Should fail if lifecycle hook only has notificationTargetARN, but not roleARN",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AWSLifecycleHooks: []AWSLifecycleHook{
						{
							Name:                  "the-hook",
							LifecycleTransition:   LifecycleHookTransitionInstanceTerminating,
							NotificationTargetARN: aws.String("notification-target-arn"),
						},
					},
				},
			},
			wantErrToContain: ptr.To[string]("roleARN"),
		},
		{
			name: "Should fail if the lifecycle hook heartbeat timeout is less than 30 seconds",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AWSLifecycleHooks: []AWSLifecycleHook{
						{
							Name:                  "the-hook",
							LifecycleTransition:   LifecycleHookTransitionInstanceTerminating,
							NotificationTargetARN: aws.String("notification-target-arn"),
							RoleARN:               aws.String("role-arn"),
							HeartbeatTimeout:      &metav1.Duration{Duration: 29 * time.Second},
						},
					},
				},
			},
			wantErrToContain: ptr.To[string]("heartbeatTimeout"),
		},
		{
			name: "Should fail if the lifecycle hook heartbeat timeout is more than 172800 seconds",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AWSLifecycleHooks: []AWSLifecycleHook{
						{
							Name:                  "the-hook",
							LifecycleTransition:   LifecycleHookTransitionInstanceTerminating,
							NotificationTargetARN: aws.String("notification-target-arn"),
							RoleARN:               aws.String("role-arn"),
							HeartbeatTimeout:      &metav1.Duration{Duration: 172801 * time.Second},
						},
					},
				},
			},
			wantErrToContain: ptr.To[string]("heartbeatTimeout"),
		},
		{
			name: "Should succeed on correct lifecycle hook",
			pool: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AWSLifecycleHooks: []AWSLifecycleHook{
						{
							Name:                  "the-hook",
							LifecycleTransition:   LifecycleHookTransitionInstanceTerminating,
							NotificationTargetARN: aws.String("notification-target-arn"),
							RoleARN:               aws.String("role-arn"),
							HeartbeatTimeout:      &metav1.Duration{Duration: 180 * time.Second},
						},
					},
				},
			},
			wantErrToContain: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warn, err := tt.pool.ValidateCreate()
			if tt.wantErrToContain != nil {
				g.Expect(err).ToNot(BeNil())
				if err != nil {
					g.Expect(err.Error()).To(ContainSubstring(*tt.wantErrToContain))
				}
			} else {
				g.Expect(err).To(Succeed())
			}
			// Nothing emits warnings yet
			g.Expect(warn).To(BeEmpty())
		})
	}
}

func TestAWSMachinePoolValidateUpdate(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name             string
		new              *AWSMachinePool
		old              *AWSMachinePool
		wantErrToContain *string
	}{
		{
			name: "adding tags is accepted",
			old: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
					},
				},
			},
			new: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
				},
			},
			wantErrToContain: nil,
		},
		{
			name: "adding invalid tags is rejected",
			old: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
					},
				},
			},
			new: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1":                    "value-1",
						"":                         "value-2",
						strings.Repeat("CAPI", 33): "value-3",
						"key-4":                    strings.Repeat("CAPI", 65),
					},
				},
			},
			wantErrToContain: ptr.To[string]("additionalTags"),
		},
		{
			name: "Should fail update if both subnetID and filters passed in AWSMachinePool spec",
			old: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
					},
				},
			},
			new: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
					Subnets: []infrav1.AWSResourceReference{
						{
							ID:      ptr.To[string]("subnet-id"),
							Filters: []infrav1.Filter{{Name: "filter_name", Values: []string{"filter_value"}}},
						},
					},
				},
			},
			wantErrToContain: ptr.To[string]("filter"),
		},
		{
			name: "Should pass update if either subnetID or filters passed in AWSMachinePool spec",
			old: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
					},
				},
			},
			new: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
					Subnets: []infrav1.AWSResourceReference{
						{
							ID: ptr.To[string]("subnet-id"),
						},
					},
				},
			},
			wantErrToContain: nil,
		},
		{
			name: "Should fail update if both spec.awsLaunchTemplate.SpotMarketOptions and spec.MixedInstancesPolicy are passed in AWSMachinePool spec",
			old: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					MixedInstancesPolicy: &MixedInstancesPolicy{
						Overrides: []Overrides{{InstanceType: "t3.medium"}},
					},
				},
			},
			new: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					MixedInstancesPolicy: &MixedInstancesPolicy{
						Overrides: []Overrides{{InstanceType: "t3.medium"}},
					},
					AWSLaunchTemplate: AWSLaunchTemplate{
						SpotMarketOptions: &infrav1.SpotMarketOptions{MaxPrice: ptr.To[string]("0.1")},
					},
				},
			},
			wantErrToContain: ptr.To[string]("spotMarketOptions"),
		},
		{
			name: "Should fail if MaxHealthyPercentage is set, but MinHealthyPercentage is not set",
			new: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					RefreshPreferences: &RefreshPreferences{MaxHealthyPercentage: aws.Int64(100)},
				},
			},
			wantErrToContain: ptr.To[string]("minHealthyPercentage"),
		},
		{
			name: "Should fail if the difference between MaxHealthyPercentage and MinHealthyPercentage is greater than 100",
			new: &AWSMachinePool{
				Spec: AWSMachinePoolSpec{
					RefreshPreferences: &RefreshPreferences{
						MaxHealthyPercentage: aws.Int64(150),
						MinHealthyPercentage: aws.Int64(25),
					},
				},
			},
			wantErrToContain: ptr.To[string]("minHealthyPercentage"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warn, err := tt.new.ValidateUpdate(tt.old.DeepCopy())
			if tt.wantErrToContain != nil {
				g.Expect(err).ToNot(BeNil())
				if err != nil {
					g.Expect(err.Error()).To(ContainSubstring(*tt.wantErrToContain))
				}
			} else {
				g.Expect(err).To(Succeed())
			}
			// Nothing emits warnings yet
			g.Expect(warn).To(BeEmpty())
		})
	}
}
