package seed

import (
	"context"
	"fmt"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
	"github.com/google/uuid"
)

func ptr[T any](v T) *T { return &v }

// LoadCards seeds the card repository with initial card data.
func LoadCards(ctx context.Context, repo repository.CardRepository) error {
	cards := allCards()
	for i := range cards {
		cards[i].ID = uuid.New().String()
		cards[i].CreatedAt = time.Now()
		if err := repo.Create(ctx, &cards[i]); err != nil {
			return fmt.Errorf("seed card %s: %w", cards[i].Name, err)
		}
	}
	return nil
}

func allCards() []model.Card {
	var cards []model.Card
	cards = append(cards, constellationCards()...)
	cards = append(cards, planetCards()...)
	cards = append(cards, phenomenonCards()...)
	return cards
}

func constellationCards() []model.Card {
	return []model.Card{
		{Name: "オリオン座", Category: model.CardCategoryConstellation, ReadingText: "冬の夜空に三つ星が並ぶ、狩人の姿", Description: "冬を代表する星座。ベテルギウスとリゲルが目印。三つ星は赤道近くに位置し、世界中から見える。", Magnitude: ptr(0.12), Distance: ptr("1,344光年"), BestSeason: "winter"},
		{Name: "おおぐま座", Category: model.CardCategoryConstellation, ReadingText: "北の空で柄杓の形、北極星への道しるべ", Description: "北斗七星を含む大きな星座。柄杓の先端から北極星を見つけることができる。", Magnitude: ptr(1.79), Distance: ptr("123光年"), BestSeason: "spring"},
		{Name: "さそり座", Category: model.CardCategoryConstellation, ReadingText: "赤い心臓アンタレスが光る、夏の毒虫", Description: "夏の代表的な星座。赤い一等星アンタレスは「火星に対抗するもの」という意味。", Magnitude: ptr(0.96), Distance: ptr("554光年"), BestSeason: "summer"},
		{Name: "こと座", Category: model.CardCategoryConstellation, ReadingText: "織姫星ベガが輝く、夏の大三角の一角", Description: "一等星ベガは夏の大三角の一つ。七夕の織姫星として知られる。", Magnitude: ptr(0.03), Distance: ptr("25光年"), BestSeason: "summer"},
		{Name: "はくちょう座", Category: model.CardCategoryConstellation, ReadingText: "天の川を渡る白鳥、デネブが道標", Description: "夏の大三角の一角デネブを持つ。十字形が天の川に沿って広がる。", Magnitude: ptr(1.25), Distance: ptr("2,615光年"), BestSeason: "summer"},
		{Name: "わし座", Category: model.CardCategoryConstellation, ReadingText: "彦星アルタイルが瞬く、天の川の番人", Description: "一等星アルタイルは夏の大三角の一つ。七夕の彦星として親しまれる。", Magnitude: ptr(0.76), Distance: ptr("16.7光年"), BestSeason: "summer"},
		{Name: "おうし座", Category: model.CardCategoryConstellation, ReadingText: "すばるの宝石箱を背負う、赤い目の牡牛", Description: "アルデバランの赤い輝きとプレアデス星団（すばる）が見どころ。黄道十二星座の一つ。", Magnitude: ptr(0.85), Distance: ptr("65光年"), BestSeason: "winter"},
		{Name: "ふたご座", Category: model.CardCategoryConstellation, ReadingText: "カストルとポルックス、仲良し双子の星", Description: "二つの明るい星が双子の頭に輝く。ふたご座流星群の放射点。", Magnitude: ptr(1.14), Distance: ptr("34光年"), BestSeason: "winter"},
		{Name: "しし座", Category: model.CardCategoryConstellation, ReadingText: "春の夜空に寝そべる百獣の王", Description: "春の代表的な星座。レグルスが心臓の位置に輝く。しし座流星群で有名。", Magnitude: ptr(1.35), Distance: ptr("79光年"), BestSeason: "spring"},
		{Name: "おとめ座", Category: model.CardCategoryConstellation, ReadingText: "麦の穂を持つ乙女、スピカの青白い光", Description: "春の星座。一等星スピカは「麦の穂」を意味する。おとめ座銀河団がある。", Magnitude: ptr(0.97), Distance: ptr("250光年"), BestSeason: "spring"},
		{Name: "みずがめ座", Category: model.CardCategoryConstellation, ReadingText: "水瓶から流れ出す星の水、秋の淡い光", Description: "秋の黄道十二星座。明るい星は少ないが、みずがめ座流星群で知られる。", Magnitude: ptr(2.91), BestSeason: "autumn"},
		{Name: "うお座", Category: model.CardCategoryConstellation, ReadingText: "リボンでつながれた二匹の魚、春分点の星座", Description: "黄道十二星座の一つ。現在の春分点がこの星座にある。", Magnitude: ptr(3.62), BestSeason: "autumn"},
		{Name: "おひつじ座", Category: model.CardCategoryConstellation, ReadingText: "金の羊毛を持つ小さな牡羊", Description: "黄道十二星座の一つ。ギリシャ神話の黄金の羊にちなむ。", Magnitude: ptr(2.00), Distance: ptr("66光年"), BestSeason: "autumn"},
		{Name: "かに座", Category: model.CardCategoryConstellation, ReadingText: "プレセペ星団を胸に抱く、控えめな蟹", Description: "黄道十二星座の一つ。散開星団プレセペ（M44）が肉眼でも見える。", Magnitude: ptr(3.52), Distance: ptr("290光年"), BestSeason: "spring"},
		{Name: "てんびん座", Category: model.CardCategoryConstellation, ReadingText: "正義の女神の天秤、秋の入り口に輝く", Description: "黄道十二星座の一つ。もともとさそり座のハサミだった。", Magnitude: ptr(2.61), BestSeason: "summer"},
		{Name: "いて座", Category: model.CardCategoryConstellation, ReadingText: "銀河の中心を射る弓、天の川で最も賑やか", Description: "天の川銀河の中心方向にある。南斗六星を含む。多くの星雲や星団が密集。", Magnitude: ptr(1.85), Distance: ptr("228光年"), BestSeason: "summer"},
		{Name: "やぎ座", Category: model.CardCategoryConstellation, ReadingText: "山羊の頭に魚の尻尾、不思議な姿の星座", Description: "黄道十二星座の一つ。上半身が山羊で下半身が魚の不思議な姿。", Magnitude: ptr(2.87), BestSeason: "autumn"},
		{Name: "カシオペヤ座", Category: model.CardCategoryConstellation, ReadingText: "Wの形が目印、北極星を探す手がかり", Description: "北天の星座。W字型が特徴的で、北極星を見つけるのに使われる。", Magnitude: ptr(2.23), Distance: ptr("228光年"), BestSeason: "autumn"},
		{Name: "ペルセウス座", Category: model.CardCategoryConstellation, ReadingText: "メデューサの首を持つ英雄、流星雨の故郷", Description: "ペルセウス座流星群の放射点。変光星アルゴルで知られる。", Magnitude: ptr(1.79), Distance: ptr("510光年"), BestSeason: "winter"},
		{Name: "アンドロメダ座", Category: model.CardCategoryConstellation, ReadingText: "鎖につながれた王女、隣の銀河への入り口", Description: "アンドロメダ銀河（M31）を含む。肉眼で見える最も遠い天体がここにある。", Magnitude: ptr(2.06), Distance: ptr("254万光年"), BestSeason: "autumn"},
		{Name: "こぐま座", Category: model.CardCategoryConstellation, ReadingText: "尻尾の先に北極星、小さな子熊の星座", Description: "北極星（ポラリス）を持つ星座。地球の自転軸の延長上にある。", Magnitude: ptr(1.97), Distance: ptr("431光年"), BestSeason: "all"},
		{Name: "おおいぬ座", Category: model.CardCategoryConstellation, ReadingText: "全天一明るいシリウス、冬のダイヤモンド", Description: "全天で最も明るい恒星シリウスを持つ。オリオンの猟犬として描かれる。", Magnitude: ptr(-1.46), Distance: ptr("8.6光年"), BestSeason: "winter"},
		{Name: "こいぬ座", Category: model.CardCategoryConstellation, ReadingText: "プロキオンが光る、小さくても存在感のある子犬", Description: "プロキオンは冬の大三角の一つ。「犬の前に」という意味で、シリウスより先に昇る。", Magnitude: ptr(0.34), Distance: ptr("11.5光年"), BestSeason: "winter"},
		{Name: "ぎょしゃ座", Category: model.CardCategoryConstellation, ReadingText: "カペラが黄色く輝く、冬の五角形の一角", Description: "一等星カペラは全天で6番目に明るい。冬のダイヤモンドの一角。", Magnitude: ptr(0.08), Distance: ptr("43光年"), BestSeason: "winter"},
		{Name: "りゅう座", Category: model.CardCategoryConstellation, ReadingText: "北極の周りをうねる竜、かつての北極星の主", Description: "北天の星座。約5000年前にはりゅう座のトゥバンが北極星だった。", Magnitude: ptr(2.24), Distance: ptr("303光年"), BestSeason: "summer"},
		{Name: "ケンタウルス座", Category: model.CardCategoryConstellation, ReadingText: "最も近い恒星を持つ、南天の半人半馬", Description: "太陽に最も近い恒星プロキシマ・ケンタウリがある。日本からは一部しか見えない。", Magnitude: ptr(-0.27), Distance: ptr("4.2光年"), BestSeason: "spring"},
		{Name: "みなみじゅうじ座", Category: model.CardCategoryConstellation, ReadingText: "南半球のシンボル、最も小さな星座", Description: "全88星座で最も小さい。南半球の方角を知るのに使われる。", Magnitude: ptr(0.77), Distance: ptr("320光年"), BestSeason: "spring"},
		{Name: "ヘルクレス座", Category: model.CardCategoryConstellation, ReadingText: "天を支える英雄、球状星団M13の住処", Description: "ギリシャ最大の英雄。M13球状星団は双眼鏡で楽しめる。", Magnitude: ptr(2.81), Distance: ptr("112光年"), BestSeason: "summer"},
		{Name: "ペガスス座", Category: model.CardCategoryConstellation, ReadingText: "秋の四辺形が目印、翼を持つ天馬", Description: "ペガススの四辺形（秋の大四辺形）は秋の星座探しの起点。", Magnitude: ptr(2.42), Distance: ptr("133光年"), BestSeason: "autumn"},
		{Name: "うしかい座", Category: model.CardCategoryConstellation, ReadingText: "オレンジ色のアークトゥルス、春の大曲線の先", Description: "アークトゥルスは全天で4番目に明るい恒星。北斗七星の柄から伸ばして見つける。", Magnitude: ptr(-0.05), Distance: ptr("37光年"), BestSeason: "spring"},
	}
}

func planetCards() []model.Card {
	return []model.Card{
		{Name: "水星", Category: model.CardCategoryPlanet, ReadingText: "太陽に最も近い、足の速い使者の星", Description: "太陽系で最小の惑星。公転周期88日。昼は430℃、夜は-180℃と温度差が最大。", Distance: ptr("0.39 AU"), BestSeason: "all"},
		{Name: "金星", Category: model.CardCategoryPlanet, ReadingText: "明けの明星、宵の明星、最も明るい惑星", Description: "地球の姉妹惑星。厚い大気で温室効果が極端。表面温度は約460℃。", Magnitude: ptr(-4.6), Distance: ptr("0.72 AU"), BestSeason: "all"},
		{Name: "地球", Category: model.CardCategoryPlanet, ReadingText: "青い水の惑星、私たちの唯一の故郷", Description: "太陽系で唯一液体の水が表面に存在する惑星。月という大きな衛星を持つ。", Distance: ptr("1.00 AU"), BestSeason: "all"},
		{Name: "火星", Category: model.CardCategoryPlanet, ReadingText: "赤い砂漠の惑星、太陽系最大の火山がそびえる", Description: "オリンポス山は太陽系最大の火山。将来の有人探査の候補地。", Magnitude: ptr(-2.91), Distance: ptr("1.52 AU"), BestSeason: "all"},
		{Name: "木星", Category: model.CardCategoryPlanet, ReadingText: "太陽系最大のガス巨人、大赤斑が回る", Description: "太陽系最大の惑星。大赤斑は地球2個分の巨大な嵐。79個以上の衛星を持つ。", Magnitude: ptr(-2.94), Distance: ptr("5.20 AU"), BestSeason: "all"},
		{Name: "土星", Category: model.CardCategoryPlanet, ReadingText: "美しい環を持つ惑星、氷と岩のリング", Description: "特徴的な環を持つ。環は氷と岩の粒子からなる。タイタンという大きな衛星を持つ。", Magnitude: ptr(-0.55), Distance: ptr("9.54 AU"), BestSeason: "all"},
		{Name: "天王星", Category: model.CardCategoryPlanet, ReadingText: "横倒しで回る氷の巨人、青緑色の謎多き星", Description: "自転軸が98度も傾いている。メタンの大気で青緑色に見える。", Magnitude: ptr(5.32), Distance: ptr("19.2 AU"), BestSeason: "all"},
		{Name: "海王星", Category: model.CardCategoryPlanet, ReadingText: "太陽系最果ての青い惑星、風速2000kmの嵐", Description: "太陽系で最も風が強い惑星。美しい青色はメタンによる。", Magnitude: ptr(7.78), Distance: ptr("30.1 AU"), BestSeason: "all"},
		{Name: "冥王星", Category: model.CardCategoryPlanet, ReadingText: "かつての第9惑星、カイパーベルトの矮小天体", Description: "2006年に準惑星に再分類。ハート型の地形が特徴的。", Magnitude: ptr(13.65), Distance: ptr("39.5 AU"), BestSeason: "all"},
	}
}

func phenomenonCards() []model.Card {
	return []model.Card{
		{Name: "皆既日食", Category: model.CardCategoryPhenomenon, ReadingText: "月が太陽を隠す奇跡の瞬間、コロナが輝く", Description: "月が太陽を完全に隠す現象。太陽コロナやプロミネンスが見える。同じ場所では平均375年に1回。", BestSeason: "all"},
		{Name: "皆既月食", Category: model.CardCategoryPhenomenon, ReadingText: "赤銅色に染まる月、地球の影の中で", Description: "地球の影に月が入る現象。大気で屈折した赤い光で月が赤銅色に見える。", BestSeason: "all"},
		{Name: "ペルセウス座流星群", Category: model.CardCategoryPhenomenon, ReadingText: "真夏の夜の流れ星、1時間に100個の願い", Description: "毎年8月12-13日頃にピーク。1時間に最大100個の流星。スイフト・タットル彗星の塵。", BestSeason: "summer"},
		{Name: "ふたご座流星群", Category: model.CardCategoryPhenomenon, ReadingText: "12月の夜空を飾る、年間最大の流星群", Description: "毎年12月13-14日頃にピーク。年間で最も多い流星群の一つ。小惑星ファエトンが母天体。", BestSeason: "winter"},
		{Name: "しし座流星群", Category: model.CardCategoryPhenomenon, ReadingText: "流星の嵐、33年周期で大出現する獅子の涙", Description: "毎年11月17-18日頃。33年周期で流星嵐が発生。2001年には1時間に数千個の記録。", BestSeason: "autumn"},
		{Name: "オーロラ", Category: model.CardCategoryPhenomenon, ReadingText: "太陽の風が描く光のカーテン、極地の奇跡", Description: "太陽風の荷電粒子が大気と衝突して発光する現象。緑や赤の光のカーテンが揺れる。", BestSeason: "winter"},
		{Name: "天の川", Category: model.CardCategoryPhenomenon, ReadingText: "数千億の星が織りなす光の帯、銀河の断面図", Description: "私たちの銀河系の円盤を内側から見た姿。夏の夜空に最も明るく見える。", BestSeason: "summer"},
		{Name: "スーパームーン", Category: model.CardCategoryPhenomenon, ReadingText: "最も近づいた満月、14%大きく30%明るく", Description: "月が地球に最も近づいた時の満月。通常より大きく明るく見える。", BestSeason: "all"},
		{Name: "彗星", Category: model.CardCategoryPhenomenon, ReadingText: "太陽に近づくと尾を引く、氷と塵の旅人", Description: "氷と塵からなる天体。太陽に近づくとガスや塵が放出され美しい尾を形成する。", BestSeason: "all"},
		{Name: "金星の太陽面通過", Category: model.CardCategoryPhenomenon, ReadingText: "太陽の顔を横切る金星、100年に2回の天文ショー", Description: "金星が太陽の前を通過する現象。次回は2117年12月。非常に稀な天文現象。", BestSeason: "all"},
	}
}
