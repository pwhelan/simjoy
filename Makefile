all:
	gb build vjoy/vjoygenerate
	gb build vkbd/vkbdgenerate
	gb generate
	gb build

clean:
	rm -f src/vjoy/const_linux.go
	rm -f src/vkbd/const_linux.go
	rm bin/*

cross:
	./cross.sh
