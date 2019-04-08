#!/usr/bin/env python3
# coding=utf8

import os
import sys
from pushover_open_client import Client

class Pushover:
    def __init__(self, config=None, ledhat=None):
        if config == None or config == '':
            return

        self._ledhat = ledhat
        self._lemon_config = config
        self._client = Client(self._lemon_config)
        if self._client.secret == None or self._client.secret == "":
            self._client.login()
        if self._client.deviceID == None or self._client.deviceID == "":
            self._client.registerDevice('lemon')
        self._client.writeConfig(self._lemon_config)
        self._client.getWebSocketMessages(self._pushoverClientCallback)

    def _pushoverClientCallback(self, messageList):
        try:
            if(messageList):
                for message in messageList:
                    print(" id={0}, umid={1}, title={2}, message={3}, app={4}, aid={5}, icon={6}, date={7}, priority={8}, sound={9}, url={10}, url_title={11}, acked={12}, receipt={13}, html={14}, ".format(
                    message.id,
                    message.umid,
                    message.title,
                    message.message,
                    message.app,
                    message.aid,
                    message.icon,
                    message.date,
                    message.priority,
                    message.sound,
                    message.url,
                    message.url_title,
                    message.acked,
                    message.receipt,
                    message.html
                    ))

                    # Only show messages with priority greater/equal 0
                    if message.priority >= 0:
                        the_text = ''

                        if message.app != None and message.app != '':
                            the_text = the_text + message.app + ': '
                        if message.title != None and message.title != '':
                            the_text = the_text + message.title
                        if message.message != None and message.message != '':
                            the_text = the_text + ' -> ' + message.message.replace('\n', ' ')

                        self._ledhat.icon('pushover-prio-' + str(message.priority))
                        self._ledhat.text(the_text)
                    if(message.priority >= 2):
                        print("Ack message")
                        self._client.acknowledgeEmergency(message.receipt)
                print("Deleting message")
                self._client.deleteMessages(messageList[-1].id)
        except:
            print("Exception:")
            print(sys.exc_info())
