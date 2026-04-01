// Package travel implements the travel domain pack for DecisionBox.
// It registers itself as "travel" via init() so services can select it
// based on the project's domain field.
//
// This pack provides:
//   - AI Discovery: analysis areas, prompts, and profile schemas for travel
//   - Categories: vacation_rental (MVP), hotel/experience/flights planned
//
// Usage:
//
//	import _ "github.com/decisionbox-io/decisionbox/domain-packs/travel/go"
//	// Then: domainpack.Get("travel")
package travel

import (
	"github.com/decisionbox-io/decisionbox/libs/go-common/domainpack"
)

func init() {
	domainpack.Register("travel", NewPack())
}

// TravelPack implements domainpack.Pack and domainpack.DiscoveryPack
// for the travel and hospitality domain.
type TravelPack struct{}

// NewPack creates a new travel domain pack.
func NewPack() *TravelPack {
	return &TravelPack{}
}

func (p *TravelPack) Name() string { return "travel" }
