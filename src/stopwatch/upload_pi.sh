#! /bin/bash

./build_pi.sh

echo kill stopwatch
ssh filou@timestamp "killall stopwatch wish"
echo copy stopwatch
scp stopwatch filou@timestamp:/home/filou
echo run
ssh filou@timestamp "DISPLAY=:0 nohup /home/filou/stopwatch"
echo upload done

#DISPLAY=:0 ./stopwatch
