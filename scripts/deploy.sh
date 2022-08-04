#!/bin/sh
set -e

if [ -z "$VERSION" ]
then
    echo "version can not be empty."
    exit 1
fi;

echo "deploying new version $VERSION"

sed -i -r "s/:v[0-9]+\.[0-9]+\.[0-9]+/:v$VERSION/g" docker-compose.prod.yml
docker-compose -f docker-compose.prod.yml up