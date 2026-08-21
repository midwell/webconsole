// SPDX-FileCopyrightText: 2026 Forsway Scandinavia AB
// SPDX-License-Identifier: Apache-2.0

package configapi

import (
	"testing"

	"github.com/omec-project/webconsole/configmodels"
)

// Guaranteed rates are configured in the rule's bitrate-unit, like the maximum rates, and must be
// normalised to bps on the same path. Left unconverted, a rule written in Mbps reaches the PCF a
// million times too small.
func TestNormalizeConvertsGuaranteedBitRatesToBps(t *testing.T) {
	slice := &configmodels.Slice{
		ApplicationFilteringRules: []configmodels.SliceApplicationFilteringRules{{
			RuleName:       "gbr-rule",
			BitrateUnit:    "Mbps",
			AppMbrUplink:   50,
			AppMbrDownlink: 50,
			AppGbrUplink:   10,
			AppGbrDownlink: 20,
		}},
	}

	normalizeApplicationFilteringRules(slice)

	rule := slice.ApplicationFilteringRules[0]
	if rule.AppGbrUplink != 10_000_000 {
		t.Errorf("AppGbrUplink = %d, want 10000000 after normalising 10 Mbps", rule.AppGbrUplink)
	}
	if rule.AppGbrDownlink != 20_000_000 {
		t.Errorf("AppGbrDownlink = %d, want 20000000 after normalising 20 Mbps", rule.AppGbrDownlink)
	}
	if rule.AppMbrUplink != 50_000_000 {
		t.Errorf("AppMbrUplink = %d, want the maximum rates still normalised", rule.AppMbrUplink)
	}
}
