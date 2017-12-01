# Proxi - A simple go service to ingest KML and find K nearest places to a given location

# How to Build
Project dependencies are listed in `glide.yaml`, use `glide up` or simply go get each one of them. (Glide pulls the deps under vendor directory).

No build scripts are checked in yet, simply do `go build -o proxi` 

# Running
Once built the web service can be launched simply by `./proxi`. 
The service by default uses in-memory repository (with a default limit of 1000 indexable KML placemarks) and listens to port 8000. All of these can be changed by editing `src/config.yaml` file.

The service provides 4 endpoints and they are explained in the homepage `http://localhost:8000`(default)

# Mongo as backend
This is a TODO. Simply implement the Repository interface for Mongo (mongostore.go). The type `GeoPoint` is not needed for Mongo implementation as Mongo supports GeoJSON by default. Look for comments in `mongostore.go`
