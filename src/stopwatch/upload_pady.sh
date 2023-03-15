#! /bin/bash

./build_386.sh

echo kill stopwatch
ssh filou@pady "killall stopwatch wish"
echo copy stopwatch
scp stopwatch filou@pady:/home/filou
#echo run
#ssh filou@pady "DISPLAY=:0 nohup /home/filou/stopwatch"
echo upload done

#DISPLAY=:0 ./stopwatch
