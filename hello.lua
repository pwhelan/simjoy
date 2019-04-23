function init ()
	print("hello")
	toggel = false
end

vjoy = vjoys[1]

function tick(cur)
	-- print("tick: ", tostring(cur))
	if (cur % 100) == 0 then
		if toggle == true then
			toggle = false
			print("WAX ON")
		else
			toggle = true
			print("WAX OFF")
		end
		-- print("ON")
		-- midi[1]["out"]:Send(1, 144, 50, 100)
	elseif (cur % 50) == 0 then
		-- print("OFF")
		-- midi[1]["out"]:Send(1, 128, 50, 0)
	end
end

function midirecv(dev, msg)
	print("MIDI RECV'd: ", dev["name"])
	-- midi[2]["out"]:Send(14, 144, 36, 69)	
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

	print("midi: ", 
		msg.channel, 
		msg.status, 
		msg.data[1], 
		msg.data[2])
end
