nip - query [BBC iplayer](https://www.bbc.co.uk/iplayer) for selected tv categories, search for programmes by title, lookup programme info and download programmes using
[youtube-dl](https://github.com/rg3/youtube-dl).

- [INSTALLATION](#installation)
- [DESCRIPTION](#description)
- [LICENSE](#license)
- [COMMANDS](#commands)

# INSTALLATION

## Dependencies

You need to have [youtube-dl](https://github.com/rg3/youtube-dl) installed, and for now
the [Go Programming Language](https://golang.org/doc/install).

**Then run**

`go get github.com/mswift42/nip`


# DESCRIPTION

nip builds a database of the most common iplayer tv categories, and stores that 
as a json file to disk. You can search for programmes by title or category, print
a programmes synopsis, go to a programmes homepage, list related links for a programme
and download programmes using youtube-dl.

# LICENSE

nip is licensed under the [MIT License](https://github.com/mswift42/nip/blob/master/LICENSE).

# COMMANDS

-   list, l              list all available categories.
-   category, c          list all programmes for a category.
-   search, s            search for a programme.
-   show, sh             open Programmes homepage.
-   synopsis, syn        print programme's synopsis
-   links, lnk           show related links for a programme
-   download, g, d, get  use youtube-dl to download programme with index n
-   formats, f           list youtube-dl formats for programme with index n
-   refresh, r           refresh programme db




