#!/bin/bash

declare -A sizes
sizes=(
    [480]=_xs
    [736]=_s
    [980]=_m
    [1280]=_l
    [1690]=_xl
)

for i in $*; do
	id=${i%.jpg}
	mv $id.jpg ${id}_full.jpg
	for h in "${!sizes[@]}"; do
		convert ${id}_full.jpg -sampling-factor 4:2:0 -strip -resize x$h -interlace JPEG -colorspace sRGB -quality 85 ${id}${sizes[$h]}.jpg
	done
done
