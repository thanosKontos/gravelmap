**WIP**

# Gravel map

A routing engine to create off-road routes created as a composition of other services (osmium, postgis and pgrouting). With the expection to make it more autonomous in the near future.

# Installation guide

Example of building and running in Ubuntu x64:

## Drawer util

The util is used in order to create a test HTML page in order to help with manual testing the router. Database and data should be prepared for the following to work. - TBD - 

```bash
/path/to/project/cmd/map_drawer$ env GOOS=linux GOARCH=amd64 go build -o ~/gravelmap/map-drawer main.go && ~/gravelmap/map-drawer 38.0030367,23.8110783 37.9495728,23.8312455 > ~/test.html
```
