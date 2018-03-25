package userbooks

import (
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/elastic"
	. "github.com/akiraak/go-manga/model"
	"log"
	"time"
)

func userAsins(userId int64) []string {
	keywords := []string{
		"HUNTER×HUNTER",
		"ヴィンランド・サガ",
		"落日のパトス",
		"狼と香辛料",
		"食戟のソーマ",
		"小説家になる方法",
		"重版出来",
		"山と食欲と私",
		"釣り船御前丸",
		"木根さんの1人でキネマ",
		"インベスターZ",
		"波よ聞いてくれ",
		"食戟のソーマ",
		"アルキメデスの大戦",
		"ふらいんぐうぃっち",
		"亜人",
		"宇宙兄弟",
		"BLUE GIANT",
		"ベイビーステップ",
		"ハイスコアガール",
		"蛇蔵",
		"僕らはみんな河合荘",
		"のの湯",
		"百姓貴族",
		"ドロヘドロ",
		"からかい上手の高木さん",
		"乙嫁語り",
		"ばらかもん",
		"君に届け",
		"のんのんびより",
		"海街diary",
		"後遺症ラジオ",
		"ワンパンマン",
		"いぶり暮らし",
		"ヒストリエ",
		"つれづれダイアリー",
		"ダンジョン飯",
		"メイドインアビス",
		"ドメスティックな彼女",
		"東京喰種",
		"進撃の巨人",
		"アオバ自転車店",
		"ちはやふる",
		"甘々と稲妻",
		"ろんぐらいだぁす",
		"はたらく細胞",
		"猫のお寺の知恩さん",
		"ウーパ",
		"ゆるキャン",
		"平方イコルスン",
		"山賊ダイアリー",
		"銀の匙",
		"味噌汁でカンパイ",
		"放課後さいころ倶楽部",
		"ふしぎの国のバード",
		"レイリ",
		"3月のライオン",
		"コウノドリ",
		"南鎌倉高校女子自転車部",
		"あげくの果てのカノン",
		"徒然チルドレン",
		"ぐらんぶる",
		"ドラゴン桜",
		"おひ釣りさま",
		"放課後ていぼう日誌",
		"MISS CAST",
	}
	asins, _ := elastic.SearchAsins(keywords, 0, 10000)
	return asins
}

func deleteUserBooks(userId int64) {
	db.ORM.Exec("delete from user_books where user_id = ?", userId)
}

func createUserBook(userId int64, asins []string) {
	for _, asin := range asins {
		book := Book{}
		if !db.ORM.Where("asin = ?", asin).Find(&book).RecordNotFound() {
			userBook := UserBook{}
			userBook.UserID = userId
			userBook.BookID = book.ID
			db.ORM.Create(&userBook)
		}
	}
}

func updateUserBooks(userId int64, asins []string) {
	deleteUserBooks(userId)
	createUserBook(userId, asins)
}

func update_user_books() {
	userId := int64(1)
	asins := userAsins(userId)
	updateUserBooks(userId, asins)
}

func Run() {
	startTime := time.Now();
	update_user_books()
	endTime := time.Now();
	log.Printf("userbooks.Run(): %f秒\n",(endTime.Sub(startTime)).Seconds())
}
