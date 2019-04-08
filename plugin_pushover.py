#!/usr/bin/env python3
# coding=utf8

import os
from pushover_open_client import Client

class Pushover:
    def __init__(self, config=None, ledhat=None):
        if config == None or config == '':
            return

        self._ledhat = ledhat
        self._lemon_config = config
        self._client = Client(self._lemon_config)
        self._client.login()
        self._client.registerDevice('lemon')
        self._client.writeConfig(self._lemon_config)
        self._client.getWebSocketMessages(self._pushoverClientCallback)

    def _pushoverClientCallback(self, messageList):
        if(messageList):
            for message in messageList:
                print(message)
                print(message.id)
                print(message.uuid)
                print(message.title)
                print(message.message)
                print(message.app)
                print(message.aid)
                print(message.icon)
                print(message.data)
                print(message.priority)
                print(message.sound)
                print(message.url)
                print(message.url_title)
                print(message.acked)
                print(message.receipt)
                print(message.html)

                # Only show messages with priority greater/equal 0
                if message.priority >= 0:
                    self._ledhat.icon('pushover-prio-' + message.priority)
                    self._ledhat.text(message.app + ': ' + message.title + ' -> ' + message.message)

                if(message.priority >= 2):
                    pushover_client.acknowledgeEmergency(message.receipt)
            pushover_client.deleteMessages(messageList[-1].id)
