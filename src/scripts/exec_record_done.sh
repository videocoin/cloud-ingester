#!/bin/bash

on_die ()
{
    pkill -KILL -P $$
}

ffmpeg -i $1 -vframes 1 -an -f image2 -s 640x480 -y /tmp/records/$2.jpg && \
rm -rf $1