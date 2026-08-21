// SPDX-FileCopyrightText: 2026 Forsway Scandinavia AB
// SPDX-License-Identifier: Apache-2.0

package nfconfig

import (
	"testing"

	"github.com/omec-project/webconsole/configmodels"
)

func ruleWithRates(mbrUl, mbrDl, gbrUl, gbrDl int32) configmodels.SliceApplicationFilteringRules {
	return configmodels.SliceApplicationFilteringRules{
		RuleName:       "test-rule",
		AppMbrUplink:   mbrUl,
		AppMbrDownlink: mbrDl,
		AppGbrUplink:   gbrUl,
		AppGbrDownlink: gbrDl,
		TrafficClass:   &configmodels.TrafficClassInfo{Qci: 2, Arp: 1},
	}
}

// Rates reach this point already normalised to bps by normalizeApplicationFilteringRules, and
// ConvertToString picks the unit back by magnitude — so 10 here really is 10 bps.
func TestBuildPccQosCarriesGuaranteedBitRate(t *testing.T) {
	qos := buildPccQos(ruleWithRates(50, 50, 10, 20))

	if !qos.HasGbrUl() || !qos.HasGbrDl() {
		t.Fatal("a configured guaranteed rate must reach the policy served to the PCF")
	}
	if got := qos.GetGbrUl(); got != "10 bps" {
		t.Errorf("gbrUl = %q, want %q", got, "10 bps")
	}
	if got := qos.GetGbrDl(); got != "20 bps" {
		t.Errorf("gbrDl = %q, want %q", got, "20 bps")
	}
	if got := qos.GetMaxBrUl(); got != "50 bps" {
		t.Errorf("maxBrUl = %q, want the maximum rates unaffected", got)
	}
}

// A rule with no guaranteed rate must not acquire one. Non-GBR flows are the common case and a
// zero must not serialise as a guarantee of zero.
func TestBuildPccQosOmitsAnUnsetGuaranteedBitRate(t *testing.T) {
	qos := buildPccQos(ruleWithRates(50, 50, 0, 0))

	if qos.HasGbrUl() || qos.HasGbrDl() {
		t.Error("a rule with no guaranteed rate must not report one")
	}
}

// A guarantee in one direction only is plausible where the return link is the scarce one, so it
// must survive rather than being dropped for being incomplete.
func TestBuildPccQosCarriesAOneDirectionalGuarantee(t *testing.T) {
	qos := buildPccQos(ruleWithRates(50, 50, 10, 0))

	if !qos.HasGbrUl() {
		t.Error("an uplink guarantee must be carried")
	}
	if qos.HasGbrDl() {
		t.Error("an unset downlink guarantee must stay unset")
	}
}
