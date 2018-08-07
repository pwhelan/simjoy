all:
	gb build vjoy/vjoygenerate
	gb generate
	gb build

clean:
	rm -f src/vjoy/const_linux.go
	rm bin/*

cross:
	./cross.sh
