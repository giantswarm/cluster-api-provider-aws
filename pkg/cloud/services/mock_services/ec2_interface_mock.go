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

// Code generated by MockGen. DO NOT EDIT.
// Source: sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services (interfaces: EC2Interface)

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	ec2 "github.com/aws/aws-sdk-go/service/ec2"
	gomock "github.com/golang/mock/gomock"
	types "k8s.io/apimachinery/pkg/types"
	v1beta2 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	v1beta20 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	scope "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// MockEC2Interface is a mock of EC2Interface interface.
type MockEC2Interface struct {
	ctrl     *gomock.Controller
	recorder *MockEC2InterfaceMockRecorder
}

// MockEC2InterfaceMockRecorder is the mock recorder for MockEC2Interface.
type MockEC2InterfaceMockRecorder struct {
	mock *MockEC2Interface
}

// NewMockEC2Interface creates a new mock instance.
func NewMockEC2Interface(ctrl *gomock.Controller) *MockEC2Interface {
	mock := &MockEC2Interface{ctrl: ctrl}
	mock.recorder = &MockEC2InterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEC2Interface) EXPECT() *MockEC2InterfaceMockRecorder {
	return m.recorder
}

// CreateInstance mocks base method.
func (m *MockEC2Interface) CreateInstance(arg0 *scope.MachineScope, arg1 []byte, arg2 string) (*v1beta2.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateInstance", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1beta2.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateInstance indicates an expected call of CreateInstance.
func (mr *MockEC2InterfaceMockRecorder) CreateInstance(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInstance", reflect.TypeOf((*MockEC2Interface)(nil).CreateInstance), arg0, arg1, arg2)
}

// CreateLaunchTemplate mocks base method.
func (m *MockEC2Interface) CreateLaunchTemplate(arg0 scope.LaunchTemplateScope, arg1 *string, arg2 types.NamespacedName, arg3 []byte, arg4 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLaunchTemplate", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLaunchTemplate indicates an expected call of CreateLaunchTemplate.
func (mr *MockEC2InterfaceMockRecorder) CreateLaunchTemplate(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLaunchTemplate", reflect.TypeOf((*MockEC2Interface)(nil).CreateLaunchTemplate), arg0, arg1, arg2, arg3, arg4)
}

// CreateLaunchTemplateVersion mocks base method.
func (m *MockEC2Interface) CreateLaunchTemplateVersion(arg0 string, arg1 scope.LaunchTemplateScope, arg2 *string, arg3 types.NamespacedName, arg4 []byte, arg5 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLaunchTemplateVersion", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLaunchTemplateVersion indicates an expected call of CreateLaunchTemplateVersion.
func (mr *MockEC2InterfaceMockRecorder) CreateLaunchTemplateVersion(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLaunchTemplateVersion", reflect.TypeOf((*MockEC2Interface)(nil).CreateLaunchTemplateVersion), arg0, arg1, arg2, arg3, arg4, arg5)
}

// DeleteBastion mocks base method.
func (m *MockEC2Interface) DeleteBastion() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBastion")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBastion indicates an expected call of DeleteBastion.
func (mr *MockEC2InterfaceMockRecorder) DeleteBastion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBastion", reflect.TypeOf((*MockEC2Interface)(nil).DeleteBastion))
}

// DeleteLaunchTemplate mocks base method.
func (m *MockEC2Interface) DeleteLaunchTemplate(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLaunchTemplate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLaunchTemplate indicates an expected call of DeleteLaunchTemplate.
func (mr *MockEC2InterfaceMockRecorder) DeleteLaunchTemplate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLaunchTemplate", reflect.TypeOf((*MockEC2Interface)(nil).DeleteLaunchTemplate), arg0)
}

// DetachSecurityGroupsFromNetworkInterface mocks base method.
func (m *MockEC2Interface) DetachSecurityGroupsFromNetworkInterface(arg0 []string, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachSecurityGroupsFromNetworkInterface", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DetachSecurityGroupsFromNetworkInterface indicates an expected call of DetachSecurityGroupsFromNetworkInterface.
func (mr *MockEC2InterfaceMockRecorder) DetachSecurityGroupsFromNetworkInterface(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachSecurityGroupsFromNetworkInterface", reflect.TypeOf((*MockEC2Interface)(nil).DetachSecurityGroupsFromNetworkInterface), arg0, arg1)
}

// DiscoverLaunchTemplateAMI mocks base method.
func (m *MockEC2Interface) DiscoverLaunchTemplateAMI(arg0 scope.LaunchTemplateScope) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiscoverLaunchTemplateAMI", arg0)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DiscoverLaunchTemplateAMI indicates an expected call of DiscoverLaunchTemplateAMI.
func (mr *MockEC2InterfaceMockRecorder) DiscoverLaunchTemplateAMI(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiscoverLaunchTemplateAMI", reflect.TypeOf((*MockEC2Interface)(nil).DiscoverLaunchTemplateAMI), arg0)
}

// GetAdditionalSecurityGroupsIDs mocks base method.
func (m *MockEC2Interface) GetAdditionalSecurityGroupsIDs(arg0 []v1beta2.AWSResourceReference) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdditionalSecurityGroupsIDs", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdditionalSecurityGroupsIDs indicates an expected call of GetAdditionalSecurityGroupsIDs.
func (mr *MockEC2InterfaceMockRecorder) GetAdditionalSecurityGroupsIDs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdditionalSecurityGroupsIDs", reflect.TypeOf((*MockEC2Interface)(nil).GetAdditionalSecurityGroupsIDs), arg0)
}

// GetCoreSecurityGroups mocks base method.
func (m *MockEC2Interface) GetCoreSecurityGroups(arg0 *scope.MachineScope) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoreSecurityGroups", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoreSecurityGroups indicates an expected call of GetCoreSecurityGroups.
func (mr *MockEC2InterfaceMockRecorder) GetCoreSecurityGroups(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoreSecurityGroups", reflect.TypeOf((*MockEC2Interface)(nil).GetCoreSecurityGroups), arg0)
}

// GetInstanceSecurityGroups mocks base method.
func (m *MockEC2Interface) GetInstanceSecurityGroups(arg0 string) (map[string][]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstanceSecurityGroups", arg0)
	ret0, _ := ret[0].(map[string][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstanceSecurityGroups indicates an expected call of GetInstanceSecurityGroups.
func (mr *MockEC2InterfaceMockRecorder) GetInstanceSecurityGroups(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstanceSecurityGroups", reflect.TypeOf((*MockEC2Interface)(nil).GetInstanceSecurityGroups), arg0)
}

// GetLaunchTemplate mocks base method.
func (m *MockEC2Interface) GetLaunchTemplate(arg0 string) (*v1beta20.AWSLaunchTemplate, string, *types.NamespacedName, *string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLaunchTemplate", arg0)
	ret0, _ := ret[0].(*v1beta20.AWSLaunchTemplate)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(*types.NamespacedName)
	ret3, _ := ret[3].(*string)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// GetLaunchTemplate indicates an expected call of GetLaunchTemplate.
func (mr *MockEC2InterfaceMockRecorder) GetLaunchTemplate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLaunchTemplate", reflect.TypeOf((*MockEC2Interface)(nil).GetLaunchTemplate), arg0)
}

// GetLaunchTemplateID mocks base method.
func (m *MockEC2Interface) GetLaunchTemplateID(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLaunchTemplateID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLaunchTemplateID indicates an expected call of GetLaunchTemplateID.
func (mr *MockEC2InterfaceMockRecorder) GetLaunchTemplateID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLaunchTemplateID", reflect.TypeOf((*MockEC2Interface)(nil).GetLaunchTemplateID), arg0)
}

// GetLaunchTemplateLatestVersion mocks base method.
func (m *MockEC2Interface) GetLaunchTemplateLatestVersion(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLaunchTemplateLatestVersion", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLaunchTemplateLatestVersion indicates an expected call of GetLaunchTemplateLatestVersion.
func (mr *MockEC2InterfaceMockRecorder) GetLaunchTemplateLatestVersion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLaunchTemplateLatestVersion", reflect.TypeOf((*MockEC2Interface)(nil).GetLaunchTemplateLatestVersion), arg0)
}

// GetRunningInstanceByTags mocks base method.
func (m *MockEC2Interface) GetRunningInstanceByTags(arg0 *scope.MachineScope) (*v1beta2.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRunningInstanceByTags", arg0)
	ret0, _ := ret[0].(*v1beta2.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRunningInstanceByTags indicates an expected call of GetRunningInstanceByTags.
func (mr *MockEC2InterfaceMockRecorder) GetRunningInstanceByTags(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRunningInstanceByTags", reflect.TypeOf((*MockEC2Interface)(nil).GetRunningInstanceByTags), arg0)
}

// InstanceIfExists mocks base method.
func (m *MockEC2Interface) InstanceIfExists(arg0 *string) (*v1beta2.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceIfExists", arg0)
	ret0, _ := ret[0].(*v1beta2.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceIfExists indicates an expected call of InstanceIfExists.
func (mr *MockEC2InterfaceMockRecorder) InstanceIfExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceIfExists", reflect.TypeOf((*MockEC2Interface)(nil).InstanceIfExists), arg0)
}

// LaunchTemplateNeedsUpdate mocks base method.
func (m *MockEC2Interface) LaunchTemplateNeedsUpdate(arg0 scope.LaunchTemplateScope, arg1, arg2 *v1beta20.AWSLaunchTemplate) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LaunchTemplateNeedsUpdate", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LaunchTemplateNeedsUpdate indicates an expected call of LaunchTemplateNeedsUpdate.
func (mr *MockEC2InterfaceMockRecorder) LaunchTemplateNeedsUpdate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LaunchTemplateNeedsUpdate", reflect.TypeOf((*MockEC2Interface)(nil).LaunchTemplateNeedsUpdate), arg0, arg1, arg2)
}

// ModifyInstanceMetadataOptions mocks base method.
func (m *MockEC2Interface) ModifyInstanceMetadataOptions(arg0 string, arg1 *v1beta2.InstanceMetadataOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModifyInstanceMetadataOptions", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ModifyInstanceMetadataOptions indicates an expected call of ModifyInstanceMetadataOptions.
func (mr *MockEC2InterfaceMockRecorder) ModifyInstanceMetadataOptions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModifyInstanceMetadataOptions", reflect.TypeOf((*MockEC2Interface)(nil).ModifyInstanceMetadataOptions), arg0, arg1)
}

// PruneLaunchTemplateVersions mocks base method.
func (m *MockEC2Interface) PruneLaunchTemplateVersions(arg0 string) (*ec2.LaunchTemplateVersion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PruneLaunchTemplateVersions", arg0)
	ret0, _ := ret[0].(*ec2.LaunchTemplateVersion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PruneLaunchTemplateVersions indicates an expected call of PruneLaunchTemplateVersions.
func (mr *MockEC2InterfaceMockRecorder) PruneLaunchTemplateVersions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PruneLaunchTemplateVersions", reflect.TypeOf((*MockEC2Interface)(nil).PruneLaunchTemplateVersions), arg0)
}

// ReconcileBastion mocks base method.
func (m *MockEC2Interface) ReconcileBastion() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileBastion")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReconcileBastion indicates an expected call of ReconcileBastion.
func (mr *MockEC2InterfaceMockRecorder) ReconcileBastion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileBastion", reflect.TypeOf((*MockEC2Interface)(nil).ReconcileBastion))
}

// TerminateInstance mocks base method.
func (m *MockEC2Interface) TerminateInstance(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TerminateInstance", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// TerminateInstance indicates an expected call of TerminateInstance.
func (mr *MockEC2InterfaceMockRecorder) TerminateInstance(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TerminateInstance", reflect.TypeOf((*MockEC2Interface)(nil).TerminateInstance), arg0)
}

// TerminateInstanceAndWait mocks base method.
func (m *MockEC2Interface) TerminateInstanceAndWait(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TerminateInstanceAndWait", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// TerminateInstanceAndWait indicates an expected call of TerminateInstanceAndWait.
func (mr *MockEC2InterfaceMockRecorder) TerminateInstanceAndWait(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TerminateInstanceAndWait", reflect.TypeOf((*MockEC2Interface)(nil).TerminateInstanceAndWait), arg0)
}

// UpdateInstanceSecurityGroups mocks base method.
func (m *MockEC2Interface) UpdateInstanceSecurityGroups(arg0 string, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInstanceSecurityGroups", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInstanceSecurityGroups indicates an expected call of UpdateInstanceSecurityGroups.
func (mr *MockEC2InterfaceMockRecorder) UpdateInstanceSecurityGroups(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInstanceSecurityGroups", reflect.TypeOf((*MockEC2Interface)(nil).UpdateInstanceSecurityGroups), arg0, arg1)
}

// UpdateResourceTags mocks base method.
func (m *MockEC2Interface) UpdateResourceTags(arg0 *string, arg1, arg2 map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateResourceTags", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateResourceTags indicates an expected call of UpdateResourceTags.
func (mr *MockEC2InterfaceMockRecorder) UpdateResourceTags(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateResourceTags", reflect.TypeOf((*MockEC2Interface)(nil).UpdateResourceTags), arg0, arg1, arg2)
}
