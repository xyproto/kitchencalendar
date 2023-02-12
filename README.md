# Kitchen Calendar

Generate per-week calendars that can be printed out and hung up in the kitchen.

For instance, print out calendars for the next 10 weeks and attach them to the cupboard doors with Blu Tack or another adhesive.

It gives a very nice and hands-on overview of the coming weeks, and it's easy to sync and align what's happening with other people.

Currently, the calendars are only in Norwegian, but pull requests are welcome.

### Example PDF

![kitchen calendar](img/kitchencalendar_februar_2023.png)

### Getting started

Install the utility, for Go >= 1.17

    go install github.com/xyproto/kitchencalendar@latest

For creating a `calendar.pdf` file

    kitchencalendar -names Bob,Alice,Mallory,Judy -year 2023 -week 8

For generating calendars for week 7 to 17 (with 2 weeks on each PDF), for this year

    for x in $(seq 7 2 17); do kitchencalendar -names Bob,Alice,Mallory,Judy -week $x -o week$x.pdf; done

### General info

* Version: 0.0.1
* License: BSD-3
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
