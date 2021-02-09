package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/avimitin/osuapi/internal/api"
	"github.com/avimitin/osuapi/internal/database"
)

type PlayerData interface {
	GetPlayerStat(name string) (string, error)
}

type OsuServer struct {
	Data PlayerData
}

func (osuSer *OsuServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/api/v1/players/")
	stat, err := osuSer.Data.GetPlayerStat(player)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}
	fmt.Fprint(w, stat)
}

type OsuPlayerData struct {
	LatestData *api.User
	LocalData  *database.User
}

type Player struct {
	Data *api.User  `json:"latest_data"`
	Diff *Different `json:"diff"`
}

type Different struct {
	PlayCount string `json:"play_count"`
	Rank      string `json:"rank"`
	PP        string `json:"pp"`
	Acc       string `json:"acc"`
	TotalPlay string `json:"total_play"`
}

func (opd *OsuPlayerData) GetPlayerStat(name string) (string, error) {
	p, e := getPlayerDataByName(name)
	if e != nil {
		return "", e
	}
	data, e := json.Marshal(p)
	if e != nil {
		return "", e
	}
	return string(data), nil
}

func getPlayerDataByName(name string) (*Player, error) {
	u, e := api.GetUsers(name)
	if e != nil {
		return nil, e
	}
	if len(u) <= 0 {
		return nil, errors.New("user %s not found")
	}
	p := &Player{
		Data: u[0],
		Diff: &Different{
			Acc: "1",
		},
	}
	return p, nil
}
