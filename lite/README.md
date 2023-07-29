
An all-in-one image for iWF server.

(Internally, it just runs Temporalite as backend)

IWF service: http://localhost:8801/
WebUI: http://localhost:8233/
## How to use
```shell
docker run -p 8801:8801 -p 7933:7933 -p 8233:8233 -e AUTO_FIX_WORKER_URL=host.docker.internal --add-host host.docker.internal:host-gateway -it iworkflowio/iwf-server-lite:latest
```

## How to build
Make sure you are at the root directory of this project (parent of current):
```shell
docker build . -t iworkflowio/iwf-server-lite:<yourTag> -f lite/Dockerfile
```

