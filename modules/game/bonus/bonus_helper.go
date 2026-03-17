package bonus

var bonusNameMap = map[string]string{
	AntiPoisonSpellAlias:      AntiPoisonSpellName,
	InstantHypnosisSpellAlias: InstantHypnosisSpellName,
	DeathSpellAlias:           DeathSpellName,
	DeathTeleportAlias:        DeathTeleportName,
	InstantMovementAlias:      InstantMovementName,
	InstantRecoveryAlias:      InstantRecoveryName,
	LuckyStoneAlias:           LuckyStoneName,
	MagicDuckAlias:            MagicDuckName,
	WandAlias:                 WandName,
	MagicRingAlias:            MagicRingName,
}

var aliasFromNameMap = map[string]string{
	AntiPoisonSpellName:      AntiPoisonSpellAlias,
	InstantHypnosisSpellName: InstantHypnosisSpellAlias,
	DeathSpellName:           DeathSpellAlias,
	DeathTeleportName:        DeathTeleportAlias,
	InstantMovementName:      InstantMovementAlias,
	InstantRecoveryName:      InstantRecoveryAlias,
	LuckyStoneName:           LuckyStoneAlias,
	MagicDuckName:            MagicDuckAlias,
	WandName:                 WandAlias,
	MagicRingName:            MagicRingAlias,
}

func GetBonusNameByAlias(alias string) string {
	if name, exists := bonusNameMap[alias]; exists {
		return name
	}
	return ""
}

func GetBonusAliasByName(name string) string {
	if alias, exists := aliasFromNameMap[name]; exists {
		return alias
	}
	return ""
}

func GetAllBonusAliases() []string {
	aliases := make([]string, 0, len(bonusNameMap))
	for alias := range bonusNameMap {
		aliases = append(aliases, alias)
	}
	return aliases
}

func GetAllBonusNames() []string {
	names := make([]string, 0, len(aliasFromNameMap))
	for name := range aliasFromNameMap {
		names = append(names, name)
	}
	return names
}
