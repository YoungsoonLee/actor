package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type Inventory struct {
	Bottles int
}

func NewInventory(bottles int) actor.Producer {
	return func() actor.Receiver {
		return &Inventory{
			Bottles: bottles,
		}
	}
}

func (i *Inventory) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		_ = msg
		fmt.Println("Inventory started. (read state from db)")
		c.Engine().Subscribe(c.PID()) // subscribe to event stream

	case actor.Stopped:
		fmt.Println("Inventory stopped. (write state to db)")
	case drinkBottle:
		fmt.Println("Inventory drank a bottle. (write state to db)")
		i.Bottles -= msg.amount
	case MyEvent:
		fmt.Println("Inventory received event stream", msg.foo)
	}
}

type Player struct {
	HP           int
	InventoryPID *actor.PID
}

type drinkBottle struct {
	amount int
}

func NewPlayer(hp int) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			HP: hp,
		}
	}
}

func (p *Player) Receive(c *actor.Context) {
	//switch msg := ctx.Message().(type) {
	switch msg := c.Message().(type) {
	case actor.Started:
		fmt.Println("Player started. (read state from db)")
		p.InventoryPID = c.SpawnChild(NewInventory(10), "inventory")

		c.Engine().Subscribe(c.PID()) // subscribe to event stream

	case actor.Stopped:
		fmt.Println("Player stopped. (write state to db)")
	case drinkBottle:
		//_ = msg
		//fmt.Println("player drinkBottle")
		//ctx.Send(p.InventoryPID, msg)
		c.Forward(p.InventoryPID)
	case string:
		fmt.Println("player received event stream")
	case MyEvent:
		fmt.Println("player received event stream", msg.foo)
	}
}

type MyEvent struct {
	foo string
}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatalf("failed to create engine: %v", err)
	}

	pid := e.Spawn(NewPlayer(100), "player", actor.WithID("myuserid"))

	e.Send(pid, drinkBottle{amount: 1})
	time.Sleep(time.Second * 2)

	//e.BroadcastEvent("event to all players")
	e.BroadcastEvent(MyEvent{foo: "bar"})
	time.Sleep(time.Second * 2)

	//<-e.Poison(pid).Done()

	// ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	// <-e.PoisonCtx(ctx, pid).Done()

	//_ = pid

	// msg := takeDamage{amount: 999}
	// for i := 0; i < 100; i++ {
	// 	e.Send(pid, msg)
	// }

	//fmt.Println("process pid:", pid)

}
