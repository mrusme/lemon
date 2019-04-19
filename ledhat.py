#!/usr/bin/env python3
# coding=utf8

import os
import time
import colorsys
import threading
import _thread
from sys import exit, argv

try:
    from PIL import Image, ImageDraw, ImageFont
except ImportError:
    exit('Lemon requires the pillow module\nInstall with: sudo pip3 install pillow')

try:
    import unicornhathd as unicorn
    print("unicorn hat hd detected")
except ImportError:
	from unicorn_hat_sim import unicornhathd as unicorn

class LedHat:
    def __init__(self):
        self._lock_ui = _thread.allocate_lock()

        unicorn.rotation(270)
        unicorn.brightness(1)
        self._unicorn_width, self._unicorn_height = unicorn.get_shape()

        self._default_font_file = 'fonts/Hack-Regular.ttf'
        self._default_font_size = 12
        self._default_font = ImageFont.truetype(self._default_font_file, self._default_font_size)

    def _animate_icon(self, image, repeat=3, cycle_time=0.10):
        if image == None or type(image) is str == False:
            print("Not a string:", image)
            return

        self._lock_ui.acquire()

        for i in range(0, repeat):
            # this is the original pimoroni function for drawing sprites
            for o_x in range(int(image.size[0] / self._unicorn_width)):
                for o_y in range(int(image.size[1] / self._unicorn_height)):
                    valid = False

                    for x in range(self._unicorn_width):
                        for y in range(self._unicorn_height):
                            pixel = image.getpixel(((o_x * self._unicorn_width) + y, (o_y * self._unicorn_height) + x))
                            r, g, b = int(pixel[0]), int(pixel[1]), int(pixel[2])
                            if r or g or b:
                                valid = True
                            unicorn.set_pixel((self._unicorn_height - y - 1), x, r, g, b)

                    if valid:
                        unicorn.show()
                        time.sleep(cycle_time)
        unicorn.off()

        self._lock_ui.release()

    def _animate_text(self, line, cycle_time=0.10, font=None):
        if line == None or type(line) is str == False:
            print("Not a string:", line)
            return

        self._lock_ui.acquire()

        text_font = self._default_font

        if font != None:
            # TODO: Check whether file exists
            custom_self._default_font_file = 'fonts/' + font + '.ttf'
            try:
                text_font = ImageFont.truetype(custom_self._default_font_file, self._default_font_size)
            except IOError:
                text_font = self._default_font

        text_width = self._unicorn_width
        text_height = 0
        text_x = self._unicorn_width
        text_y = 2

        w, h = text_font.getsize(line)
        text_width += w + self._unicorn_width
        text_height = max(text_height, h)

        text_width += self._unicorn_width + text_x + 1

        image = Image.new('RGB', (text_width, max(16, text_height)), (0, 0, 0))
        draw = ImageDraw.Draw(image)

        offset_left = 0

        draw.text((text_x + offset_left, text_y), line, font=self._default_font, fill=(255,255,255,255))
        offset_left += self._default_font.getsize(line)[0] + self._unicorn_width

        for scroll in range(text_width - self._unicorn_width):
            for x in range(self._unicorn_width):
                for y in range(self._unicorn_height):
                    pixel = image.getpixel((x + scroll, y))
                    r, g, b = [int(n) for n in pixel]
                    unicorn.set_pixel(self._unicorn_width - 1 - x, y, r, g, b)

            unicorn.show()
            time.sleep(cycle_time / 5)
        unicorn.off()

        self._lock_ui.release()

    def icon(self, name, repeat=3, cycle_time=0.10):
        # TODO: Check whether file exists
        img = Image.open('icons/' + name + '.png')

        self._thread_icon = threading.Thread(target=self._animate_icon, args=(img,repeat,cycle_time,))
        self._thread_icon.daemon = True
        self._thread_icon.start()
        time.sleep(.05)

    def text(self, line, cycle_time=0.10, font=None):
        self._thread_text = threading.Thread(target=self._animate_text, args=(line,cycle_time,font,))
        self._thread_text.daemon = True
        self._thread_text.start()
        time.sleep(.05)
