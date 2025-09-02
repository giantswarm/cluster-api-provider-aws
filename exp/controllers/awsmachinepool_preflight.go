package controllers

import (
	"context"

	"github.com/blang/semver/v4"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api/util/version"
)

// runPreflightChecks determines if it's safe to update the machine pool's LaunchTemplate.
// Heavily based on the MachineSet PreflightChecks in Cluster API.
func (r *AWSMachinePoolReconciler) runPreflightChecks(ctx context.Context, machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper) (bool, error) {
	// log := ctrl.LoggerFrom(ctx)
	// If the MachinePoolPreflightChecks feature gate is disabled return early.
	if !feature.Gates.Enabled(feature.MachinePoolPreflightChecks) {
		return true, nil
	}

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
	cpVersion, ok, err := unstructured.NestedString(controlPlane.UnstructuredContent(), "spec", "version")
	if err != nil {
		return false, err
	}
	if !ok {
		return true, nil
	}
	_, err = semver.ParseTolerant(cpVersion)
	if err != nil {
		return false, err
	}

	// Run the control-plane-stable preflight check.
	ok, err = r.controlPlaneStablePreflightChecks(controlPlane)
	if !ok || err != nil {
		return false, err
	}

	return true, nil
}

// controlPlaneStablePreflightChecks makes sure that the Control Plane is in a stable state.
func (r *AWSMachinePoolReconciler) controlPlaneStablePreflightChecks(controlPlane *unstructured.Unstructured) (bool, error) {
	// Check that the control plane is not provisioning.
	// We can know if the control plane was previously created or is being created for the first
	// time by looking at controlplane.status.version. If the version in status is set to a valid
	// value then the control plane was already provisioned at a previous time. If not, we can
	// assume that the control plane is being created for the first time.
	_, ok, err := unstructured.NestedString(controlPlane.UnstructuredContent(), "status", "version")
	if !ok || err != nil {
		return false, err
	}

	// Check that the control plane is not upgrading.
	// A control plane is considered upgrading if spec.version is greater than status.version.
	// A control plane is considered not upgrading if the status or status.version is not set.
	isUpgrading, err := controlPlaneIsUpgrading(controlPlane)
	if err != nil {
		return false, err
	}
	if isUpgrading {
		return false, nil
	}

	return true, nil
}

func controlPlaneIsUpgrading(obj *unstructured.Unstructured) (bool, error) {
	specVersion, ok, err := unstructured.NestedString(obj.UnstructuredContent(), "spec", "version")
	if !ok || err != nil {
		return false, errors.Wrap(err, "failed to get control plane spec version")
	}
	specV, err := semver.ParseTolerant(specVersion)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse control plane spec version")
	}
	statusVersion, ok, err := unstructured.NestedString(obj.UnstructuredContent(), "status", "version")
	if !ok && err == nil { // status version is not yet set
		// If the status.version is not yet present in the object, it implies the
		// first machine of the control plane is provisioning. We can reasonably assume
		// that the control plane is not upgrading at this stage.
		return false, err
	}
	if !ok || err != nil {
		return false, err
	}
	statusV, err := semver.ParseTolerant(statusVersion)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse control plane status version")
	}

	// NOTE: we are considering the control plane upgrading when the version is greater
	// or when the version has a different build metadata.
	return version.Compare(specV, statusV, version.WithBuildTags()) >= 1, nil
}
