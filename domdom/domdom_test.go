package domdom

import (
	"os"
	"testing"
)

func TestGetAnimeListDoesNotError(t *testing.T) {
	_, err := GetAnimeList()
	if err != nil {
		t.FailNow()
	}
}
func TestGetAnimeListDoesNotReturnEmptySlice(t *testing.T) {
	animes, _ := GetAnimeList()
	if len(animes) == 0 {
		t.FailNow()
	}
}
func TestListEpisodesDoesNotError(t *testing.T) {
	key := os.Getenv("DOMDOM_KEY")
	//Yes, this assumes that Sword Art Online will always have episodes... I think that's a safe assumption, right?
	_, err := ListEpisodes("Sword Art Online", key)
	if err != nil {
		t.FailNow()
	}
}
func TestListEpisodesDoesNotReturnEmptySlice(t *testing.T) {
	key := os.Getenv("DOMDOM_KEY")
	//Yes, this assumes that Sword Art Online will always have episodes... I think that's a safe assumption, right?
	episodes, _ := ListEpisodes("Sword Art Online", key)
	if len(episodes) == 0 {
		t.FailNow()
	}
}
func TestUpdateAnimeListDoesNotError(t *testing.T) {
	filename := "test_update_anime_list"
	err := UpdateAnimeList(filename)
	if err != nil {
		t.FailNow()
	}
}
func TestUpdateAnimeListWritesFile(t *testing.T) {
	filename := "test_update_anime_list"
	UpdateAnimeList(filename)
	if _, err := os.Stat(filename); err != nil {
		t.FailNow()
	}
}
func TestLoadAnimeListDoesNotError(t *testing.T) {
	filename := "test_update_anime_list"
	_, err := LoadAnimeList(filename)
	if err != nil {
		t.FailNow()
	}
}
func TestLoadAnimeListDoesNotReturnEmptySlice(t *testing.T) {
	filename := "test_update_anime_list"
	animes, _ := LoadAnimeList(filename)
	if len(animes) == 0 {
		t.Fail()
	}
	os.Remove(filename)
}
func TestGetDownloadLinksDoesNotError(t *testing.T) {
	key := os.Getenv("DOMDOM_KEY")

	episode := Episode{ //Assumes this episode remains present serverside.
		SeriesName: "Sword Art Online",
		FileName:   "Sword_Art_Online_Special_01_Offline.mkv",
		FileSize:   "33115009",
	}

	_, err := GetDownloadLinks(episode, key)
	if err != nil {
		t.FailNow()
	}
}
func TestGetDownloadLinksDoesNotReturnEmptySlice(t *testing.T) {
	key := os.Getenv("DOMDOM_KEY")

	episode := Episode{ //Assumes this episode remains present serverside.
		SeriesName: "Sword Art Online",
		FileName:   "Sword_Art_Online_Special_01_Offline.mkv",
		FileSize:   "33115009",
	}

	links, _ := GetDownloadLinks(episode, key)
	if len(links) == 0 {
		t.FailNow()
	}
}
