cwc -> Cyclists With Cameras

For more information on CyclistsWithCameras see [here](https://github.com/jehiah/safe_streets/blob/master/cyclists_with_cameras.md)

This repository is a set of command line utilities for managing a database of complaints. It expects complaints to be one per directory in a directory named `YYYYMMDD_HHMM_$LICENSE`. The default location for these is `~/Documents/cyclists_with_cameras`.

A directory is a spot to combine your notes on the complaint (in a `notes.txt` file) and all the images or video related to that complaint. A typical notes.txt follows the structure

```
YYYY/MM/DD hh:mm:p [Taxi|FHV] $LICENSE $LOCATION

At <LOCATION> I observed <VEHICLE> <VIOLATION>. Pictures included.
  
C-1-1-123445678
```

Included are the following utilities:

* cwc - a multi-purpose search & report tool
* `reports.sh` generate statistics (i.e [this report](https://on.jehiah.cz/29J6lIX))

## cwc search

cwc is a flexible tool for searching for keywords in complaints. This is helpful when reviewing cases with the TLC as you can search for the TLC complaint number, the 311 number, the license plate, or anything that is present in the notes file.

```
$ cwc search 4K45
Searching for: "4K45"
opening: 4K45 - Fri Jul 29 2016 5:48pm /Users/jehiah/Documents/cyclists_with_cameras/20160729_1748_4K45
also found: 4K45 - Fri Jul 8 2016 4:20pm /Users/jehiah/Documents/cyclists_with_cameras/20160708_1620_4K45
```

## cwc new - complaint creation tool

`cwc new` is a tool to walk through generating a detailed consistent complaint. It helps you gather information, make sure you have a clear complaint and a clearly identified violation.

An example of how this works is below:

```
$ cwc new
Date (YYYYMMDD) or Filename: /Users/jehiah/Downloads/IMG_0421.JPG 

> using EXIF time 2016/08/31 6:00pm

License Plate: 8L56

Taxi [y/n] (Default is y): 

Where: West 21st St between 5th and 6th Ave

> creating /Users/jehiah/Documents/cyclists_with_cameras/20160831_1800_8L56

Violation: 

1. no driving in bike lane
2. no stopping in bike lane
3. no pickup or discharge of passengers in bike lane
4. no parking on sidewalks
5. blocking intersection and crosswalks
6. no u-turns in business district
7. no honking in non-danger situations
8. no driving in bus & right turn only lane
9. yield sign violation
10. failing to yield right of way
11. traffic signal violation
12. improper passing
13. unsafe lane change
14. no right from center lane
15. no left from center lane when both two-way streets
16. no left from center lane at one-way street
17. no passing zone
18. license plate must not be obstructed
19. no side window tint below 70%
20. cell-phone use while driving
21. threats, harassment, abuse
22. use or threat of physical force

Enter a number: 2

1. At <LOCATION> I observed <VEHICLE> <VIOLATION>. Pictures included.
2. <VEHICLE> stopped in bike lane, dangerously forcing bikers (including myself) into traffic lane <VIOLATION>. Pictures included.
3. <VEHICLE> stopped in bike lane, obstructing my use of bike lane <VIOLATION>. Pictures included.
4. While near <LOCATION> I observed <VEHICLE> stopped in bike lane <VIOLATION>. Pictures included.

Enter a number (Default is 1): 2

> opening https://www1.nyc.gov/apps/311universalintake/form.htm?serviceName=TLC+FHV+Driver+Unsafe+Driving
> done
```

## cwc report

`cwc report` provides a summarized view of activity based on complaint directories

## Building from Source

This project uses [getgb.io](https://getgb.io/). The following commands will build the binaries and place them in a `bin` directory. Use `vendor.sh` to load dependencies into a `vendor` directory.

```
go get github.com/constabulary/gb/...
./vendor.sh
gb build
```