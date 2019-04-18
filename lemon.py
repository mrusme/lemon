#!/usr/bin/env python3
# coding=utf8

import os
import falcon

from ledhat import LedHat

from middleware_json import MiddlewareJson

from resource_github import ResourceGitHub
from resource_api import ResourceApi

class Lemon:
    def __init__(self):
        PUSHOVER_CONFIG = os.getenv('LEMON_PUSHOVER_CONFIG', None)
        INFLUXDB_SERVER = os.getenv('LEMON_INFLUXDB_SERVER', None)
        INFLUXDB_PORT   = os.getenv('LEMON_INFLUXDB_PORT', 8086)
        INFLUXDB_UDP    = (os.getenv('LEMON_INFLUXDB_UDP', "0") == "1")
        INFLUXDB_DB     = os.getenv('LEMON_INFLUXDB_DB', None)
        INFLUXDB_USER   = os.getenv('LEMON_INFLUXDB_USER', None)
        INFLUXDB_PASS   = os.getenv('LEMON_INFLUXDB_PASS', None)

        self._ledhat = LedHat()
        self._ledhat.icon('lemon')

        self._app = falcon.API(middleware=[
            MiddlewareJson(),
        ])

        self._ledhat.text('Lemon')

        self._influx_client = None
        if INFLUXDB_SERVER != None and INFLUXDB_SERVER != '':
            from plugin_influx import Influx
            self._influx_client = Influx(host=INFLUXDB_SERVER, port=INFLUXDB_PORT, udp=INFLUXDB_UDP, database=INFLUXDB_DB, username=INFLUXDB_USER, password=INFLUXDB_PASS)

        self._github = ResourceGitHub(ledhat=self._ledhat, influx=self._influx_client)
        self._api = ResourceApi(ledhat=self._ledhat, influx=self._influx_client)

        self._app.add_route('/api', self._api)
        self._app.add_route('/ifttt', self._api)
        self._app.add_route('/zapier', self._api)
        self._app.add_route('/github', self._github)

        self._pushover_client = None
        if PUSHOVER_CONFIG != None and PUSHOVER_CONFIG != '':
            from plugin_pushover import Pushover
            self._pushover_client = Pushover(config=PUSHOVER_CONFIG, ledhat=self._ledhat, influx=self._influx_client)

lemon = Lemon()
app = lemon._app
