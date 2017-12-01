package main

import (
	"fmt"
	"math/rand"
	"time"
)

type User struct {
	id   int
	char int
}

func (u User) enter() {
	lan_house.using += 1
	fmt.Printf("Usuário %v entrou na lan house\n", u.id)

	// User using machine
	rand_time := rand.Intn(120-15) + 15
	time.Sleep(time.Second * time.Duration(rand_time))
	fmt.Printf("Usuário %v usou a máquina por %v minutos\n", u.id, rand_time)
	u.exit()
}

func (u User) exit() {
	lan_house.using -= 1

	if lan_house.using == 0 && len(queue) == 0 {
		close(queue)
		waiting_clients = false
	}
}

type LanHouse struct {
	capacity int
	using    int
}

func (lh LanHouse) has_char() bool {
	return lh.capacity > lh.using
}

func (lh LanHouse) open() {
	last_id := 0

	for len(queue) < cap(queue) || waiting_clients {
		// Sets flag to false when first client came
		if waiting_clients = true; len(queue) > 0 {
			waiting_clients = false
		}

		// Every random time (0-5 sec) a new User enters to queue
		time.Sleep(time.Second * time.Duration(rand.Intn(5)))

		usr := User{last_id + 1, -1}
		queue <- usr

		fmt.Printf("Usuário %v entrou na fila\n", usr.id)
		last_id += 1

		if last_id >= cap(queue) {
			break
		}
	}
}

var waiting_clients = true
var queue chan User
var lan_house = LanHouse{8, 0}

func main() {
	queue = make(chan User, 26)
	rand.Seed(time.Now().Unix())

	go lan_house.open()

	// Waiting setup machines and first users came
	time.Sleep(time.Second * time.Duration(5))

	for len(queue) > 0 || waiting_clients {
		if lan_house.has_char() {
			next_usr := <-queue
			go next_usr.enter()

			if next_usr.id >= cap(queue) {
				break
			}

			if len(queue) < 1 {
				waiting_clients = true
			}
		}

		time.Sleep(time.Second * time.Duration(1))
	}

	for lan_house.using > 0 {
		time.Sleep(time.Second * time.Duration(1))
	}

	fmt.Println("Todos os clientes foram atendidos. :)")
}
