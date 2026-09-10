// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package firecracker

import (
	"testing"

	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"
)

func TestMachineConfigurationVcpuLimit(t *testing.T) {
	t.Run("accepts 64 vCPUs", func(t *testing.T) {
		vcpuCount := int64(64)
		memory := int64(128)
		cfg := models.MachineConfiguration{
			MemSizeMib: &memory,
			VcpuCount:  &vcpuCount,
		}
		require.NoError(t, cfg.Validate(strfmt.Default))
	})

	t.Run("rejects more than 64 vCPUs", func(t *testing.T) {
		vcpuCount := int64(65)
		memory := int64(128)
		cfg := models.MachineConfiguration{
			MemSizeMib: &memory,
			VcpuCount:  &vcpuCount,
		}
		require.Error(t, cfg.Validate(strfmt.Default))
	})
}
