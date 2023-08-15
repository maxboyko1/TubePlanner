# TubePlanner
Program for planning trips on London's commuter transit system. Transit data is taken from schedules available on the [Transport for London site](https://tfl.gov.uk/) as well as real-time trip planning applications. Full system map below.

![System Map](https://upload.wikimedia.org/wikipedia/commons/1/13/London_Underground_Overground_DLR_Crossrail_map.svg)

To use the program, build using `make` and run with two command-line arguments, specifying desired start and end locations for the journey. Surround multi-word station names in quotes. If both are valid locations, program will print a series of directions to complete the fastest possible trip between the two stations.

Terminal usage example below.

```
maxboyko:~/Documents/github/tubeplanner $ make
go build -o tubeplanner transitdata.go tubeplanner.go
maxboyko:~/Documents/github/tubeplanner $ ./tubeplanner
USAGE: ./tubeplanner <start> <destination>
maxboyko:~/Documents/github/tubeplanner $ ./tubeplanner Crikeyshire Hammersmith
ERROR: Crikeyshire is not a valid initial station
maxboyko:~/Documents/github/tubeplanner $ ./tubeplanner "Heathrow Terminal 4" Bonkersbury
ERROR: Bonkersbury is not a valid destination
maxboyko:~/Documents/github/tubeplanner $ ./tubeplanner Waterloo Waterloo
Already at destination!
maxboyko:~/Documents/github/tubeplanner $ ./tubeplanner "Queen's Park" "Canary Wharf"
1) Begin journey at Queen's Park station. (0 minutes)
2) Travel on the Bakerloo line, through station stops:
- Kilburn Park (1 minutes)
- Maida Vale (3 minutes)
- Warwick Avenue (5 minutes)
- Paddington (7 minutes)
3) Get off at Paddington and interchange to the Elizabeth line. (13 minutes)
4) Travel on the Elizabeth line, through station stops:
- Bond Street (16 minutes)
- Tottenham Court Road (19 minutes)
- Farringdon (22 minutes)
- Liverpool Street (25 minutes)
- Whitechapel (28 minutes)
- Canary Wharf (31 minutes)
5) Reach destination at Canary Wharf station. (31 minutes)
maxboyko:~/Documents/github/tubeplanner $ ./tubeplanner "High Street Kensington" "Canada Water"
1) Begin journey at High Street Kensington station. (0 minutes)
2) Travel on the Circle line, through station stops:
- Gloucester Road (2 minutes)
- South Kensington (6 minutes)
- Sloane Square (9 minutes)
- Victoria (11 minutes)
- St. James's Park (12 minutes)
- Westminster (14 minutes)
3) Get off at Westminster and interchange to the Jubilee line. (18 minutes)
4) Travel on the Jubilee line, through station stops:
- Waterloo (20 minutes)
- Southwark (21 minutes)
- London Bridge (23 minutes)
- Bermondsey (25 minutes)
- Canada Water (27 minutes)
5) Reach destination at Canada Water station. (27 minutes)
maxboyko:~/Documents/github/tubeplanner $ ./tubeplanner Uxbridge "Woolwich Arsenal"
1) Begin journey at Uxbridge station. (0 minutes)
2) Travel on the Metropolitan line, through station stops:
- Hillingdon (3 minutes)
- Ickenham (5 minutes)
- Ruislip (8 minutes)
- Ruislip Manor (9 minutes)
- Eastcote (11 minutes)
- Rayners Lane (15 minutes)
- West Harrow (17 minutes)
- Harrow-on-the-Hill (19 minutes)
- Northwick Park (22 minutes)
- Preston Road (24 minutes)
- Wembley Park (27 minutes)
- Finchley Road (34 minutes)
- Baker Street (40 minutes)
- Great Portland Street (42 minutes)
- Euston Square (44 minutes)
- King's Cross St. Pancras (46 minutes)
- Farringdon (49 minutes)
3) Get off at Farringdon and interchange to the Elizabeth line. (54 minutes)
4) Travel on the Elizabeth line, through station stops:
- Liverpool Street (57 minutes)
- Whitechapel (60 minutes)
- Canary Wharf (63 minutes)
- Custom House for ExCeL (67 minutes)
- Woolwich (71 minutes)
5) From Woolwich, interchange on foot to nearby Woolwich Arsenal station. (77 minutes)
6) Reach destination at Woolwich Arsenal station. (77 minutes)
```
