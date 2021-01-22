#!/bin/sh

n=0
until [ "$n" -ge 5 ]
do
   python main.py
   n=$((n+1))
   sleep 15
done