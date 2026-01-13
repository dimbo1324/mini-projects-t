package main

import (
	"fmt"
	"log"
	"slices"
	"time"
)

func main() {
	inputData := []string{
		"harry.dubois@mail.ru",
		"k.kitsuragi@mail.ru",
		"d.vader@mail.ru",
		"noname@mail.ru",
		"e.musk@mail.ru",
		"spiderman@mail.ru",
		"red_prince@mail.ru",
		"tomasangelo@mail.ru",
		"batman@mail.ru",
		"bruce.wayne@mail.ru",
	}

	result := []string{}
	RunPipeline(
		cmd(newCatStrings(inputData, 0)),
		cmd(SelectUsers),
		cmd(SelectMessages),
		cmd(CheckSpam),
		cmd(CombineResults),
		cmd(newCollectStrings(&result)),
	)

	expectedOutput := []string{
		"true 221945221381252775",
		"true 357347175551886490",
		"true 1595319133252549342",
		"true 1877225754447839300",
		"true 4652873815360231330",
		"true 5108368734614700369",
		"true 7829088386935944034",
		"true 8065084208075053255",
		"true 9323185346293974544",
		"true 10463884548348336960",
		"true 11204847394727393252",
		"true 12026159364158506481",
		"true 12386730660396758454",
		"true 12556782602004681106",
		"true 12728377754914798838",
		"true 13245035231559086127",
		"true 14107154567229229487",
		"true 16476037061321929257",
		"true 16728486308265447483",
		"true 17087986564527251681",
		"true 17259218828069106373",
		"true 17696166526272393238",
		"false 26236336874602209",
		"false 59892029605752939",
		"false 221962074543525747",
		"false 378045830174189628",
		"false 2803967521226628027",
		"false 6652443725402098015",
		"false 7594744397141820297",
		"false 9656111811170476016",
		"false 10167774218733491071",
		"false 10462184946173556768",
		"false 10493933060383355848",
		"false 10523043777071802347",
		"false 11512743696420569029",
		"false 12792092352287413255",
		"false 12975933273041759035",
		"false 14498495926778052146",
		"false 15161554273155698590",
		"false 15262116397886015961",
		"false 15728889559763622673",
		"false 15784986543485231004",
	}

	if !slices.Equal(expectedOutput, result) {
		log.Fatal("итоговый результат отличается от ожидаемого")
	} else {
		log.Println("MISSION PASSED! RESPECT +")
	}
}

func newCatStrings(strs []string, pauses time.Duration) func(in, out chan interface{}) {
	return func(in, out chan interface{}) {
		for _, email := range strs {
			out <- email
			if pauses != 0 {
				time.Sleep(pauses)
			}
		}
	}
}

func newCollectStrings(strs *[]string) func(in, out chan interface{}) {
	return func(in, out chan interface{}) {
		for dataRaw := range in {
			data := fmt.Sprintf("%v", dataRaw)
			*strs = append(*strs, data)
		}
	}
}
