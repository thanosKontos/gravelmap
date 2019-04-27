**WIP**

# Gravel map

A routing engine to create off-road routes created as a composition of other services (osmium, postgis and pgrouting). With the expection to make it more autonomous in the near future.

## Prerequisites

You need to have the following tools installed in order to run the routing:

- Postgres (with postGIS and pgRouting extensions)
- Osmium
- osm2pgrouting

## Installation guide

Example of building and running in Ubuntu x64:

```bash
git clone git@github.com:thanosKontos/gravelmap.git
cd gravelmap
env GOOS=linux GOARCH=amd64 go build -o /tmp/gravelmap cmd/main.go && /tmp/gravelmap version
```

## Import data

Example of building and running in Ubuntu x64

```bash
env GOOS=linux GOARCH=amd64 go build -o ~/gravelmap/gravelmap cmd/main.go && ~/gravelmap/gravelmap add-data --input /path/to/some/osm/attiki.osm --database routing --tag-cost-config profiles/pgrouting/mt_bike.xml
```

## Drawer util

The util is used in order to create a test HTML page in order to help with manual testing the router. Database and data should be prepared for the following to work. - TBD - 

```bash
/path/to/project/cmd/map_drawer$ env GOOS=linux GOARCH=amd64 go build -o ~/gravelmap/map-drawer main.go && ~/gravelmap/map-drawer 38.0030367,23.8110783 37.9495728,23.8312455 > ~/test.html
```
