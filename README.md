nip - query [BBC iplayer](https://www.bbc.co.uk/iplayer) for selected tv categories, search for programmes by title, lookup programme info and download programmes using
[youtube-dl](https://github.com/rg3/youtube-dl).

- [INSTALLATION](#installation)
- [DESCRIPTION](#description)
- [LICENSE](#license)

# INSTALLATION

## Dependencies

You need to have [youtube-dl](https://github.com/rg3/youtube-dl) installed, and for now
the [Go Programming Language](https://golang.org/doc/install).

**Then run**

`go get -u github.com/mswift42/nip`


# DESCRIPTION

nip builds a database of the most common iplayer tv categories, and stores that 
as a json file to disk. You can search for programmes by title or category, print
a programmes synopsis, go to a programmes homepage, list related links for a programme
and download programmes using youtube-dl.

# LICENSE

nip is licensed under the [MIT License](https://github.com/mswift42/nip/blob/master/LICENSE).

