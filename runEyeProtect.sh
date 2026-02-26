#!/bin/bash

exeDir="`dirname \"$0\"`"
cd $exeDir
./daemon -i 10 -w 20 -r "閉上眼睛" -t  
