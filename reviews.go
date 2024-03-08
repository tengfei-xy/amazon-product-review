package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/tengfei-xy/go-log"
)

type review struct {
	review_id string
	title     string
	name      string
	time      string
	star      float32
	color     string
	body      string
	img_list  []string
}

// 返回值
// 0 => 存在且不更新
// 1 => 需要插入
// 2 => 需要更新
func (r *review) exist_review_id(asin, checksum string) (int, error) {
	result := app.db.QueryRow("select  checksum  from review where review_id=? and asin=?", r.review_id, asin)
	if err := result.Err(); err != nil {
		return 0, err
	}
	var out_checksum string
	if err := result.Scan(&out_checksum); err != nil {
		if err == sql.ErrNoRows {
			return 1, nil
		}
		return 0, err
	}

	if out_checksum != checksum {
		return 2, nil
	}
	return 0, nil
}
func (r *review) insert(asin, checksum, language string) error {
	img_list := strings.Join(r.img_list, ",")
	// 如果评论ID存在
	_, err := app.db.Exec("insert into review(asin,review_id,language,title,name,time,star,color,body,img_list,checksum) values(?,?,?,?,?,?,?,?,?,?,?)", asin, r.review_id, language, r.title, r.name, r.time, r.star, r.color, r.body, img_list, checksum)

	if err != nil {
		return err
	}
	// _, err := result.RowsAffected()
	// if err != nil {
	// 	return err
	// }
	log.Infof("插入成功 评论ID:%s", r.review_id)
	return nil
}
func (r *review) update(asin, checksum, language string) error {
	img_list := strings.Join(r.img_list, ",")
	_, err := app.db.Exec("update review set asin=?,language=?,title=?,name=?,time=?,star=?,color=?,body=?,img_list=?,checksum=? where review_id=?", asin, language, r.title, r.name, r.time, r.star, r.color, r.body, img_list, checksum, r.review_id)
	if err != nil {
		return err
	}
	log.Infof("更新成功 评论ID:%s", r.review_id)

	return nil

}
func (r *review) sha256() string {
	var buffer = sha256.New()
	buffer.Write([]byte(r.name))
	buffer.Write([]byte(r.time))
	buffer.Write([]byte(fmt.Sprintf("%.1f", r.star)))
	buffer.Write([]byte(r.title))
	buffer.Write([]byte(r.color))
	buffer.Write([]byte(r.body))
	if r.img_list != nil {
		for _, img := range r.img_list {
			buffer.Write([]byte(img))
		}
	}
	// 计算哈希值
	hashValue := buffer.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashString := hex.EncodeToString(hashValue)
	// log.Infof("校验码:%s", hashString)
	return hashString

}
func reviews_main(asin string) {

	doc, err := get_reviews_html(asin)
	if err != nil {
		log.Error(err)
		return
	}
	rs, l := get_reviews_doc(doc)
	for _, r := range rs {
		checksum := r.sha256()
		seq, err := r.exist_review_id(asin, checksum)
		if err != nil {
			log.Error(err)
			break
		}
		switch seq {
		case 0:
			log.Infof("已最新  评论ID: %s", r.review_id)
		case 1:
			log.Infof("开始插入  评论ID: %s", r.review_id)
			err := r.insert(asin, checksum, l)
			if err != nil {
				log.Error(err)
			}
		case 2:
			log.Infof("开始更新  评论ID: %s", r.review_id)
			err := r.update(asin, checksum, l)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func get_reviews_html(asin string) (*goquery.Document, error) {
	url := fmt.Sprintf("https://www.amazon.de/product-reviews/%s", asin)
	client := get_client()
	cookie, err := app.get_cookie()
	if err != nil {
		return nil, err
	}
	// curl 'https://www.amazon.de/product-reviews/B0C9ZV7BX6' \
	// -H ': ' \
	// -H ': ' \
	// -H 'accept-language: zh-CN,zh;q=0.9' \
	// -H 'cache-control: max-age=0' \
	// -H 'cookie: session-id=258-8222509-0214168; x-amz-captcha-1=1709882958766900; x-amz-captcha-2=G7DrXeag6X/d/SaC09CQPQ==; ubid-acbde=257-3991201-7608257; i18n-prefs=EUR; sp-cdn="L5Z9:CN"; lc-acbde=de_DE; csm-hit=tb:6CEDB1Z3RDPFXSQZWHCV+s-6CEDB1Z3RDPFXSQZWHCV|1709876892694&t:1709876892694&adb:adblk_yes; session-token=Y0omIofilyTyUqOCAOfuojVOf0yGZ2e18U+SODyORF6+LPM9af+fM502ABgHw+/Bc7JDbzsb+/mhemfq1QVf1cfOXisZynMzf1uthV0QkHOrtEeH001bnjGdaKRfWqt8ANPf7xo1xqNCMsBhOKwE8LQL4Hk6PrCLezssN0TM1G9Tlz2hvMWWuKbLKjOYK/EVoip2D3P961x/2H7qiX+GjLQDBWLXyNQQ0pk0LHy50Bq18IGUC89gLQlCFsTrekuBcht8M7OcdZ+AUjLRmTQn8N+9HLZWDdKqpNZ5lIted3bRVaKyUPtA+bb/pqqPIIStNE434/LaB1AaKxG/rq4h0/g06P14/jyG; session-id-time=2082754801l' \
	// -H 'device-memory: 8' \
	// -H 'downlink: 2.4' \
	// -H 'dpr: 2' \
	// -H 'ect: 4g' \
	// -H 'rtt: 250' \
	// -H 'upgrade-insecure-requests: 1' \
	// -H ': ' \
	// -H 'viewport-width: 1024'

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err

	}
	req.Header.Set("authority", app.Domain)
	req.Header.Set("accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`)
	req.Header.Set("accept-language", `zh-CN,zh;q=0.9`)
	req.Header.Set("cache-control", `max-age=0`)
	req.Header.Set("cookie", cookie)
	req.Header.Set("upgrade-insecure-requests", `1`)
	req.Header.Set("user-agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36`)

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("内部错误:%v", err)
		return nil, err

	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		break
	case 404:
		return nil, fmt.Errorf("状态码:%d, 访问链接:%s", resp.StatusCode, url)
	case 503:
		return nil, fmt.Errorf("状态码:%d, 访问链接:%s", resp.StatusCode, url)
	default:
		return nil, fmt.Errorf("状态码:%d, 访问链接:%s", resp.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("内部错误:%v", err)
	}
	return doc, nil
}
func get_reviews_doc(doc *goquery.Document) ([]review, string) {
	var r_list []review
	language, _ := doc.Find("html").Attr("lang")

	switch language {
	case "de-de":
		language = "de"
	default:
		language = "en"
	}

	doc.Find("div[id=cm_cr-review_list]").Find("div[data-hook=review]").Each(func(i int, s *goquery.Selection) {
		var r review

		review_id, _ := s.Attr("id")
		r.review_id = review_id
		log.Infof("发现评论ID: %s", review_id)

		name := s.Find("span[class=a-profile-name]").Text()
		r.name = name
		log.Infof("发现评论者: %s", name)

		time := s.Find("span[data-hook=review-date]").Text()
		time = trans_time_de(time)
		r.time = time
		log.Infof("发现评论时间: %s", time)

		star := s.Find("a[data-hook=review-title]>i").Text()
		star_float := trans_star_de(star)
		r.star = star_float
		log.Infof("发现评论星级: %.1f", star_float)

		title := s.Find("a[data-hook=review-title]>span+span").Text()
		title = strings.TrimSpace(title)
		r.title = title
		log.Infof("发现评论标题: %s", title)

		color := s.Find("a[data-hook=format-strip]").Text()
		color = trans_color_de(color)
		r.color = color
		log.Infof("发现评论款式: %s", color)

		body := s.Find("span[data-hook=review-body]").Text()
		body = strings.TrimSpace(body)
		r.body = body
		log.Infof("发现评论内容长度: %d", len(body))

		img := s.Find("img[data-hook=review-image-tile]")
		if img.Length() == 0 {
			r.img_list = nil
		} else {
			img_list := make([]string, img.Length())
			img.Each(func(i int, s *goquery.Selection) {
				img, _ := s.Attr("src")
				log.Infof("发现评论图片:% v", img)
				img_list[i] = img
			})
			r.img_list = img_list
		}

		log.Info("--------------------------")
		r_list = append(r_list, r)
	})
	return r_list, language
}
func trans_color_de(color string) string {
	color = strings.TrimLeft(color, "Farbe: ")
	return color
}
func trans_star_de(star string) float32 {
	s := strings.TrimSuffix(star, ",0 von 5 Sternen")
	// 转换为小数
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 10
	}
	return float32(f)
}
func trans_time_de(dateStr string) string {
	dateStr = strings.TrimLeft(dateStr, "Rezension aus Deutschland vom ")
	t := strings.Split(dateStr, " ")
	day := strings.TrimRight(t[0], ".")
	if len(day) == 1 {
		day = "0" + day
	}
	month := t[1]
	switch month {
	case "Januar":
		month = "01"
	case "Februar":
		month = "02"
	case "März":
		month = "03"
	case "April":
		month = "04"
	case "Mai":
		month = "05"
	case "Juni":
		month = "06"
	case "Juli":
		month = "07"
	case "August":
		month = "08"
	case "September":
		month = "09"
	case "Oktober":
		month = "10"
	case "November":
		month = "11"
	case "Dezember":
		month = "12"
	}
	year := t[2]
	return year + "-" + month + "-" + day
}
