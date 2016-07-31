cwc -> Cyclists With Cameras

For more information on CyclistsWithCameras see [here](https://github.com/jehiah/safe_streets/blob/master/cyclists_with_cameras.md)

This repository is a set of command line utilities for managing a database of complaints. It expects complaints to be one per directory in a directory named `YYYYMMDD_HHMM_$LICENSE`. The default location for these is `~/Documents/cyclists_with_cameras`.

A directory is a spot to combine your notes on the complaint (in a `notes.txt` file) and all the images or video related to that complaint. A typical notes.txt follows the structure

```
YYYY/MM/DD hh:mm:p [Taxi|FHV] $LICENSE $LOCATION

At <LOCATION> I observed <VEHICLE> <VIOLATION>. Pictures included.
```

Included are the following utilities:

* ncwc (new cyclists with cameras report)
* cwc (search tool)
* `reports.sh` generate statistics (i.e [this report](https://on.jehiah.cz/29J6lIX))

## cwc search

cwc is a flexible tool for searching for keywords in complaints. This is helpful when reviewing cases with the TLC as you can search for the TLC complaint number, the 311 number, the license plate, or anything that might be in the notes file.

```
$ cwc 4K45
Searching for: "4K45"
opening: 4K45 - Fri Jul 29 2016 5:48pm /Users/jehiah/Documents/cyclists_with_cameras/20160729_1748_4K45
also found: 4K45 - Fri Jul 8 2016 4:20pm /Users/jehiah/Documents/cyclists_with_cameras/20160708_1620_4K45
```

## nwcw complaint creation tool

`ncwc` is a tool to walk through generating a detailed consistent complaint. It helps you gather information, make sure you have a clear complaint and a clearly identified violation.
