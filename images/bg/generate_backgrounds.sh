#!/bin/bash

declare -A sizes
sizes=(
    [480]=_xs
    [736]=_s
    [980]=_m
    [1280]=_l
    [1690]=_xl
)

for i in `ls`; do
	for h in "${!sizes[@]}"; do
		convert $i -sampling-factor 4:2:0 -strip -resize x$h -interlace JPEG -colorspace sRGB -quality 85 ${i}${sizes[$h]}
	done
done
