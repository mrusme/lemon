#!/bin/sh

if [ "$1" = "" ] || [ "$2" = "" ] || [ "$3" = "" ]
then
  printf "usage: %s <email> <password> <2fa>\n" "$0"
  exit 1
fi

rsp_login=$(curl -s \
  --form-string "email=$1" \
  --form-string "password=$2" \
  --form-string "twofa=$3" \
  https://api.pushover.net/1/users/login.json)

secret=$(printf "%s" "$rsp_login" | jq --raw-output '.secret')

rsp_devices=$(curl -s \
  --form-string "secret=$secret" \
  --form-string "name=lemon" \
  --form-string "os=O" \
  https://api.pushover.net/1/devices.json)

device_id=$(printf "%s" "$rsp_devices" | jq --raw-output '.id')

printf "export PUSHOVER_DEVICE_ID=\"%s\" PUSHOVER_SECRET=\"%s\"\n" \
  "$device_id" "$secret"

