#!/bin/bash
cd /home/sajith/stats
cd $1
auth=$3
git log --author=$2 --since $3 --until $4 --numstat --pretty=tformat:  --no-merges| gawk '{ add += $1; subs += $2; loc += $1 - $2;name=$auth} END { 
printf "{\"name\":\"%s\",\n\"added\": %s,\n\"removed\": %s,\n\"total\": %s}\n", name,add, subs, loc }' 
