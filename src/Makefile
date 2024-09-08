# BASH settings per https://tech.davis-hansson.com/p/make/
SHELL := /bin/bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules


ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

NAME = lxt

SRC = $(shell find . -type f -name '*.go')
PWD = $(shell pwd)
TARGET = ../bin/$(NAME)

$(TARGET): $(SRC)
> /usr/local/go/bin/go build -tags=dev -o $(TARGET) cmd/$(NAME).go; 

clean:
> -@rm -rf ../bin/$(NAME)

dist: $(SRC) gen
> /usr/local/go/bin/go build -o $(TARGET) cmd/$(NAME).go; 

.PHONY: all clean gen test
