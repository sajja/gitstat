#!/bin/bash
cd /home/sajith/stats
cd $1
git log --author=$2 --since $3 --until $4 --oneline --shortstat | gawk '{ add += $1; subs += $2; loc += $1 - $2 } END { printf "added lines: %s removed lines: %s total lines: %s\n", add, subs, loc }' -