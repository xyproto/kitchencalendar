# Kitchen Calendar

Generate per-week calendars that can be printed out and hung up in the kitchen.

For instance, print out calendars for the next 10 week and attach them to the cupboard doors with Blu Tack or another adhesive.

Currently, the calendars are only in Norwegian, but pull requests are welcome.

### Example PDF

![kitchen calendar](img/kitchencalendar_februar_2023.png)

### Example use for Go >= 1.17

    go install github.com/xyproto/kitchencalendar@latest
    kitchencalendar -names Bob,Alice,Mallory,Judy -year 2023 -week 8

### General info

* Version: 0.0.1
* License: BSD-3
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
