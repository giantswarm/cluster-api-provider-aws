/*
Copyright The Kubernetes Authors.

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

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	capiannotations "sigs.k8s.io/cluster-api/util/annotations"
)

// AWSNetworkReconciler reconciles a AWSNetwork object
type AWSNetworkReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsnetworks,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsnetworks/status,verbs=get;update;patch

func (r *AWSNetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.FromContext(ctx)

	// Fetch the AWSNetwork instance
	awsNetwork := &expinfrav1.AWSNetwork{}
	err := r.Get(ctx, req.NamespacedName, awsNetwork)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if capiannotations.HasPaused(awsNetwork) {
		log.Info("AWSNetwork is marked as paused, won't reconcile")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("awsNetwork", klog.KObj(awsNetwork))

	awsNetworkScope := &scope.AWSNetworkScope{
		Logger: *log,
	}

	// Handle deleted network
	if !awsNetwork.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(awsNetworkScope)
	}

	// Handle non-deleted network
	return r.reconcileNormal()
}

func (r *AWSNetworkReconciler) reconcileDelete(awsNetworkScope *scope.AWSNetworkScope) (ctrl.Result, error) {
	// TODO
	return ctrl.Result{}, nil
}

func (r *AWSNetworkReconciler) reconcileNormal() (ctrl.Result, error) {
	// TODO
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSNetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.AWSNetwork{}).
		Complete(r)
}
