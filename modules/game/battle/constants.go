package battle

const (
	AttackingEnemy  = "enemies"
	AttackingPlayer = "player"

	ChainMailProtection = 4

	WeaponHand          = "Руки"
	WeaponSword         = "Меч-кладенец"
	WeaponLightning     = "Молнии"
	WeaponBallLightning = "Шаровые молнии"
)

var Weapon = map[string]string{
	"hand":           "Руки",
	"sword":          "Меч-кладенец",
	"lightning":      "Молнии",
	"ball lightning": "Шаровые молнии",
}
