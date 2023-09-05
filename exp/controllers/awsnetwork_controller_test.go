/*
Copyright 2023 The Kubernetes Authors.

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

package controllers

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	"sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func newAWSNetwork(name, namespace string) expinfrav1.AWSNetwork {
	return expinfrav1.AWSNetwork{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: expinfrav1.AWSNetworkSpec{
			Foo: "yipeee",
		},
	}
}

func TestAWSNetworkReconcilerIntegration(t *testing.T) {
	var (
		reconciler AWSNetworkReconciler
		mockCtrl   *gomock.Controller
		ctx        context.Context
	)

	setup := func(t *testing.T) {
		t.Helper()
		mockCtrl = gomock.NewController(t)
		reconciler = AWSNetworkReconciler{
			Client: testEnv.Client,
		}
		ctx = context.TODO()
	}

	teardown := func() {
		mockCtrl.Finish()
	}

	t.Run("Should successfully reconcile AWSNetwork for managed VPC", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		elbMock := mocks.NewMockELBAPI(mockCtrl)
		expect := func(m *mocks.MockEC2APIMockRecorder, e *mocks.MockELBAPIMockRecorder) {
			// TODO add mocked calls here
		}
		expect(ec2Mock.EXPECT(), elbMock.EXPECT())

		setup(t)
		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())

		awsNetwork := newAWSNetwork("my-network", ns.Name)

		g.Expect(testEnv.Create(ctx, &awsNetwork)).To(Succeed())
		defer teardown()
		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, &awsNetwork, ns)).To(Succeed())
		})

		_, err = reconciler.Reconcile(ctx, ctrl.Request{
			NamespacedName: client.ObjectKey{
				Namespace: awsNetwork.Namespace,
				Name:      awsNetwork.Name,
			},
		})
		g.Expect(err).To(BeNil())
	})
}
