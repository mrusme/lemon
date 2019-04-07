#!/bin/sh
BASEDIR=$(dirname "$0")

if [ -z "$PORT" ]; then
    export PORT=20001
fi

cd $BASEDIR && gunicorn -b "0.0.0.0:$PORT" lemon:app
