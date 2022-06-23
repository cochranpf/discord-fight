package data

import (
	"errors"
	"math/rand"
)

var Fighters []string
var Powers []string
var Nerfs []string

func add(input string, list *[]string) {
	if *list == nil {
		*list = make([]string, 0)
	}
	*list = append(*list, input)
}

func AddFighter(fighter string) {
	add(fighter, &Fighters)
}

func AddPower(power string) {
	add(power, &Powers)
}

func AddNerf(nerf string) {
	add(nerf, &Nerfs)
}

func getRandom(list []string) (string, error) {
	if len(list) == 0 {
		return "", errors.New("no elements in list")
	}
	return list[rand.Intn(len(list))], nil
}

func GetFighter() (string, error) {
	return getRandom(Fighters)
}

func GetPower() (string, error) {
	return getRandom(Powers)
}

func GetNerf() (string, error) {
	return getRandom(Nerfs)
}

func Clear() {
	Fighters = make([]string, 0)
	Powers = make([]string, 0)
	Nerfs = make([]string, 0)
}

func Ready() bool {
	return (len(Fighters) >= 2) && (len(Powers) >= 2) && (len(Nerfs) >= 2)
}
