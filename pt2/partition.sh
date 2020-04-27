#!/bin/bash

cd ./$1 || exit
split --number=4 $1.txt