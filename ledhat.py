#!/usr/bin/env python3

import os
import time
import colorsys
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

lock_ui = _thread.allocate_lock()

unicorn.rotation(270)
unicorn.brightness(1)

width, height = unicorn.get_shape()
cycle_time = 0.10

font_file = 'fonts/Hack-Regular.ttf'
font_size = 12
font = ImageFont.truetype(font_file, font_size)

def animate_icon(image, repeat=3):
    lock_ui.acquire()

    for i in range(0, repeat):
        # this is the original pimoroni function for drawing sprites
        for o_x in range(int(image.size[0] / width)):
            for o_y in range(int(image.size[1] / height)):
                valid = False

                for x in range(width):
                    for y in range(height):
                        pixel = image.getpixel(((o_x * width) + y, (o_y * height) + x))
                        r, g, b = int(pixel[0]), int(pixel[1]), int(pixel[2])
                        if r or g or b:
                            valid = True
                        unicorn.set_pixel((height - y - 1), x, r, g, b)

                if valid:
                    unicorn.show()
                    time.sleep(cycle_time)
    unicorn.off()

    lock_ui.release()

def animate_text(line):
    lock_ui.acquire()

    text_width = width
    text_height = 0
    text_x = width
    text_y = 2

    w, h = font.getsize(line)
    text_width += w + width
    text_height = max(text_height, h)

    text_width += width + text_x + 1

    image = Image.new('RGB', (text_width, max(16, text_height)), (0, 0, 0))
    draw = ImageDraw.Draw(image)

    offset_left = 0

    draw.text((text_x + offset_left, text_y), line, font=font, fill=(255,255,255,255))
    offset_left += font.getsize(line)[0] + width

    for scroll in range(text_width - width):
        for x in range(width):
            for y in range(height):
                pixel = image.getpixel((x + scroll, y))
                r, g, b = [int(n) for n in pixel]
                unicorn.set_pixel(width - 1 - x, y, r, g, b)

        unicorn.show()
        time.sleep(cycle_time / 5)
    unicorn.off()

    lock_ui.release()

def icon(name, repeat=3):
    img = Image.open('icons/' + name + '.png')

    _thread.start_new_thread(animate_icon, (img,repeat,))

def text(line):
    _thread.start_new_thread(animate_text, (line,))
