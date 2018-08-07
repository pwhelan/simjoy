function init ()
	print("hello")
end

if starting == true then
	init()
	toggle = 0
end

vjoy = vjoys[1]
vjoy:SetButton(joystick.BTN_A, toggle)

if toggle == 0 then
	toggle = 1
	print("tock: ", tick)
else
	toggle = 0
	print("tick: ", tick)
end
