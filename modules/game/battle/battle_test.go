package battle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	// Test constant values
	assert.Equal(t, "enemies", AttackingEnemy)
	assert.Equal(t, "player", AttackingPlayer)
	assert.Equal(t, 4, ChainMailProtection)
	assert.Equal(t, "Руками", WeaponHand)
	assert.Equal(t, "Меч-кладенец", WeaponSword)
	assert.Equal(t, "Молнии", WeaponLightning)
	assert.Equal(t, "Шаровые молнии", WeaponBallLightning)
}

func TestWeaponMap(t *testing.T) {
	// Test weapon map
	assert.Equal(t, "Руками", Weapon["hand"])
	assert.Equal(t, "Меч-кладенец", Weapon["sword"])
	assert.Equal(t, "Молнии", Weapon["lightning"])
	assert.Equal(t, "Шаровые молнии", Weapon["ball lightning"])

	// Test unknown weapon
	assert.Equal(t, "", Weapon["unknown"])
}
