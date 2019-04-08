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

        self._ledhat = LedHat()
        self._ledhat.icon('lemon')

        self._app = falcon.API(middleware=[
            MiddlewareJson(),
        ])

        self._ledhat.text('Lemon')

        self._github = ResourceGitHub(ledhat=self._ledhat)
        self._api = ResourceApi(ledhat=self._ledhat)

        self._app.add_route('/api', self._api)
        self._app.add_route('/ifttt', self._api)
        self._app.add_route('/zapier', self._api)
        self._app.add_route('/github', self._github)

        self._pushover_client = None
        if PUSHOVER_CONFIG != None and PUSHOVER_CONFIG != '':
            from plugin_pushover import Pushover
            self._pushover_client = Pushover(config=PUSHOVER_CONFIG, ledhat=self._ledhat)

lemon = Lemon()
app = lemon._app
