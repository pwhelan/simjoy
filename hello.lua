function init ()
	toggle = true
end

vjoy = vjoys[1]
curtick = 0

function tick(cur)
	if (curtick % 100) == 0 then
		print("CURTICK:", curtick)
	end
	curtick = cur
	if (cur % 100) == 0 then
		if toggle == true then
			toggle = false
			print("toggled off")
		else
			toggle = true
			print("toggled on")
		end
	end
end

function mouse(ev)
	print("Event: ", ev['event'], " button: ", ev['button'])
end

function midirecv(dev, msg)
	print("MIDI received: ", dev["id"], " -> ", dev["name"])
	print(msg.channel, 
		msg.status, 
		msg.data[1], 
		msg.data[2])

	if msg.data[1] == 48 then
		print("48!")
		if msg.status == 144 then
			vjoy:SetButton(joystick.BTN_A, 1)
		elseif msg.status == 128 then
			vjoy:SetButton(joystick.BTN_A, 0)
		end
	end

	if msg.status == 176 and msg.data[1] == 16 then
		vjoy:SetAxis(joystick.ABS_X, (63356 * (msg.data[2] / 127)) - 32768)
	end

end

function joystickrecv(ev)
	if ev['type'] == joystick.EVENT_BUTTON then
		print("BUTTON: " + ev['button'] + " = " + ev['value'])
	end
	if ev['type'] == joystick.EVENT_AXIS then
		print("AXIS: " + ev['axis'] + " = " + ev['value'])
	end
end
