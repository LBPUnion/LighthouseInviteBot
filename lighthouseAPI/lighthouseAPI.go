package lighthouseapi

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Zaprit/LighthouseInviteBot/common"
)

func GetInviteURL() string {

	request, err := http.NewRequest("POST", common.LoadConfig().Lighthouse.ServerAPIURL+"/user/inviteToken", nil)

	if err != nil {
		log.Println(err.Error())
	}
	request.Header.Add("Authorization", "Basic "+common.LoadConfig().Lighthouse.APIKey)

	resp, er2 := http.DefaultClient.Do(request)
	if er2 != nil {
		log.Println(er2.Error())
	}
	reader := bufio.NewScanner(resp.Body)
	reader.Scan()

	fmt.Println(reader.Text())
	return common.LoadConfig().Lighthouse.ServerURL + "/register?token=" + strings.ReplaceAll(reader.Text(), "\"", "")
}
