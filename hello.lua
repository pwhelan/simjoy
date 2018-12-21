function init ()
	print("hello")
end

if starting == true then
	init()
	toggle = 0
	midi[1]["out"]:Send(14, 144, 36, 69)
end

vjoy = vjoys[1]

-- midi[1]["out"]:Send(14, 144, 36, 127)
-- midi[1]["out"]:Send(14, 144, 36, 0)

if midi[1]["in"] ~= nil then
	m = midi[1]["in"]
	if m.data[1] == 36 then
		if m.status == 144 then
			vjoy:SetButton(joystick.BTN_A, 1)
		elseif m.status == 128 then
			vjoy:SetButton(joystick.BTN_A, 0)
		end
	end

	if m.status == 176 and m.data[1] == 16 then
		vjoy:SetAxis(joystick.ABS_X, (63356 * (m.data[2] / 127)) - 32768)
	end

	print("midi: ", 
		m.channel, 
		m.status, 
		m.data[1], 
		m.data[2])
end
