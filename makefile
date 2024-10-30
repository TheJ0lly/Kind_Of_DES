SRC := $(filter-out src/tests.c, $(wildcard src/*.c))
TEST_SRC := $(filter-out src/main.c, $(wildcard src/*.c))

help:
	@echo "Usage: make <compile/run>" 

compile: $(SRC)
	gcc -o bin/sdes $(SRC) -Wall -Wswitch -Werror -pedantic -std=c11

run: compile
	./bin/sdes

test: $(TEST_SRC)
	gcc -o bin/sdes_test $(TEST_SRC) -Wall -Wswitch -Werror -pedantic -std=c11
	./bin/sdes_test
	rm bin/sdes_test