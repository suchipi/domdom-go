package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/go-homedir"
	"github.com/suchipi/domdom-go/domdom"
	"github.com/suchipi/yukkuri/download"
	"github.com/suchipi/yukkuri/unzip"
)

func check(err error) {
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "domdom-go"
	app.Usage = "Download Anime!"
	app.Version = "0.1.1"
	app.Author = "Suchipi Izumi"
	app.Email = "me@suchipi.com"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "outputdir, o",
			Value:  "~/Downloads/Anime",
			Usage:  "Directory to save into. Directories for series names will be created within this directory.",
			EnvVar: "DOMDOM_OUTPUTDIR",
		},
		cli.StringFlag{
			Name:   "key, k",
			Value:  "",
			Usage:  "DomDomSoft Anime Downloader key. Without a key, there is a limit of 5 episodes downloaded per 24 hours.",
			EnvVar: "DOMDOM_KEY",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "search",
			ShortName: "s",
			Usage:     "Search for an anime by title",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "regex, r, term, t",
					Value: "",
					Usage: "Search term. It will be evaluated as a regular expression.",
				},
			},
			Action: searchAction,
		},
		{
			Name:      "listepisodes",
			ShortName: "l",
			Usage:     "List available episodes for a series",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "title, t",
					Value: "",
					Usage: "Title of series",
				},
			},
			Action: listepisodesAction,
		},
		{
			Name:      "download",
			ShortName: "d",
			Usage:     "Download episode(s)",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "title, t",
					Value: "",
					Usage: "Title of series to download",
				},
				cli.StringFlag{
					Name:  "episode, e",
					Value: "",
					Usage: "Episode filename to download",
				},
				cli.StringFlag{
					Name:  "episodeid, i",
					Value: "",
					Usage: "Episode id to download",
				},
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Download entire series",
				},
				cli.BoolFlag{
					Name:  "redownload, r",
					Usage: "Redownload existing files",
				},
				cli.BoolFlag{
					Name:  "keep-zip-files, z",
					Usage: "Keep zip files instead of removing them after extraction",
				},
			},
			Action: downloadAction,
		},
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Run(os.Args)
}

func updateAction(c *cli.Context) {
	savepath, err := homedir.Expand(c.GlobalString("animelist"))
	check(err)
	if savepath == "" {
		fmt.Println(`No save location for the anime list was specified. 
If you used the -l or --animelist flags, please check your input. 
Otherwise, check the contents of the $DOMDOM_ANIMELIST environment variable.
You can unset the environment variable to use the application default value, which is ~/.domdom_anime_list.`)
	}

	fmt.Println("Downloading anime list...")
	err = domdom.UpdateAnimeList(savepath)
	check(err)

	fmt.Printf("Anime list was downloaded and written to %s.\n", savepath)
}

func listepisodesAction(c *cli.Context) {
	title := c.String("title")
	key := c.GlobalString("key")

	if title == "" {
		fmt.Println("Please specify a series title to list episodes of.")
		cli.ShowCommandHelp(c, "listepisodes")
		return
	}

	episodes, err := domdom.ListEpisodes(title, key)
	check(err)

	fmt.Printf("Episodes for %s:\n", title)
	for index, episode := range episodes {
		sizeInt, err := strconv.Atoi(episode.FileSize)
		check(err)
		niceSize := humanize.Bytes(uint64(sizeInt))
		fmt.Printf("%d: %s (%s)\n", index+1, episode.FileName, niceSize)
	}
}

func downloadAction(c *cli.Context) {
	title := c.String("title")
	filename := c.String("episode")
	id := c.String("episodeid")
	key := c.GlobalString("key")
	redownload := c.Bool("redownload")
	keepzipfiles := c.Bool("keep-zip-files")
	all := c.Bool("all")

	if title == "" {
		fmt.Println("Please specify a series title to download episodes from.")
		cli.ShowCommandHelp(c, "download")
		return
	}
	if filename == "" && id == "" && all == false {
		fmt.Println("Please specify an episode to download.")
		cli.ShowCommandHelp(c, "download")
		return
	}
	if key == "" {
		fmt.Println("Warning: No key was specified. You will only be able to get download links 5 times per 24 hours.")
	}

	var episodes []domdom.Episode
	var err error

	if filename != "" {
		episode, err := domdom.FindEpisodeByName(title, filename, key)
		check(err)
		episodes = make([]domdom.Episode, 1)
		episodes[0] = episode
	} else if id != "" {
		episode, err := domdom.FindEpisodeByIndex(title, id, key)
		check(err)
		episodes = make([]domdom.Episode, 1)
		episodes[0] = episode
	} else if all {
		episodes, err = domdom.ListEpisodes(title, key)
		check(err)
	}

	outputdir, err := homedir.Expand(c.GlobalString("outputdir"))
	check(err)

	for _, episode := range episodes {
		downloadEpisode(episode, outputdir, key, redownload, keepzipfiles)
	}

	fmt.Println("Download complete!")
}

func downloadEpisode(episode domdom.Episode, outputdir, key string, redownload, keepzipfiles bool) error {
	links, err := domdom.GetDownloadLinks(episode, key)
	check(err)

	outputdir = filepath.Join(outputdir, episode.SeriesName)

	if _, err := os.Stat(filepath.Join(outputdir, episode.FileName)); err == nil {
		fmt.Printf("File %s already exists... ", episode.FileName)
		if redownload {
			fmt.Print("redownload requested.\n")
		} else {
			fmt.Print("not redownloading.\n")
			return nil
		}
	}

	zip_parts := make([]string, len(links))

	fmt.Printf("Downloading file %s from series %s...\n", episode.FileName, episode.SeriesName)
	for index, link := range links {
		dl, err := download.New(link, outputdir)
		check(err)

		finalPath := filepath.Join(dl.OutputDir, dl.FileName)
		shouldDownload := true

		fmt.Printf("Part %d of %d: %s (%s)\n", index+1, len(links), dl.FileName, humanize.Bytes(dl.FileSize))

		if _, err := os.Stat(finalPath); err == nil {
			fmt.Printf("File %s exists... ", dl.FileName)
			if redownload {
				fmt.Print("redownload requested.\n")
				shouldDownload = true
			} else {
				fmt.Print("not redownloading.\n")
				shouldDownload = false
			}
		}

		if shouldDownload {
			err = dl.Run()
			check(err)
		}
		zip_parts[index] = filepath.Join(outputdir, dl.FileName)
	}

	fmt.Printf("Unzipping file %s from %d parts...\n", episode.FileName, len(zip_parts))
	err = unzip.Multiple(zip_parts, outputdir)
	check(err)

	if keepzipfiles != true {
		for _, file := range zip_parts {
			os.Remove(file)
		}
	}
	return nil
}

func searchAction(c *cli.Context) {
	term := c.String("regex")
	if term == "" {
		fmt.Println("Please specify a search term.")
		cli.ShowCommandHelp(c, "search")
		return
	}
	regex, err := regexp.Compile(term)
	check(err)
	fmt.Println("Fetching anime list...")
	animes, err := domdom.GetAnimeList()
	check(err)
	output := make([]string, 0)
	for _, anime := range animes {
		if regex.MatchString(anime.Title) {
			output = append(output, fmt.Sprintf("%s (%s episodes)", anime.Title, anime.NumFiles))
		}
	}

	for _, line := range output {
		fmt.Println(line)
	}

	fmt.Printf("%d total result(s).\n", len(output))
}