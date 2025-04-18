package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type Player struct {
	HP int
}

type takeDamage struct {
	amount int
}

func NewPlayer(hp int) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			HP: hp,
		}
	}
}

func (p *Player) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		fmt.Println("Player started. (read state from db)")
	case actor.Stopped:
		fmt.Println("Player stopped. (write state to db)")
	case takeDamage:
		fmt.Println("player took damage: ", msg.amount)
	}
}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatalf("failed to create engine: %v", err)
	}

	pid := e.Spawn(NewPlayer(100), "player", actor.WithID("myuserid"))

	msg := takeDamage{amount: 999}
	for i := 0; i < 100; i++ {
		e.Send(pid, msg)
	}

	time.Sleep(time.Second * 2)

	//fmt.Println("process pid:", pid)

}
