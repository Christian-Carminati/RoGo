package main

import "fmt"

func loadFuncs() {
	moves = []Move{
		{
			Name:    "self-heal",
			Allowed: []int{classNameToId("Mage")},
			Desc:    "heals the caster",
			Move: func(caster *Character, chs *[]Character, queue *Queue) error {
				// caster heals himself
				if (*caster).Hp+(10*int((*caster).Lvl)) > int((*caster).MaxHp) {
					(*caster).Hp = int((*caster).MaxHp)
					return nil
				}
				(*caster).Hp += 10 * int((*caster).Lvl)
				return nil
			},
		},
		{
			Name:    "attack",
			Allowed: []int{classNameToId("Mage"), classNameToId("Ranger"), classNameToId("Warrior"), classNameToId("Rogue")},
			Desc:    "use your weapon to attack one enemy",
			Move: func(caster *Character, chs *[]Character, queue *Queue) error {
				var Damage = weapons[(*caster).Weapon].Damage
				var DamageType = weapons[(*caster).Weapon].DamageType
				// character uses his melee weapon to attack an enemy

				attacked := SingleSelector("who do you want to attack?\n", chs, struct{ caster *Character }{caster: caster}, func(enemy Character, input struct{ caster *Character }) bool {

					return enemy.Friendly != input.caster.Friendly && enemy.Hp > 0-int(enemy.MaxHp)
				})

				(*chs)[attacked].Hp -= int(calculateDamageProtection(&DamageType, &(*chs)[attacked]) * float64(Damage) * float64((*caster).Lvl))
				return nil
			},
		},
		{
			Name:    "fireball",
			Allowed: []int{classNameToId("Mage")},
			Desc:    "the mage casts a huge fireball, hitting all the enemies",
			Move: func(caster *Character, chs *[]Character, queue *Queue) error {

				var Damage = 5
				var DamageType = []int{DmgTypeId("Fire")}
				var max_targets = 3
				var title = "Who do you want to attack? (select multiple targets up to " + fmt.Sprint(max_targets) + " )\n"
				// fireball deals AOE damage, it also targets the dead

				attacked := multipleSelector(title, chs, max_targets, struct{ caster *Character }{caster: caster}, func(enemy Character, input struct{ caster *Character }) bool {
					return enemy.Friendly != input.caster.Friendly
				})

				for _, v := range attacked {

					(*chs)[v].Hp -= Damage * int(calculateDamageProtection(&DamageType, &(*chs)[v])) * int((*caster).Lvl)
				}
				return nil
			},
		},
		{
			Name:    "mind control",
			Allowed: []int{classNameToId("Mage")},
			Desc:    "the mage controls the mind of the enemy for ⌊lvl/2⌋+1 turns",
			Move: func(caster *Character, chs *[]Character, queue *Queue) error {

				if (*caster).Focus == true {
					return fmt.Errorf("caster does not have focus")
				}

				i := SingleSelector("who do you want to control?\n", chs, struct{ caster *Character }{caster: caster}, func(enemy Character, input struct{ caster *Character }) bool {
					_, okMindC := enemy.Status[1]
					_, okCMind := enemy.Status[2]
					return enemy.Friendly != caster.Friendly && enemy.Hp > 0-int(enemy.MaxHp) && !okMindC && !okCMind
				})

				(*chs)[i].Friendly = !(*chs)[i].Friendly
				(*chs)[i].Focus = true

				if (*chs)[i].Status == nil {
					(*chs)[i].Status = make(map[int]int)
				}
				if (*caster).Status == nil {
					(*caster).Status = make(map[int]int)
				}
				(*chs)[i].Status[1] = int(caster.Lvl/2) + 1
				(*caster).Status[2] = int((*chs)[i].Id)
				(*caster).Focus = true

				return nil
			},
		},
		{
			Name:    "poisonus dart",
			Allowed: []int{classNameToId("Rogue")},
			Desc:    "the attacker launches a poisoned dart, dealing 5 dmg and posioning the subject for 2*Lvl stacks",
			Move: func(caster *Character, chs *[]Character, queue *Queue) error {

				var DamageArrow = weapons[(*caster).Weapon].Damage
				var DamageTypeArrow = weapons[(*caster).Weapon].DamageType
				var stack = 5

				i := SingleSelector("who do you want to attack?\n", chs, struct{ caster *Character }{caster: caster}, func(enemy Character, input struct{ caster *Character }) bool {

					return enemy.Friendly != input.caster.Friendly && enemy.Hp > 0-int(enemy.MaxHp)
				})

				(*chs)[i].Hp -= int(calculateDamageProtection(&DamageTypeArrow, &(*chs)[i]) * float64(DamageArrow) * float64((*caster).Lvl))

				if (*chs)[i].Status == nil {
					(*chs)[i].Status = make(map[int]int)
				}
				(*chs)[i].Status[0] = stack * int(caster.Lvl)

				return nil
			},
		},
	}
	statusEffects = []StatusEffect{
		{
			name: "poison",
			desc: "the character is poisoned, taking damage every turn",
			effect: func(key int, caster *Character, chs *[]Character, queue *Queue) error {

				if (*caster).Status[key] <= 0 {
					statusEffects[key].endEffect(key, caster, chs, queue)
				}

				(*caster).Hp -= (*caster).Status[key]

				(*caster).Status[key]--

				return nil
			},
			endEffect: func(key int, caster *Character, chs *[]Character, queue *Queue) {
				delete((*caster).Status, key)
			},
		},
		{
			name: "mind control",
			desc: "the character changes factions",
			effect: func(key int, caster *Character, chs *[]Character, queue *Queue) error {

				(*caster).Focus = true

				if (*caster).Status[key] <= 0 {
					statusEffects[key].endEffect(key, caster, chs, queue)
				}
				(*caster).Status[key]--

				return nil
			},
			endEffect: func(key int, caster *Character, chs *[]Character, queue *Queue) {
				(*caster).Friendly = !(*caster).Friendly
				(*caster).Focus = false

				for i := range *chs {
					if val, ok := (*chs)[i].Status[2]; ok && val == int((*caster).Id) {
						statusEffects[2].endEffect(2, &((*chs)[i]), chs, queue)
					}
				}

				delete((*caster).Status, key)
			},
		},
		{
			name: "controlling mind",
			desc: "caster is controlling the mind of another character",
			effect: func(key int, caster *Character, chs *[]Character, queue *Queue) error {
				return nil
			},
			endEffect: func(key int, caster *Character, chs *[]Character, queue *Queue) {
				(*caster).Focus = false
				delete((*caster).Status, key)
			},
		},
	}

}
