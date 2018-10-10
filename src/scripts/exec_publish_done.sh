#!/bin/bash

on_die ()
{
    pkill -KILL -P $$
}

rm -rf /tmp/hls/$1
rm -rf /tmp/records/$1*
