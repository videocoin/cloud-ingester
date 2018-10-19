# Ingester

Ingester is responsible for accepting streams from client (webcam) validaiting the stream and passing it along to Manager.

## Local Development

### Requirements

- Docker (tested on Docker version 18.06.1-ce, build e68fc7a)
- golang (tested on go version go1.10.4 linux/amd64)

### Running

Start nginx container

```shell

    docker build -t ingester -f Dockerfile.local .
    docker run  --network=host  -t ingester

```

Start hookd service

```shell

    cd hookd/cmd
    go run main.go

```

Create a test stream


```bash

    ffmpeg -f x11grab -s 720x480 -framerate 15 -i :0.0 -f pulse -ac 1 -i default -c:v libx264 -preset fast -pix_fmt yuv420p -s 1920x1080 -c:a aac -b:a 160k -ar 44100 -threads 0 -f flv "rtmp://0.0.0.0:1935/live/1-cameracamera"
    # "rtmp://0.0.0.0:1935/live/[userid int]-[camera id 16 charcter string]

```