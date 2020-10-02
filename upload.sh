#!/bin/bash

HEAD="v0.0."
git fetch
TAGS=$(git describe --tags --long)
MINOR=$(echo $TAGS | sed 's/-.*//g' | sed "s/^${HEAD}//g")
((MINOR=MINOR+1))
NEWTAG=${HEAD}${MINOR}
echo "Our new tag uploaded is = $NEWTAG"
make upload tag=$NEWTAG
