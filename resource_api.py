#!/usr/bin/env python3
# coding=utf8

import falcon

class ResourceApi(object):
    def __init__(self, ledhat=None, influx=None):
        self._ledhat = ledhat
        self._influx = influx

    def on_post(self, req, resp):
        if req.context['request']:
            body = req.context['request']
            print(body)

            if 'icon' in body and 'text' in body:
                icon = body['icon']
                text = body['text']
                category = 'undefined'
                repeat = 3
                icon_cycle_time = 0.10
                text_cycle_time = 0.10
                text_font = None

                if 'category' in body:
                    category = body['category']

                if 'icon_repeat' in body:
                    repeat = int(body['icon_repeat'])

                if 'icon_cycle_time' in body:
                    icon_cycle_time = float(body['icon_cycle_time'])

                if 'text_font' in body:
                    text_font = body['text_font']

                self._ledhat.icon(icon, repeat=repeat, cycle_time=icon_cycle_time)
                self._ledhat.text(text, cycle_time=text_cycle_time, font=text_font)

                if self._influx != None:
                    self._influx.write(resource=req.path.replace('/','_'), icon=icon, category=category)

                resp.status = falcon.HTTP_204
            else:
                resp.status = falcon.HTTP_400
        else:
            resp.status = falcon.HTTP_500
