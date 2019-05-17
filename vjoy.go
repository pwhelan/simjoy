package main

import (
	"github.com/pwhelan/simjoy/vjoy"
	lua "github.com/yuin/gopher-lua"
)

const vJoyName = "vjoy"

// Registers VJoy type
func registerVJoy(L *lua.LState) {
	mt := L.NewTypeMetatable(vJoyName)
	// methods
	L.SetField(mt,
		"__index",
		L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"SetButton": func(L *lua.LState) int {
				ud := L.CheckUserData(1)
				if vj, ok := ud.Value.(*vjoy.VJoy); !ok {
					L.ArgError(1, "vjoy expected")
				} else {
					if L.GetTop() == 3 {
						btn := L.CheckInt(2)
						state := L.CheckInt(3)
						vj.SetButton(btn, state)
					} else {
						L.ArgError(1, "two arguments expected")
					}
				}
				return 0
			},
			"SetAxis": func(L *lua.LState) int {
				ud := L.CheckUserData(1)
				if vj, ok := ud.Value.(*vjoy.VJoy); !ok {
					L.ArgError(1, "vjoy expected")
				} else {
					if L.GetTop() == 3 {
						axis := L.CheckInt(2)
						pos := L.CheckInt(3)
						vj.SetAxis(axis, pos)
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
	L.SetMetatable(ud, L.GetTypeMetatable(vJoyName))
	return ud
}
