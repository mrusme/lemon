#!/usr/bin/env python3
# coding=utf8

import os
import sys
import socket
import datetime
from influxdb import InfluxDBClient

class Influx:
    def __init__(self, host=None, port=8086, udp=False, database=None, username=None, password=None):
        self._client = None
        self._hostname = socket.gethostname()

        if host == None or host == "" or database == None or database == "":
            print("Not activating InfluxDB because no host/database was set!")
            return

        if udp == False:
            self._client = InfluxDBClient(host=host, port=port, database=database, username=username, password=password)
        else:
            self._client = InfluxDBClient(host=host, use_udp=True, udp_port=port, database=database, username=username, password=password)

    def write(self, resource, icon, category="undefined"):
        json_data = []

        if self._client == None:
            return False

        json_data.append({
            "measurement": "lemon-notifications",
            "tags": {
                "host": self._hostname,
                "resource": resource
            },
            "time": datetime.datetime.now(datetime.timezone.utc).astimezone().isoformat(),
            "fields": {
                "resource": resource,
                "icon": icon,
                "category": category
            }
        })

        return self._client.write_points(json_data)
