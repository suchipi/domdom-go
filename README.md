# domdom-go
##### Background
[DomDomSoft Anime Downloader](http://domdomsoft.com/domdomsoft-anime-downloader-overview) is a Windows application that facilitates downloading anime.

The Windows application is a SOAP client and download manager. The owner of DomDomSoft runs a SOAP server that gives download links back to the client.

##### So what's this?
`domdom-go` is a command-line client to the same SOAP server. However, it is cross-platform (should work on Windows, OS X, Linux, *BSD, pretty much anywhere Go is available), and can be scripted and run headless.

### Usage

I'll just give some shell output here; should be pretty self-explanatory.

<pre>
suchipi@debian:~/anime# domdom-go 
NAME:
   domdom-go - Download Anime!

USAGE:
   domdom-go [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
  Suchipi Izumi - &lt;me@suchipi.com&gt;

COMMANDS:
   update, u            Update the anime list
   listepisodes, l      List available episodes for a series
   download, d          Download episode(s)
   help, h              Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --outputdir, -o "~/Downloads/Anime"          Directory to save into. Directories for series names will be created within this directory. [$DOMDOM_OUTPUTDIR]
   --key, -k                                    DomDomSoft Anime Downloader key. Without a key, there is a limit of 5 episodes downloaded per 24 hours. [$DOMDOM_KEY]
   --animelist, -l "~/.domdom_anime_list"       Location to save/load anime list [$DOMDOM_ANIMELIST]
   --help, -h                                   show help
   --version, -v                                print the version
   
suchipi@debian:~/anime# echo $DOMDOM_KEY
&lt;REDACTED&gt;
suchipi@debian:~/anime# echo $DOMDOM_OUTPUTDIR
.
suchipi@debian:~/anime# domdom-go listepisodes
Please specify a series title to list episodes of.
NAME:
   listepisodes - List available episodes for a series

USAGE:
   command listepisodes [command options] [arguments...]

OPTIONS:
   --title, -t  Title of series
   
suchipi@debian:~/anime# domdom-go listepisodes -t "Kantai Collection"
Episodes for Kantai Collection:
1: [HorribleSubs] Kantai Collection - 01 [720p].mkv (352MB)
2: [HorribleSubs] Kantai Collection - 02 [720p].mkv (352MB)
3: [HorribleSubs] Kantai Collection - 03 [720p].mkv (353MB)
4: [HorribleSubs] Kantai Collection - 04 [720p].mkv (352MB)
5: [HorribleSubs] Kantai Collection - 05 [720p].mkv (353MB)
6: [HorribleSubs] Kantai Collection - 06 [720p].mkv (353MB)
7: [HorribleSubs] Kantai Collection - 07 [720p].mkv (353MB)
8: [HorribleSubs] Kantai Collection - 08 [720p].mkv (353MB)
9: [HorribleSubs] Kantai Collection - 09 [720p].mkv (352MB)
suchipi@debian:~/anime# domdom-go download
Please specify a series title to download episodes from.
NAME:
   download - Download episode(s)

USAGE:
   command download [command options] [arguments...]

OPTIONS:
   --title, -t          Title of series to download
   --episode, -e        Episode filename to download
   --episodeid, -i      Episode id to download
   --all, -a            Download entire series
   --redownload, -r     Redownload existing files
   --keep-zip-files, -z Keep zip files instead of removing them after extraction
   
suchipi@debian:~/anime# domdom-go download -t "Kantai Collection" -i 4
Downloading file [HorribleSubs] Kantai Collection - 04 [720p].mkv from series Kantai Collection...
Part 1 of 2: [HorribleSubs] Kantai Collection - 04 [720p].00000.zip.part (314MB)
Part 2 of 2: [HorribleSubs] Kantai Collection - 04 [720p].00001.zip.part (39MB)
Unzipping file [HorribleSubs] Kantai Collection - 04 [720p].mkv from 2 parts...
Download complete!
suchipi@debian:~/anime# ls
Kantai Collection
suchipi@debian:~/anime# ls Kantai\ Collection/
[HorribleSubs] Kantai Collection - 04 [720p].mkv
suchipi@debian:~/anime# 
</pre>

### Building

Building is the same as (almost) any go project:
<pre>
sudo apt-get install golang-go //or equivalent for your platform
mkdir ~/go
export GOPATH=~/go
go get github.com/suchipi/domdom-go
go install github.com/suchipi/domdom-go
</pre>
After these steps, the `domdom-go` binary will reside in `$GOPATH/bin`.

### Notes

If you don't have a key for DomDomSoft Anime Downloader, you will only be able to download 5 episodes per 24 hours.

### TODO

* Progress bar for downloads
* Ability to override final save dir (in case you want to put eg "Sword Art Online II" in "Sword Art Online/Season 2" instead)
