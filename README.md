# Gmail to gCal
Take a email of shifts from riteq and put them into google calandar

### Sample
>Your schedule for 24 Dec 2018 through to 6 Jan 2019 is shown below
>
>| Day | Date   | Start Work | End Work | Total Hours | Breaks | Pay   | Org Level                           |
>|-----|--------|------------|----------|-------------|--------|-------|-------------------------------------|
>| Mon | 24 Dec | 12:45      | 18:15    | 05:30       | 00:30  | 05:00 | AAAA\Dry Operations\Snr CSO         |
>| Thu | 27 Dec | 13:45      | 20:00    | 06:15       | 00:30  | 05:45 | AAAA\Dry Operations\Snr CSO         |
>| Fri | 28 Dec | 13:45      | 20:00    | 06:15       | 00:30  | 05:45 | AAAA\Dry Operations\Snr CSO         |
>| Sat | 29 Dec | 06:00      | 12:15    | 06:15       | 00:00  | 06:15 | AAAA\Dry Operations\Dry Ops Officer |
>| Sun | 30 Dec | 13:00      | 18:15    | 05:15       | 00:00  | 05:15 | AAAA\Dry Operations\Snr CSO         |
>| Thu | 03 Jan | 13:15      | 21:15    | 08:00       | 00:00  | 08:00 | AAAA\Dry Operations\Dry Ops Officer |
>| Fri | 04 Jan | 07:30      | 14:00    | 06:30       | 00:30  | 06:00 | AAAA\Dry Operations\Snr CSO         |

### v0.1
    ~~~decode email~~~

### v0.2
    ~~~check first line if the year dates are the same or different~~~
    take the date from the date column and the year from the first line to create a full date
    place all vaules into a map[]?

### v0.3
    use map[]? values into sqlite
    send each event to gcal
    enter successful to sqlite

### v04
    check for failures
    Notify of failures