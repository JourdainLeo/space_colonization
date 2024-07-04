package strategy

// EnumAlliance
type EnumAlliancEnumTYPE string
type EnumAlliancEnumVALUES struct {
	UNDEFINED EnumAlliancEnumTYPE
	NEVER     EnumAlliancEnumTYPE
	ALWAYS    EnumAlliancEnumTYPE
	IF_HI     EnumAlliancEnumTYPE
	IF_LO     EnumAlliancEnumTYPE
}

var EnumAlliance = EnumAlliancEnumVALUES{
	UNDEFINED: "undefined",
	NEVER:     "never",
	ALWAYS:    "always",
	IF_HI:     "if_hi",
	IF_LO:     "if_lo",
}

// EnumObservateur
type EnumObservationTYPE string
type EnumObservationVALUES struct {
	UNDEFINED   EnumObservationTYPE
	QUICK       EnumObservationTYPE
	SLOW        EnumObservationTYPE
	IF_HI_QUICK EnumObservationTYPE
	IF_LO_QUICK EnumObservationTYPE
	IF_HI_SLOW  EnumObservationTYPE
	IF_LO_SLOW  EnumObservationTYPE
}

var EnumObservation = EnumObservationVALUES{
	UNDEFINED:   "undefined",
	QUICK:       "quick",
	SLOW:        "slow",
	IF_HI_QUICK: "if_hi_quick",
	IF_LO_QUICK: "if_lo_quick",
	IF_HI_SLOW:  "if_hi_slow",
	IF_LO_SLOW:  "if_lo_slow",
}

// EnumLaunch
type EnumLaunchTYPE string
type EnumLaunchVALUES struct {
	UNDEFINED    EnumLaunchTYPE
	QUICK        EnumLaunchTYPE
	MEDIUM       EnumLaunchTYPE
	SLOW         EnumLaunchTYPE
	IF_HI_QUICK  EnumLaunchTYPE
	IF_LO_QUICK  EnumLaunchTYPE
	IF_HI_MEDIUM EnumLaunchTYPE
	IF_LO_MEDIUM EnumLaunchTYPE
	IF_HI_SLOW   EnumLaunchTYPE
	IF_LO_SLOW   EnumLaunchTYPE
}

var EnumLaunch = EnumLaunchVALUES{
	UNDEFINED:    "undefined",
	QUICK:        "quick",
	MEDIUM:       "medium",
	SLOW:         "slow",
	IF_HI_QUICK:  "if_hi_quick",
	IF_LO_QUICK:  "if_lo_quick",
	IF_HI_SLOW:   "if_hi_slow",
	IF_LO_SLOW:   "if_lo_slow",
	IF_LO_MEDIUM: "if_lo_medium",
	IF_HI_MEDIUM: "if_hi_medium",
}

// EnumSkipping
type EnumSkippingEnumTYPE string
type EnumSkippingEnumVALUES struct {
	UNDEFINED EnumAlliancEnumTYPE
	NEVER     EnumSkippingEnumTYPE
	ALWAYS    EnumSkippingEnumTYPE
	IF_HI     EnumSkippingEnumTYPE
	IF_LO     EnumSkippingEnumTYPE
}

var EnumSkipping = EnumSkippingEnumVALUES{
	UNDEFINED: "undefined",
	NEVER:     "never",
	ALWAYS:    "always",
	IF_HI:     "if_hi",
	IF_LO:     "if_lo",
} // EnumSkipping
type EnumReactionEnumTYPE string
type EnumReactionEnumVALUES struct {
	UNDEFINED   EnumReactionEnumTYPE
	FASTER      EnumReactionEnumTYPE
	SLOWER      EnumReactionEnumTYPE
	SAME        EnumReactionEnumTYPE
	BALANCED_LO EnumReactionEnumTYPE
	BALANCED_HI EnumReactionEnumTYPE
}

var EnumReaction = EnumReactionEnumVALUES{
	UNDEFINED:   "undefined",
	FASTER:      "faster",
	SLOWER:      "slower",
	SAME:        "same",
	BALANCED_LO: "balanced_lo",
	BALANCED_HI: "balanced_hi",
}
