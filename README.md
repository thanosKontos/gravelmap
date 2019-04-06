# Gravel map

WIP

# Installation guide

Example of building and running in Ubuntu x64:

```
cd $PROJECT_PATH/cmd/extractor
env GOOS=linux GOARCH=amd64 go build -o ~/gravelmap/extractor main.go
~/gravelmap/extractor /tmp/berlin-latest.osm.pbf
```

