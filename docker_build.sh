#!/bin/bash
docker build --rm -t dlza-manager-storage-handler:latest .
docker tag dlza-manager-storage-handler:latest registry.localhost:5001/dlza-manager-storage-handler
docker push registry.localhost:5001/dlza-manager-storage-handler
