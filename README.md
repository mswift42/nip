nip - query [BBC iplayer](https://www.bbc.co.uk/iplayer) for selected tv categories, search for programmes by title, lookup programme info and download programmes using
[youtube-dl](https://github.com/rg3/youtube-dl).

- [INSTALLATION](#installation)
- [DESCRIPTION](#description)
- [LICENSE](#license)
- [COMMANDS](#commands)
- [USAGE](#usage)

# INSTALLATION

## Dependencies

If you want to download iplayer programmes, You need to install [youtube-dl](https://github.com/rg3/youtube-dl).

If you have the Go programming language installed, you can run 

`go get -u github.com/mswift42/nip`

to install nip.

Else, you can grab a binary from the [releases](https://github.com/mswift42/nip/releases) section.


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
-   show, sh             open a programme's homepage.
-   synopsis, syn        print programme's synopsis
-   links, lnk           show related links for a programme
-   download, g, d, get  use youtube-dl to download programme with index n
-   formats, f           list youtube-dl formats for programme with index n
-   refresh, r           refresh programme db

# USAGE

enter `nip` followed by the command you want to run.

**Examples**

`nip l` will list all categories,

`nip c crime` will all programmes in category crime.

`nip --help` will show the help output.

`nip s pride` will list all programmes with "pride" in their title.

`nip g 133` will download programme with index 133 in the best available index.

`nip g 133 <youtube-dl format>` will download it with the specified youtube-dl format.
