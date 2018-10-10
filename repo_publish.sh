#!/bin/bash

set -e

curl --fail --user $REPO_USER:$REPO_PWD -v -XPOST -F file=@$PACKAGE_FILE $REPO_ADDR/api/files/$PACKAGE_NAME
if [ $? -ne 0 ]; then
    echo "failed to upload package deb file"
    exit 1
fi;

curl --fail --user $REPO_USER:$REPO_PWD -v -XPOST $REPO_ADDR/api/repos/liveplanet-internal/file/$PACKAGE_NAME
if [ $? -ne 0 ]; then
    echo "failed to add package to repo database"
    exit 1
fi;

curl --fail --user $REPO_USER:$REPO_PWD -v -XPUT -H 'Content-Type: application/json' --data '{}' $REPO_ADDR/api/publish/gs:liveplanet-internal:liveplanet-internal/jessie
if [ $? -ne 0 ]; then
    echo "failed to publish repo"
    exit 1
fi;
