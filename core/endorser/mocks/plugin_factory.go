/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Code generated by mockery v1.0.0
package mocks

import endorsement "github.com/hyperledger/mchain/core/handlers/endorsement/api"
import mock "github.com/stretchr/testify/mock"

// PluginFactory is an autogenerated mock type for the PluginFactory type
type PluginFactory struct {
	mock.Mock
}

// New provides a mock function with given fields:
func (_m *PluginFactory) New() endorsement.Plugin {
	ret := _m.Called()

	var r0 endorsement.Plugin
	if rf, ok := ret.Get(0).(func() endorsement.Plugin); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(endorsement.Plugin)
		}
	}

	return r0
}
