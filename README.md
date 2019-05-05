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

## Use server

```bash
env GOOS=linux GOARCH=amd64 go build -o ~/gravelmap/gravelmap cmd/main.go && ~/gravelmap/gravelmap create-server
```

Open example_webasite.html to test routing

