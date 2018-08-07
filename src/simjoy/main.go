package main

import (
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"

	"vjoy"
)

const luaVJOyTypeName = "person"

// Registers my person type to given L.
func registerVJoy(L *lua.LState) {
	mt := L.NewTypeMetatable(luaVJOyTypeName)
	// static attributes
	//L.SetField(mt, "new", L.NewFunction(newPerson))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"SetButton": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			if vj, ok := ud.Value.(*vjoy.VJoy); !ok {
				L.ArgError(1, "vjoy expected")
				return 0
			} else {
				if L.GetTop() == 3 {
					btn := L.CheckInt(2)
					state := L.CheckInt(3)
					vj.SetButton(btn, state)
					return 0
				} else {
					L.ArgError(1, "two arguments expected")
				}
			}
			return 0
		},
	}))
}

// Constructor
func userdataVJoy(L *lua.LState, vj *vjoy.VJoy) lua.LValue {
	ud := L.NewUserData()
	ud.Value = vj
	L.SetMetatable(ud, L.GetTypeMetatable(luaVJOyTypeName))
	return ud
}

func main() {
	vjoys := make([]*vjoy.VJoy, 1)
	vj, err := vjoy.OpenVJoy(0)
	if err != nil {
		panic(err)
	}

	vjoys[0] = vj

	quit := make(chan int)
	go func() {
		ft := 1000 * time.Millisecond / 60.0
		toggle := 0
		L := lua.NewState()
		defer L.Close()

		tick := 0

		registerVJoy(L)
		vjoy.Lua(L)
		vjoystable := L.NewTable()
		for _, vj := range vjoys {
			vjoystable.Append(userdataVJoy(L, vj))
		}

		L.SetGlobal("starting", lua.LBool(true))
		L.SetGlobal("vjoys", vjoystable)

		for {
			select {
			case <-quit:
				fmt.Println("quit")
				return
			default:
				L.SetGlobal("tick", lua.LNumber(float64(tick)))
				//vj.SetButton(vjoy.BTN_A, toggle)
				if toggle == 1 {
					toggle = 0
				} else {
					toggle = 1
				}
				if err := L.DoFile("hello.lua"); err != nil {
					panic(err)
				}
				vj.Tick()
				tick++
				if tick == 1 {
					L.SetGlobal("starting", lua.LBool(false))
				}
				time.Sleep(ft)
			}
		}
	}()

	fmt.Println("Press enter to exit")
	var input string
	fmt.Scanln(&input)
	fmt.Print(input)
	quit <- 0

	vj.Close()
}
