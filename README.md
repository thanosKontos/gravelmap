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

## Import OSM data

The command below will use osm2pgrouting in order to add ways to the DB. It reads only extracted OSM XML files at the moment.

```bash
env GOOS=linux GOARCH=amd64 go build -o /tmp/gravelmap cmd/main.go && /tmp/gravelmap import-osm --tag-cost-config profiles/pgrouting_mt_bike.xml --input /path/to/osm/greece_E21N37.osm
```

## Import elevation data

The command below will import the elevation file into the database. It reads only ascii SRTM files at the moment.

```bash
env GOOS=linux GOARCH=amd64 go build -o /tmp/gravelmap cmd/main.go && /tmp/gravelmap import-elevation --input /path/to/asc/N37E021.asc
```

## Apply elevation cost to OSM data

If you have ran the import OSM and import elevation for a part of the earth, then you will need to grade the ways in terms of elevation

```bash
env GOOS=linux GOARCH=amd64 go build -o /tmp/gravelmap cmd/main.go && /tmp/gravelmap grade-ways
```

## Create web-server

This is not part of the actual toolkit. It is just an example of how you may use the data from the above commands.

Plus is a nice way for me to debug the result in a nice interface.

```bash
env GOOS=linux GOARCH=amd64 go build -o /tmp/gravelmap cmd/main.go && /tmp/gravelmap create-server
```

Open example_website.html to test routing

![](resources/example_website.png)
