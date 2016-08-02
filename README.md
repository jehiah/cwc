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

An example of how this works is below:

```
$ ncwc
Date (YYYYMMDD) or Filename: **/Users/jehiah/Downloads/IMG_9024.JPG**
> using EXIF time 2016/08/01 5:21pm
License Plate: **T660390C**
Where? **8th Ave and West 25th St**
> creating /Users/jehiah/Documents/cyclists_with_cameras/20160801_1721_T660390C
1: no driving in bike lane
2: no stopping in bike lane
3: no pickup or discharge of passengers in bike lane
4: no parking on sidewalks
5: blocking intersection and crosswalks
6: no u-turns in business district
7: no honking in non-danger situations
8: yield sign violation
9: failing to yield right of way
10: traffic signal violation
11: improper passing
12: unsafe lane change
13: no right from center lane
14: no left from center lane when both two-way streets
15: no left from center lane at one-way street
16: no passing zone
17: license plate must not be obstructed
18: no side window tint below 70%
19: threats, harassment, abuse
20: use or threat of physical force
Violation (comma separate multiple): **10**
1: At <LOCATION> I observed <VEHICLE> <VIOLATION>. Pictures included.
2: At <LOCATION> I observed <VEHICLE> run red light <VIOLATION>. Pictures included. Pictures show light red and vehicle before intersection, and then vehicle proceeding through intersection on red.
Template: **2**
> opening https://www1.nyc.gov/apps/311universalintake/form.htm?serviceName=TLC+FHV+Driver+Unsafe+Driving
> done
```