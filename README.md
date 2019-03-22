**Notice: This repository is deprecated-- [code](http://github.com/kenellorando/cadence) has been integrated with Cadence 4.0**

# Cadence Database Populator
The `cadence-dbp` "Cadence Database Populator" is a [Cadence Radio](http://cadenceradio.com/) auxillary tool for populating a database with music metadata used in searches and song requests.

Though metadata population was previously achieved using a webserver function, the newest iteration of the server is written with Python, which is particularly slow for file metadata and database tasks. This populator is written using Golang as necessitated by speed.

## Usage
```
go run populate.go MUSIC_DIR SERVER_DIR
# Example:
# go run populate.go /home/user/music/ /home/user/cadence/server/
```
where 
* `MUSIC_DIR` is the directory containing the music files to be inserted into the database.
* `SERVER_DIR` is the directory containing a valid Cadence Radio server configuration file

The cadence-dbp recursively searches a given music directory and examines all audio-type files for metadata including song title, artist, album, genre, year, and filesystem path. The database login credentials, hostname, port, table names, and table columns should be correctly set in the Cadence Radio server configuration file. The populator will use this data for execution of database insert statements.

**Please note:** The populator recursively searches folders, then files, as they are sorted on the filesystem. The process will **stop** when an audio file containing blank metadata is encountered, and an action must be taken on it (add metadata, remove the file) before the populator can examine any files that come after it.
