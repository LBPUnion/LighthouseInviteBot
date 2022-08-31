package lighthouseapi

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/LBPUnion/LighthouseInviteBot/common"
	"github.com/sirupsen/logrus"
)

func GetInviteURL() string {

	request, err := http.NewRequest("POST", common.GetAPIURL()+"/user/inviteToken", nil)

	if err != nil {
		logrus.WithError(err).Errorln("Failed to create lighthouse request")
	}
	request.Header.Add("Authorization", "Basic "+common.LoadConfig().Lighthouse.APIKey)

	resp, er2 := http.DefaultClient.Do(request)
	if er2 != nil {
		logrus.WithError(er2).Errorln("Failed to do lighthouse request")
	}
	reader := bufio.NewScanner(resp.Body)
	reader.Scan()

	return common.LoadConfig().Lighthouse.ServerURL + "/register?token=" + strings.ReplaceAll(reader.Text(), "\"", "")
}

type LighthouseStatistics struct {
	RecentMatches int
	Slots         int
	Users         int
	TeamPicks     int
	Photos        int
}

func GetStatistics() (LighthouseStatistics, error) {
	resp, err := http.Get(common.GetAPIURL() + "/statistics")
	if err != nil {
		return LighthouseStatistics{}, errors.New("failed to get statistics")
	}
	bytes, er2 := io.ReadAll(resp.Body)
	if er2 != nil {
		return LighthouseStatistics{}, errors.New("failed to read HTTP response body")
	}

	var stats LighthouseStatistics

	er3 := json.Unmarshal(bytes, &stats)
	if er3 != nil {

		return LighthouseStatistics{}, errors.New("failed to unmartial JSON, maybe the server errored")
	}

	return stats, nil
}
