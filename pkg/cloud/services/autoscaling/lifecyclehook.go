/*
Copyright 2024 The Kubernetes Authors.

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

package asg

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// DescribeLifecycleHooks returns the lifecycle hooks for the given AutoScalingGroup after retrieving them from the AWS API.
func (s *Service) DescribeLifecycleHooks(asgName string) ([]*expinfrav1.AWSLifecycleHook, error) {
	input := &autoscaling.DescribeLifecycleHooksInput{
		AutoScalingGroupName: aws.String(asgName),
	}

	out, err := s.ASGClient.DescribeLifecycleHooksWithContext(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe lifecycle hooks for AutoScalingGroup: %q", asgName)
	}

	hooks := make([]*expinfrav1.AWSLifecycleHook, len(out.LifecycleHooks))
	for i, hook := range out.LifecycleHooks {
		hooks[i] = s.SDKToLifecycleHook(hook)
	}

	return hooks, nil
}

func getPutLifecycleHookInput(asgName string, hook *expinfrav1.AWSLifecycleHook) (ret *autoscaling.PutLifecycleHookInput) {
	ret = &autoscaling.PutLifecycleHookInput{
		AutoScalingGroupName: aws.String(asgName),
		LifecycleHookName:    aws.String(hook.Name),
		LifecycleTransition:  aws.String(hook.LifecycleTransition.String()),

		// Optional
		RoleARN:               hook.RoleARN,
		NotificationTargetARN: hook.NotificationTargetARN,
		NotificationMetadata:  hook.NotificationMetadata,
	}

	// Optional parameters
	if hook.DefaultResult != nil {
		ret.DefaultResult = aws.String(hook.DefaultResult.String())
	}

	if hook.HeartbeatTimeout != nil {
		timeoutSeconds := hook.HeartbeatTimeout.Duration.Seconds()
		ret.HeartbeatTimeout = aws.Int64(int64(timeoutSeconds))
	}

	return
}

// CreateLifecycleHook creates a lifecycle hook for the given AutoScalingGroup.
func (s *Service) CreateLifecycleHook(ctx context.Context, asgName string, hook *expinfrav1.AWSLifecycleHook) error {
	input := getPutLifecycleHookInput(asgName, hook)

	if _, err := s.ASGClient.PutLifecycleHookWithContext(ctx, input); err != nil {
		return errors.Wrapf(err, "failed to create lifecycle hook %q for AutoScalingGroup: %q", hook.Name, asgName)
	}

	return nil
}

// UpdateLifecycleHook updates a lifecycle hook for the given AutoScalingGroup.
func (s *Service) UpdateLifecycleHook(ctx context.Context, asgName string, hook *expinfrav1.AWSLifecycleHook) error {
	input := getPutLifecycleHookInput(asgName, hook)

	if _, err := s.ASGClient.PutLifecycleHookWithContext(ctx, input); err != nil {
		return errors.Wrapf(err, "failed to update lifecycle hook %q for AutoScalingGroup: %q", hook.Name, asgName)
	}

	return nil
}

// DeleteLifecycleHook deletes a lifecycle hook for the given AutoScalingGroup.
func (s *Service) DeleteLifecycleHook(ctx context.Context, asgName string, hook *expinfrav1.AWSLifecycleHook) error {
	input := &autoscaling.DeleteLifecycleHookInput{
		AutoScalingGroupName: aws.String(asgName),
		LifecycleHookName:    aws.String(hook.Name),
	}

	if _, err := s.ASGClient.DeleteLifecycleHookWithContext(ctx, input); err != nil {
		return errors.Wrapf(err, "failed to delete lifecycle hook %q for AutoScalingGroup: %q", hook.Name, asgName)
	}

	return nil
}

// SDKToLifecycleHook converts an AWS SDK LifecycleHook to the CAPA lifecycle hook type.
func (s *Service) SDKToLifecycleHook(hook *autoscaling.LifecycleHook) *expinfrav1.AWSLifecycleHook {
	timeoutDuration := time.Duration(*hook.HeartbeatTimeout) * time.Second
	metav1Duration := metav1.Duration{Duration: timeoutDuration}
	defaultResult := expinfrav1.LifecycleHookDefaultResult(*hook.DefaultResult)
	lifecycleTransition := expinfrav1.LifecycleTransition(*hook.LifecycleTransition)

	return &expinfrav1.AWSLifecycleHook{
		Name:                  *hook.LifecycleHookName,
		DefaultResult:         &defaultResult,
		HeartbeatTimeout:      &metav1Duration,
		LifecycleTransition:   lifecycleTransition,
		NotificationTargetARN: hook.NotificationTargetARN,
		RoleARN:               hook.RoleARN,
		NotificationMetadata:  hook.NotificationMetadata,
	}
}

func getLifecycleHookSpecificationList(lifecycleHooks []expinfrav1.AWSLifecycleHook) (ret []*autoscaling.LifecycleHookSpecification) {
	for _, hook := range lifecycleHooks {
		spec := &autoscaling.LifecycleHookSpecification{
			LifecycleHookName:   aws.String(hook.Name),
			LifecycleTransition: aws.String(hook.LifecycleTransition.String()),

			// Optional
			RoleARN:               hook.RoleARN,
			NotificationTargetARN: hook.NotificationTargetARN,
			NotificationMetadata:  hook.NotificationMetadata,
		}

		// Optional parameters
		if hook.DefaultResult != nil {
			spec.DefaultResult = aws.String(hook.DefaultResult.String())
		}

		if hook.HeartbeatTimeout != nil {
			timeoutSeconds := hook.HeartbeatTimeout.Duration.Seconds()
			spec.HeartbeatTimeout = aws.Int64(int64(timeoutSeconds))
		}

		ret = append(ret, spec)
	}

	return
}

// ReconcileLifecycleHooks reconciles lifecycle hooks for an ASG
// by creating missing hooks, updating mismatching hooks and
// deleting extraneous hooks (except those specified in
// ignoreLifecycleHooks).
func ReconcileLifecycleHooks(ctx context.Context, asgService services.ASGInterface, asgName string, wantedLifecycleHooks []expinfrav1.AWSLifecycleHook, ignoreLifecycleHooks map[string]bool, storeConditionsOnObject conditions.Setter, log logger.Wrapper) error {
	existingHooks, err := asgService.DescribeLifecycleHooks(asgName)
	if err != nil {
		return err
	}

	for i := range wantedLifecycleHooks {
		if ignoreLifecycleHooks[wantedLifecycleHooks[i].Name] {
			log.Info("Not reconciling lifecycle hook since it's on the ignore list")
			continue
		}

		if err := reconcileLifecycleHook(ctx, asgService, asgName, &wantedLifecycleHooks[i], existingHooks, storeConditionsOnObject, log); err != nil {
			return err
		}
	}

	for _, existingHook := range existingHooks {
		found := false
		if ignoreLifecycleHooks[existingHook.Name] {
			continue
		}
		for _, wantedHook := range wantedLifecycleHooks {
			if existingHook.Name == wantedHook.Name {
				found = true
				break
			}
		}
		if !found {
			log.Info("Deleting extraneous lifecycle hook", "hook", existingHook.Name)
			if err := asgService.DeleteLifecycleHook(ctx, asgName, existingHook); err != nil {
				conditions.MarkFalse(storeConditionsOnObject, expinfrav1.LifecycleHookReadyCondition, expinfrav1.LifecycleHookDeletionFailedReason, clusterv1.ConditionSeverityError, err.Error())
				return err
			}
		}
	}

	return nil
}

func lifecycleHookNeedsUpdate(existing *expinfrav1.AWSLifecycleHook, expected *expinfrav1.AWSLifecycleHook) bool {
	return ptr.Deref(existing.DefaultResult, expinfrav1.LifecycleHookDefaultResultAbandon) != ptr.Deref(expected.DefaultResult, expinfrav1.LifecycleHookDefaultResultAbandon) ||
		ptr.Deref(existing.HeartbeatTimeout, metav1.Duration{Duration: 3600 * time.Second}) != ptr.Deref(expected.HeartbeatTimeout, metav1.Duration{Duration: 3600 * time.Second}) ||
		existing.LifecycleTransition != expected.LifecycleTransition ||
		existing.NotificationTargetARN != expected.NotificationTargetARN ||
		existing.NotificationMetadata != expected.NotificationMetadata
}

func reconcileLifecycleHook(ctx context.Context, asgService services.ASGInterface, asgName string, wantedHook *expinfrav1.AWSLifecycleHook, existingHooks []*expinfrav1.AWSLifecycleHook, storeConditionsOnObject conditions.Setter, log logger.Wrapper) error {
	log = log.WithValues("hook", wantedHook.Name)

	log.Info("Checking for existing lifecycle hook")
	var existingHook *expinfrav1.AWSLifecycleHook
	for _, h := range existingHooks {
		if h.Name == wantedHook.Name {
			existingHook = h
			break
		}
	}

	if existingHook == nil {
		log.Info("Creating lifecycle hook")
		if err := asgService.CreateLifecycleHook(ctx, asgName, wantedHook); err != nil {
			conditions.MarkFalse(storeConditionsOnObject, expinfrav1.LifecycleHookReadyCondition, expinfrav1.LifecycleHookCreationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return err
		}
		return nil
	}

	if lifecycleHookNeedsUpdate(existingHook, wantedHook) {
		log.Info("Updating lifecycle hook")
		if err := asgService.UpdateLifecycleHook(ctx, asgName, wantedHook); err != nil {
			conditions.MarkFalse(storeConditionsOnObject, expinfrav1.LifecycleHookReadyCondition, expinfrav1.LifecycleHookUpdateFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return err
		}
	}

	conditions.MarkTrue(storeConditionsOnObject, expinfrav1.LifecycleHookReadyCondition)
	return nil
}
