package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// 构造客户端
	client := http.Client{}
	// 构造Get请求
	URL := "https://api.bilibili.com/x/v2/reply/main?csrf=9a70a0c269d742c67c583eba29fa37b9&mode=3&next=0&oid=634501161&plat=1&type=1"
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println("req err", err)
	}
	// 添加请求头
	req.Header.Set("authority", "api.bilibili.com")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("referer", "https://www.bilibili.com/bangumi/play/ss39462?spm_id_from=333.337.0.0")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败", err)
	}
	defer resp.Body.Close()
	// 读取响应json到内存bodyTest中
	bodyTest, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io err", err)
	}
	var resultList KingRankResp
	_ = json.Unmarshal(bodyTest, &resultList)
	for _, commit := range resultList.Data.Replies {
		fmt.Println("一级评论", commit.Content.Message)
		for _, reply := range commit.Replies {
			fmt.Println("二级评论", reply.Content.Message)
		}
		fmt.Println("########################################################")
	}
}

type KingRankResp struct {
	Code int64 `json:"code"`
	Data struct {
		Replies []struct {
			Content struct {
				Emote struct {
					Doge struct {
						Attr      int64  `json:"attr"`
						ID        int64  `json:"id"`
						JumpTitle string `json:"jump_title"`
						Meta      struct {
							Size int64 `json:"size"`
						} `json:"meta"`
						Mtime     int64  `json:"mtime"`
						PackageID int64  `json:"package_id"`
						State     int64  `json:"state"`
						Text      string `json:"text"`
						Type      int64  `json:"type"`
						URL       string `json:"url"`
					} `json:"[doge]"`
					妙啊 struct {
						Attr      int64  `json:"attr"`
						ID        int64  `json:"id"`
						JumpTitle string `json:"jump_title"`
						Meta      struct {
							Size int64 `json:"size"`
						} `json:"meta"`
						Mtime     int64  `json:"mtime"`
						PackageID int64  `json:"package_id"`
						State     int64  `json:"state"`
						Text      string `json:"text"`
						Type      int64  `json:"type"`
						URL       string `json:"url"`
					} `json:"[妙啊]"`
					捂眼 struct {
						Attr      int64  `json:"attr"`
						ID        int64  `json:"id"`
						JumpTitle string `json:"jump_title"`
						Meta      struct {
							Size int64 `json:"size"`
						} `json:"meta"`
						Mtime     int64  `json:"mtime"`
						PackageID int64  `json:"package_id"`
						State     int64  `json:"state"`
						Text      string `json:"text"`
						Type      int64  `json:"type"`
						URL       string `json:"url"`
					} `json:"[捂眼]"`
					滑稽 struct {
						Attr      int64  `json:"attr"`
						ID        int64  `json:"id"`
						JumpTitle string `json:"jump_title"`
						Meta      struct {
							Size int64 `json:"size"`
						} `json:"meta"`
						Mtime     int64  `json:"mtime"`
						PackageID int64  `json:"package_id"`
						State     int64  `json:"state"`
						Text      string `json:"text"`
						Type      int64  `json:"type"`
						URL       string `json:"url"`
					} `json:"[滑稽]"`
					疑惑 struct {
						Attr      int64  `json:"attr"`
						ID        int64  `json:"id"`
						JumpTitle string `json:"jump_title"`
						Meta      struct {
							Size int64 `json:"size"`
						} `json:"meta"`
						Mtime     int64  `json:"mtime"`
						PackageID int64  `json:"package_id"`
						State     int64  `json:"state"`
						Text      string `json:"text"`
						Type      int64  `json:"type"`
						URL       string `json:"url"`
					} `json:"[疑惑]"`
					笑哭 struct {
						Attr      int64  `json:"attr"`
						ID        int64  `json:"id"`
						JumpTitle string `json:"jump_title"`
						Meta      struct {
							Size int64 `json:"size"`
						} `json:"meta"`
						Mtime     int64  `json:"mtime"`
						PackageID int64  `json:"package_id"`
						State     int64  `json:"state"`
						Text      string `json:"text"`
						Type      int64  `json:"type"`
						URL       string `json:"url"`
					} `json:"[笑哭]"`
					羞羞 struct {
						Attr      int64  `json:"attr"`
						ID        int64  `json:"id"`
						JumpTitle string `json:"jump_title"`
						Meta      struct {
							Size int64 `json:"size"`
						} `json:"meta"`
						Mtime     int64  `json:"mtime"`
						PackageID int64  `json:"package_id"`
						State     int64  `json:"state"`
						Text      string `json:"text"`
						Type      int64  `json:"type"`
						URL       string `json:"url"`
					} `json:"[羞羞]"`
					藏狐 struct {
						Attr      int64  `json:"attr"`
						ID        int64  `json:"id"`
						JumpTitle string `json:"jump_title"`
						Meta      struct {
							Size int64 `json:"size"`
						} `json:"meta"`
						Mtime     int64  `json:"mtime"`
						PackageID int64  `json:"package_id"`
						State     int64  `json:"state"`
						Text      string `json:"text"`
						Type      int64  `json:"type"`
						URL       string `json:"url"`
					} `json:"[藏狐]"`
				} `json:"emote"`
				JumpURL struct {
					希琳王妃 struct {
						AppName        string `json:"app_name"`
						AppPackageName string `json:"app_package_name"`
						AppURLSchema   string `json:"app_url_schema"`
						ClickReport    string `json:"click_report"`
						ExposureReport string `json:"exposure_report"`
						Extra          struct {
							GoodsClickReport    string `json:"goods_click_report"`
							GoodsCmControl      int64  `json:"goods_cm_control"`
							GoodsExposureReport string `json:"goods_exposure_report"`
							GoodsShowType       int64  `json:"goods_show_type"`
							IsWordSearch        bool   `json:"is_word_search"`
						} `json:"extra"`
						IconPosition int64  `json:"icon_position"`
						IsHalfScreen bool   `json:"is_half_screen"`
						MatchOnce    bool   `json:"match_once"`
						PcURL        string `json:"pc_url"`
						PrefixIcon   string `json:"prefix_icon"`
						State        int64  `json:"state"`
						Title        string `json:"title"`
						Underline    bool   `json:"underline"`
					} `json:"希琳王妃"`
					梶裕贵 struct {
						AppName        string `json:"app_name"`
						AppPackageName string `json:"app_package_name"`
						AppURLSchema   string `json:"app_url_schema"`
						ClickReport    string `json:"click_report"`
						ExposureReport string `json:"exposure_report"`
						Extra          struct {
							GoodsClickReport    string `json:"goods_click_report"`
							GoodsCmControl      int64  `json:"goods_cm_control"`
							GoodsExposureReport string `json:"goods_exposure_report"`
							GoodsShowType       int64  `json:"goods_show_type"`
							IsWordSearch        bool   `json:"is_word_search"`
						} `json:"extra"`
						IconPosition int64  `json:"icon_position"`
						IsHalfScreen bool   `json:"is_half_screen"`
						MatchOnce    bool   `json:"match_once"`
						PcURL        string `json:"pc_url"`
						PrefixIcon   string `json:"prefix_icon"`
						State        int64  `json:"state"`
						Title        string `json:"title"`
						Underline    bool   `json:"underline"`
					} `json:"梶裕贵"`
				} `json:"jump_url"`
				MaxLine int64         `json:"max_line"`
				Members []interface{} `json:"members"`
				Message string        `json:"message"`
			} `json:"content"`
			Replies []struct {
				Content struct {
					AtNameToMid struct {
						Mr_Sinsimito   int64 `json:"Mr-Sinsimito"`
						TroubleBear麻烦熊 int64 `json:"TroubleBear麻烦熊"`
						Wmgiii         int64 `json:"WMGIII"`
						索然無味           int64 `json:"_索然無味"`
						北邙流浪           int64 `json:"北邙流浪"`
						可怜体无比          int64 `json:"可怜体无比"`
						天选之人嗷嗷         int64 `json:"天选之人嗷嗷"`
						姆莱先生           int64 `json:"姆莱先生"`
						斯桉山            int64 `json:"斯桉山"`
						新来的给我站住        int64 `json:"新来的给我站住"`
						枫叶秦            int64 `json:"枫叶秦"`
						海盗战天下第几        int64 `json:"海盗战天下第_几"`
						甲烷冰            int64 `json:"甲烷冰"`
						籁lie           int64 `json:"籁lie"`
						节奏裂开来          int64 `json:"节奏裂开来"`
						闪子哥丷           int64 `json:"闪子哥丷"`
						青汁愈人           int64 `json:"青汁愈人"`
						魔法达比           int64 `json:"魔法达比"`
						鱼蛋占            int64 `json:"鱼蛋占"`
					} `json:"at_name_to_mid"`
					Emote struct {
						Doge struct {
							Attr      int64  `json:"attr"`
							ID        int64  `json:"id"`
							JumpTitle string `json:"jump_title"`
							Meta      struct {
								Size int64 `json:"size"`
							} `json:"meta"`
							Mtime     int64  `json:"mtime"`
							PackageID int64  `json:"package_id"`
							State     int64  `json:"state"`
							Text      string `json:"text"`
							Type      int64  `json:"type"`
							URL       string `json:"url"`
						} `json:"[doge]"`
						保佑 struct {
							Attr      int64  `json:"attr"`
							ID        int64  `json:"id"`
							JumpTitle string `json:"jump_title"`
							Meta      struct {
								Size int64 `json:"size"`
							} `json:"meta"`
							Mtime     int64  `json:"mtime"`
							PackageID int64  `json:"package_id"`
							State     int64  `json:"state"`
							Text      string `json:"text"`
							Type      int64  `json:"type"`
							URL       string `json:"url"`
						} `json:"[保佑]"`
						妙啊 struct {
							Attr      int64  `json:"attr"`
							ID        int64  `json:"id"`
							JumpTitle string `json:"jump_title"`
							Meta      struct {
								Size int64 `json:"size"`
							} `json:"meta"`
							Mtime     int64  `json:"mtime"`
							PackageID int64  `json:"package_id"`
							State     int64  `json:"state"`
							Text      string `json:"text"`
							Type      int64  `json:"type"`
							URL       string `json:"url"`
						} `json:"[妙啊]"`
						思考 struct {
							Attr      int64  `json:"attr"`
							ID        int64  `json:"id"`
							JumpTitle string `json:"jump_title"`
							Meta      struct {
								Size int64 `json:"size"`
							} `json:"meta"`
							Mtime     int64  `json:"mtime"`
							PackageID int64  `json:"package_id"`
							State     int64  `json:"state"`
							Text      string `json:"text"`
							Type      int64  `json:"type"`
							URL       string `json:"url"`
						} `json:"[思考]"`
						打call struct {
							Attr      int64  `json:"attr"`
							ID        int64  `json:"id"`
							JumpTitle string `json:"jump_title"`
							Meta      struct {
								Size int64 `json:"size"`
							} `json:"meta"`
							Mtime     int64  `json:"mtime"`
							PackageID int64  `json:"package_id"`
							State     int64  `json:"state"`
							Text      string `json:"text"`
							Type      int64  `json:"type"`
							URL       string `json:"url"`
						} `json:"[打call]"`
						疑惑 struct {
							Attr      int64  `json:"attr"`
							ID        int64  `json:"id"`
							JumpTitle string `json:"jump_title"`
							Meta      struct {
								Size int64 `json:"size"`
							} `json:"meta"`
							Mtime     int64  `json:"mtime"`
							PackageID int64  `json:"package_id"`
							State     int64  `json:"state"`
							Text      string `json:"text"`
							Type      int64  `json:"type"`
							URL       string `json:"url"`
						} `json:"[疑惑]"`
						笑哭 struct {
							Attr      int64  `json:"attr"`
							ID        int64  `json:"id"`
							JumpTitle string `json:"jump_title"`
							Meta      struct {
								Size int64 `json:"size"`
							} `json:"meta"`
							Mtime     int64  `json:"mtime"`
							PackageID int64  `json:"package_id"`
							State     int64  `json:"state"`
							Text      string `json:"text"`
							Type      int64  `json:"type"`
							URL       string `json:"url"`
						} `json:"[笑哭]"`
						藏狐 struct {
							Attr      int64  `json:"attr"`
							ID        int64  `json:"id"`
							JumpTitle string `json:"jump_title"`
							Meta      struct {
								Size int64 `json:"size"`
							} `json:"meta"`
							Mtime     int64  `json:"mtime"`
							PackageID int64  `json:"package_id"`
							State     int64  `json:"state"`
							Text      string `json:"text"`
							Type      int64  `json:"type"`
							URL       string `json:"url"`
						} `json:"[藏狐]"`
						辣眼睛 struct {
							Attr      int64  `json:"attr"`
							ID        int64  `json:"id"`
							JumpTitle string `json:"jump_title"`
							Meta      struct {
								Size int64 `json:"size"`
							} `json:"meta"`
							Mtime     int64  `json:"mtime"`
							PackageID int64  `json:"package_id"`
							State     int64  `json:"state"`
							Text      string `json:"text"`
							Type      int64  `json:"type"`
							URL       string `json:"url"`
						} `json:"[辣眼睛]"`
					} `json:"emote"`
					JumpURL struct{} `json:"jump_url"`
					MaxLine int64    `json:"max_line"`
					Members []struct {
						Avatar         string `json:"avatar"`
						FaceNftNew     int64  `json:"face_nft_new"`
						IsSeniorMember int64  `json:"is_senior_member"`
						LevelInfo      struct {
							CurrentExp   int64 `json:"current_exp"`
							CurrentLevel int64 `json:"current_level"`
							CurrentMin   int64 `json:"current_min"`
							NextExp      int64 `json:"next_exp"`
						} `json:"level_info"`
						Mid       string `json:"mid"`
						Nameplate struct {
							Condition  string `json:"condition"`
							Image      string `json:"image"`
							ImageSmall string `json:"image_small"`
							Level      string `json:"level"`
							Name       string `json:"name"`
							Nid        int64  `json:"nid"`
						} `json:"nameplate"`
						OfficialVerify struct {
							Desc string `json:"desc"`
							Type int64  `json:"type"`
						} `json:"official_verify"`
						Pendant struct {
							Expire            int64  `json:"expire"`
							Image             string `json:"image"`
							ImageEnhance      string `json:"image_enhance"`
							ImageEnhanceFrame string `json:"image_enhance_frame"`
							Name              string `json:"name"`
							Pid               int64  `json:"pid"`
						} `json:"pendant"`
						Rank   string   `json:"rank"`
						Senior struct{} `json:"senior"`
						Sex    string   `json:"sex"`
						Sign   string   `json:"sign"`
						Uname  string   `json:"uname"`
						Vip    struct {
							AccessStatus    int64  `json:"accessStatus"`
							AvatarSubscript int64  `json:"avatar_subscript"`
							DueRemark       string `json:"dueRemark"`
							Label           struct {
								BgColor               string `json:"bg_color"`
								BgStyle               int64  `json:"bg_style"`
								BorderColor           string `json:"border_color"`
								ImgLabelURIHans       string `json:"img_label_uri_hans"`
								ImgLabelURIHansStatic string `json:"img_label_uri_hans_static"`
								ImgLabelURIHant       string `json:"img_label_uri_hant"`
								ImgLabelURIHantStatic string `json:"img_label_uri_hant_static"`
								LabelTheme            string `json:"label_theme"`
								Path                  string `json:"path"`
								Text                  string `json:"text"`
								TextColor             string `json:"text_color"`
								UseImgLabel           bool   `json:"use_img_label"`
							} `json:"label"`
							NicknameColor string `json:"nickname_color"`
							ThemeType     int64  `json:"themeType"`
							VipDueDate    int64  `json:"vipDueDate"`
							VipStatus     int64  `json:"vipStatus"`
							VipStatusWarn string `json:"vipStatusWarn"`
							VipType       int64  `json:"vipType"`
						} `json:"vip"`
					} `json:"members"`
					Message string `json:"message"`
				} `json:"content"`
			} `json:"replies"`
		} `json:"replies"`
	} `json:"data"`
}

func minDeletions(s string) int {
	cnt := make([]int, 26)
	for _, c := range s {
		cnt[c-'a']++
	}
	res := 0
	m := make(map[int]int)
	for i := 0; i < 26; i++ {
		cur := cnt[i]
		for cur != 0 && m[cur] == 1 {
			res++
			cur--
		}
		m[cur] = 1
	}
	return res
}
