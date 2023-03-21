package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Item struct {
	Name  string
	Items []string
}

type User struct {
	SelectedRoom string
	AvBackpack   bool
	Backpack     *Backpack
}
type Backpack struct {
	Items []string
}

type Room struct {
	Name       string
	Text       string
	Access     bool
	Navigation map[string]string
	Quest      func() bool
	Items      []*Item
}

type Game struct {
	User  *User
	Rooms map[string]*Room
}

var game = &Game{}

func main() {

	initGame()

	for {
		command := IdleCommandFromUser()
		handleCommand(command)
	}

}

func initGame() {

	game.User = &User{
		SelectedRoom: "кухня",
		AvBackpack:   false,
		Backpack: &Backpack{
			Items: []string{},
		},
	}

	game.Rooms = map[string]*Room{

		"кухня": {
			Name: "кухня",
			Text: "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор",
			Navigation: map[string]string{
				"коридор": "коридор",
			},
			Access: true,
		},

		"коридор": {
			Name: "коридор",
			Text: "ничего интересного. можно пройти - кухня, комната, улица",
			Navigation: map[string]string{
				"комната": "комната",
				"кухня":   "кухня",
				"улица":   "улица",
			},
			Access: true,
		},

		"комната": {
			Name:  "комната",
			Quest: GameRoomQuest,
			Text:  "ты в своей комнате. можно пройти - коридор",
			Items: []*Item{
				{
					Name: "стол",
					Items: []string{
						"ключи",
						"конспекты",
					},
				},
				{
					Name: "стул",
					Items: []string{
						"рюкзак",
					},
				},
			},
			Navigation: map[string]string{
				"коридор": "коридор",
			},
			Access: true,
		},

		"улица": {
			Name:  "улица",
			Quest: GameRoomQuest,
			Text:  "на улице весна. можно пройти - домой",
			Navigation: map[string]string{
				"комната": "комната",
			},
			Access: false,
		},
	}
}

func IdleCommandFromUser() string {

	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	return in.Text()
}

func handleCommand(command string) {

	SplitCommand := strings.Split(command, " ")

	if len(SplitCommand) < 3 {
		SplitCommand = append(SplitCommand, " ", " ")
	}

	action := SplitCommand[0]
	items := SplitCommand[1]
	addItems := SplitCommand[2]

	switch action {
	case "осмотреться":
		game.LookAround()
	case "идти":
		game.GoTo(items)
	case "надеть":
		game.PutOn(items)
	case "взять":
		game.AddBackpack(items)
	case "применить":
		game.ApplyItem(items, addItems)

	default:
		fmt.Println("неизвестная команда")
		return
	}

}

func (game *Game) LookAround() {

	room := game.User.SelectedRoom

	if !GameRoomQuest() {
		TextRoom := game.Rooms[room].Text
		fmt.Println(TextRoom)
	}
	return
}

func (game *Game) GoTo(room string) {

	GamerRoom := game.User.SelectedRoom

	admission, access := PermissionToVisit(game.Rooms[GamerRoom].Navigation, room)

	if admission && access {
		game.User.SelectedRoom = room
		fmt.Println(game.Rooms[room].Text)
		return
	}

	if admission && !access {
		fmt.Println("дверь закрыта")
		return
	}

	if !admission {
		fmt.Println("нет пути в", room)
		return
	}

}

func PermissionToVisit(items map[string]string, value string) (bool, bool) {

	admission := game.Rooms[value].Access

	for _, room := range items {
		if room == value {
			return true, admission
		}
	}
	return false, admission
}

func GameRoomQuest() bool {
	// На столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор
	GamerRoom := game.User.SelectedRoom
	if GamerRoom == "комната" {
		fmt.Println("На столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор")
		return true
	}
	return false
}

func foundItemFromSlice(items []string, FoundItems string) int {
	for i, val := range items {
		if val == FoundItems {
			return i
		}
	}
	return -1
}

func (game *Game) PutOn(item string) {

	if item == "рюкзак" {
		presence := game.User.AvBackpack
		if !presence {
			game.User.AvBackpack = true
			fmt.Println("вы надели :рюкзак")
		}
		return
	}
	fmt.Println("ошибка")
}

func (game *Game) AddBackpack(item string) {

	room := game.User.SelectedRoom
	itemsFromRoom := game.Rooms[room].Items[0].Items
	avBackpack := game.User.AvBackpack

	if avBackpack == false {
		fmt.Println("некуда класть")
		return
	}
	if foundItemFromSlice(itemsFromRoom, item) != -1 {
		game.User.Backpack.Items = append(game.User.Backpack.Items, item)
		DeleteItems(itemsFromRoom, item)
		fmt.Println("предмет добавлен в инвентарь:", item)
		return
	}
	fmt.Println("нет такого")
}

func DeleteItems(items []string, deleteItem string) {

	NumberItem := NumberItemsFromSlice(items, deleteItem)
	if NumberItem != -1 {
		items[NumberItem] = items[len(items)-1]
		items[len(items)-1] = ""
		items = items[:len(items)-1]
	}
}

func NumberItemsFromSlice(items []string, element string) int {
	for i, val := range items {
		if val == element {
			return i
		}
	}
	return -1
}

func FoundKeyFromBackpack(item string) bool {
	backpackItems := game.User.Backpack.Items
	for _, val := range backpackItems {
		if val == item {
			return true
		}
	}
	return false
}

func (game *Game) ApplyItem(item string, addItems string) {

	if FoundKeyFromBackpack(item) {
		if addItems == "дверь" {
			game.Rooms["улица"].Access = true
			fmt.Println("дверь открыта")
			return
		}
		fmt.Println("не к чему применить")
		return
	}
	fmt.Println("нет предмета в инвентаре -", item)
}
