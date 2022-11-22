package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func serializer() {
	fmt.Println("testttt")

	armors = []Armor{
		{
			Name: "Old Rusty Chainmail",
			Resistences: map[int]float64{
				DmgTypeId("Slashing"): 0.6,
				DmgTypeId("Piercing"): 0.2,
			},
		},
		{
			Name: "Damaged Plate Armor",
			Resistences: map[int]float64{
				DmgTypeId("Piercing"):    0.6,
				DmgTypeId("Bludgeoning"): 0.4,
				DmgTypeId("Slashing"):    0.2,
			},
		},
	}
	weapons = []Weapon{
		{
			Name: "Longsword",
			DamageType: []int{
				DmgTypeId("Slashing"),
			},
			Damage: 8,
		},
		{
			Name: "Spear",
			DamageType: []int{
				DmgTypeId("Piercing"),
				DmgTypeId("Slashing"),
			},
			Damage: 20,
		},
		{
			Name: "Iron Mace",
			DamageType: []int{
				DmgTypeId("Bludgeoning"),
			},
			Damage: 15,
		},
		{
			Name: "Crossbow",
			DamageType: []int{
				DmgTypeId("Piercing"),
			},
			Damage: 10,
		},
	}
	if err := WriteJson("files/armors.json", &armors); err != nil {
		fmt.Printf("Error writing armors %e", err)
	}
	if err := WriteJson("files/weapons.json", &weapons); err != nil {
		fmt.Printf("Error writing weapons %e", err)
	}
	os.Exit(0)
}
func WriteJson[T any](FileName string, inp T) error {
	file, _ := json.MarshalIndent(inp, "", "\t")

	err := ioutil.WriteFile(FileName, file, 0644)
	return err
}
