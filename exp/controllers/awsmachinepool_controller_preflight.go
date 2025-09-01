package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// runPreflightChecks determines if it's safe to update the machine pool's LaunchTemplate.
// Heavily based on the MachineSet PreflightChecks in Cluster API.
func (r *AWSMachinePoolReconciler) runPreflightChecks(ctx context.Context, machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper) (bool, error) {
	// log := ctrl.LoggerFrom(ctx)
	//

	// Get the control plane object.
	controlPlane, err := clusterScope.UnstructuredControlPlane()
	if err != nil {
		return false, err
	}
	if controlPlane == nil {
		// If the cluster does not have a control plane reference then there is nothing to do. Return early.
		return true, nil

	}
	// cpKlogRef := klog.KRef(controlPlane.GetNamespace(), controlPlane.GetName())

	// If the Control Plane version is not set then we are dealing with a control plane that does not support version
	// or a control plane where the version is not set. In both cases we cannot perform any preflight checks as
	// we do not have enough information. Return early.
	_, ok, err := unstructured.NestedString(controlPlane.UnstructuredContent(), "spec", "version")
	if err != nil {
		return false, err
	}
	if !ok {
		return true, nil
	}
	// cpSemVer, err := semver.ParseTolerant(cpVersion)
	// if err != nil {
	// 	return false, err
	// }

	// // Run the control-plane-stable preflight check.
	// ok, err = r.controlPlaneStablePreflightChecks(ctx, machinePoolScope)

	return true, nil
}

// controlPlaneStablePreflightChecks makes sure that the Control Plane is in a stable state.
func (r *AWSMachinePoolReconciler) controlPlaneStablePreflightChecks(controlPlane *unstructured.Unstructured) (bool, error) {
	// Check that the control plane is not provisioning.
	// We can know if the control plane was previously created or is being created for the first
	// time by looking at controlplane.status.version. If the version in status is set to a valid
	// value then the control plane was already provisioned at a previous time. If not, we can
	// assume that the control plane is being created for the first time.
	// statusVersion, ok, err := unstructured.NestedString(controlPlane.UnstructuredContent(), "status", "version")
	// if err != nil {
	// }

	return false, nil
}
