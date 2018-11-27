// Copyright (c) 2018 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package descriptor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/ligato/cn-infra/logging"
	scheduler "github.com/ligato/vpp-agent/plugins/kvscheduler/api"
	"github.com/ligato/vpp-agent/plugins/vppv2/ipsecplugin/descriptor/adapter"
	"github.com/ligato/vpp-agent/plugins/vppv2/ipsecplugin/vppcalls"
	"github.com/ligato/vpp-agent/plugins/vppv2/model/ipsec"
)

const (
	// SADescriptorName is the name of the descriptor for VPP security associations.
	SADescriptorName = "vpp-ipsec-sa"
)

// A list of non-retriable errors:
var (
	// ErrSAWithoutIndex is returned when VPP security association was defined
	// without index.
	ErrSAWithoutIndex = errors.New("VPP security association defined without index")

	// ErrSAInvalidIndex is returned when VPP security association was defined
	// with non-numerical index.
	ErrSAInvalidIndex = errors.New("VPP security association defined with invalid index")
)

// IPSecSADescriptor teaches KVScheduler how to configure VPP IPSec security associations.
type IPSecSADescriptor struct {
	// dependencies
	log          logging.Logger
	ipSecHandler vppcalls.IPSecVppAPI
}

// NewIPSecSADescriptor creates a new instance of the IPSec SA descriptor.
func NewIPSecSADescriptor(ipSecHandler vppcalls.IPSecVppAPI, log logging.PluginLogger) *IPSecSADescriptor {
	return &IPSecSADescriptor{
		ipSecHandler: ipSecHandler,
		log:          log.NewLogger("ipsec-sa-descriptor"),
	}
}

// GetDescriptor returns descriptor suitable for registration (via adapter) with
// the KVScheduler.
func (d *IPSecSADescriptor) GetDescriptor() *adapter.SADescriptor {
	return &adapter.SADescriptor{
		Name:               SADescriptorName,
		KeySelector:        d.IsIPSecSAKey,
		ValueTypeName:      proto.MessageName(&ipsec.SecurityAssociation{}),
		KeyLabel:           d.IPSecSAIndexFromKey,
		ValueComparator:    d.EquivalentIPSecSAs,
		NBKeyPrefix:        ipsec.PrefixIPSecSA,
		Add:                d.Add,
		Delete:             d.Delete,
		ModifyWithRecreate: d.ModifyWithRecreate,
		IsRetriableFailure: d.IsRetriableFailure,
		Dump:               d.Dump,
	}
}

// IsIPSecSAKey returns true if the key is identifying VPP security association configuration.
func (d *IPSecSADescriptor) IsIPSecSAKey(key string) bool {
	return strings.HasPrefix(key, ipsec.PrefixIPSecSA)
}

// IPSecSAIndexFromKey returns VPP security association name from the key.
func (d *IPSecSADescriptor) IPSecSAIndexFromKey(key string) string {
	index, _ := ipsec.ParseSAIndexFromKey(key)
	return index
}

// EquivalentIPSecSAs is case-insensitive comparison function for
// ipsec.SecurityAssociation
func (d *IPSecSADescriptor) EquivalentIPSecSAs(key string, oldSA, newSA *ipsec.SecurityAssociation) bool {
	// compare base fields
	return oldSA.Spi == newSA.Spi &&
		oldSA.Protocol == newSA.Protocol &&
		oldSA.CryptoAlg == newSA.CryptoAlg &&
		oldSA.CryptoKey == newSA.CryptoKey &&
		oldSA.IntegAlg == newSA.IntegAlg &&
		oldSA.IntegKey == newSA.IntegKey &&
		oldSA.UseEsn == newSA.UseEsn &&
		oldSA.UseAntiReplay == newSA.UseAntiReplay &&
		oldSA.TunnelSrcAddr == newSA.TunnelSrcAddr &&
		oldSA.TunnelDstAddr == newSA.TunnelDstAddr &&
		oldSA.EnableUdpEncap == newSA.EnableUdpEncap
}

// IsRetriableFailure returns <false> for errors related to invalid configuration.
func (d *IPSecSADescriptor) IsRetriableFailure(err error) bool {
	nonRetriable := []error{
		ErrSAWithoutIndex,
		ErrSAInvalidIndex,
	}
	for _, nonRetriableErr := range nonRetriable {
		if err == nonRetriableErr {
			return false
		}
	}
	return true
}

// Add adds a new security association pair.
func (d *IPSecSADescriptor) Add(key string, sa *ipsec.SecurityAssociation) (metadata interface{}, err error) {
	// validate the configuration first
	err = d.validateSAConfig(sa)
	if err != nil {
		d.log.Error(err)
		return nil, err
	}

	// add security association
	err = d.ipSecHandler.AddSA(sa)
	if err != nil {
		d.log.Error(err)
	}

	return nil, err
}

// Delete removes VPP security association.
func (d *IPSecSADescriptor) Delete(key string, sa *ipsec.SecurityAssociation, metadata interface{}) error {
	err := d.ipSecHandler.DeleteSA(sa)
	if err != nil {
		d.log.Error(err)
	}
	return err
}

// ModifyWithRecreate always returns true - security associations are modified via re-creation.
func (d *IPSecSADescriptor) ModifyWithRecreate(key string, oldSA, newSA *ipsec.SecurityAssociation, metadata interface{}) bool {
	return true
}

// validateSAConfig validates VPP security association configuration.
func (d *IPSecSADescriptor) validateSAConfig(sa *ipsec.SecurityAssociation) error {
	if sa.Index == "" {
		return ErrSAWithoutIndex
	}
	if _, err := strconv.Atoi(sa.Index); err != nil {
		return ErrSAInvalidIndex
	}

	return nil
}

// Dump returns all configured VPP security associations.
func (d *IPSecSADescriptor) Dump(correlate []adapter.SAKVWithMetadata) (dump []adapter.SAKVWithMetadata, err error) {
	// dump security associations
	sas, err := d.ipSecHandler.DumpIPSecSA()
	if err != nil {
		d.log.Error(err)
		return dump, err
	}
	for _, sa := range sas {
		dump = append(dump, adapter.SAKVWithMetadata{
			Key:      ipsec.SAKey(sa.Sa.Index),
			Value:    sa.Sa,
			Metadata: sa.Meta,
			Origin:   scheduler.FromNB,
		})
	}

	var dumpList string
	for _, d := range dump {
		dumpList += fmt.Sprintf("\n - %+v", d)
	}
	d.log.Debugf("Dumping %d security associations: %v", len(dump), dumpList)

	return dump, nil
}