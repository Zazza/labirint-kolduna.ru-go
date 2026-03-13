package bonus

func GetBonusNameByAlias(alias string) string {
	switch alias {
	case AntiPoisonSpellAlias:
		return AntiPoisonSpellName
	case InstantHypnosisSpellAlias:
		return InstantHypnosisSpellName
	case AntiPoisonSpellName:
		return AntiPoisonSpellName
	case DeathSpellAlias:
		return DeathSpellName
	case DeathTeleportAlias:
		return DeathTeleportName
	case InstantMovementAlias:
		return InstantMovementName
	case InstantRecoveryAlias:
		return InstantRecoveryName
	case LuckyStoneAlias:
		return LuckyStoneName
	case MagicDuckAlias:
		return MagicDuckName
	case WandAlias:
		return WandName
	}

	return ""
}
