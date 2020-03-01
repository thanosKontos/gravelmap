**WIP**

# Gravel map

Gravelmap is a routing engine made for off-road adventurers (hikers, mountain bikers, SUVs).

## Installation guide

Example of building and running in Ubuntu x64:

```bash
git clone git@github.com:thanosKontos/gravelmap.git
cd gravelmap
env GOOS=linux GOARCH=amd64 go build -o /tmp/gravelmap cmd/main.go && /tmp/gravelmap version
```

## Import OSM data

The command below will use osm2pgrouting in order to add ways to the DB. It reads only extracted OSM XML files at the moment.

```bash
env GOOS=linux GOARCH=amd64 go build -o /tmp/gravelmap cmd/main.go && /tmp/gravelmap import-routing-data --tag-cost-config profiles/pgrouting_mt_bike.xml --input /path/to/osm/greece_E21N37.osm
```

## Create web-server

This is not part of the actual toolkit. It is just an example of how you may use the data from the above commands.

Plus is a nice way for me to debug the result in a nice interface.

```bash
env GOOS=linux GOARCH=amd64 go build -o /tmp/gravelmap cmd/main.go && /tmp/gravelmap create-web-server
```

Open example_website.html to test routing

![](resources/example_website.png)

## Special thanks

To Ryan Carrier for the initial version of the dijkstra implementation (taken from here: https://github.com/RyanCarrier/dijkstra). Deleted the longest path implementation, will probably change a bit in the future to suit my needs better.