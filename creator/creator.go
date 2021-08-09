package creator

import (
	"encoding/csv"
	"fmt"
	"github.com/duck8823/sample-dataset-creator-for-amazon-personalize/models"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

// Creator は Amazon Personalize用のデータセット作るくん
type Creator interface {
	Create() error
}

// CsvCreator は CSV形式で作るくん
type CsvCreator struct {
	Output string
}

func (c CsvCreator) Create() error {
	items, err := items()
	if err != nil {
		return xerrors.Errorf("コンテンツの生成に失敗しました: %w", err)
	}
	users, err := users(items)
	if err != nil {
		return xerrors.Errorf("ユーザーの生成に失敗しました: %w", err)
	}
	interactions, err := interactions(users, items)
	if err != nil {
		return xerrors.Errorf("行動データの生成に失敗しました: %w", err)
	}

	itf, err := os.OpenFile(path.Join(c.Output, "items.csv"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return xerrors.Errorf("ファイルの作成に失敗しました: %w", err)
	}
	iw := csv.NewWriter(itf)
	if err := iw.Write([]string{"ITEM_ID", "CATEGORY", "CREATION_TIMESTAMP"}); err != nil {
		return xerrors.Errorf("ヘッダーの出力に失敗しました: %w", err)
	}
	for _, item := range items {
		if err := iw.Write([]string{
			item.ItemID,
			string(*item.Category),
			fmt.Sprintf("%d", item.CreationTimestamp),
		}); err != nil {
			return xerrors.Errorf("コンテンツの出力に失敗しました: %w", err)
		}
	}

	uf, err := os.OpenFile(path.Join(c.Output, "users.csv"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return xerrors.Errorf("ファイルの作成に失敗しました: %w", err)
	}
	uw := csv.NewWriter(uf)
	if err := uw.Write([]string{"USER_ID", "PREFECTURES", "BOOKMARKS"}); err != nil {
		return xerrors.Errorf("ヘッダーの出力に失敗しました: %w", err)
	}
	for _, user := range users {
		prefectures := make([]string, len(user.Prefectures))
		for i, prefecture := range user.Prefectures {
			prefectures[i] = string(prefecture)
		}
		if err := uw.Write([]string{
			user.UserID,
			strings.Join(prefectures, "|"),
			strings.Join(user.Bookmarks, "|"),
		}); err != nil {
			return xerrors.Errorf("ユーザーの出力に失敗しました: %w", err)
		}
	}

	inf, err := os.OpenFile(path.Join(c.Output, "interactions.csv"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return xerrors.Errorf("ファイルの作成に失敗しました: %w", err)
	}
	inw := csv.NewWriter(inf)
	if err := inw.Write([]string{"USER_ID", "ITEM_ID", "TIMESTAMP"}); err != nil {
		return xerrors.Errorf("ヘッダーの出力に失敗しました: %w", err)
	}
	for _, interaction := range interactions {
		if err := inw.Write([]string{
			interaction.UserID,
			interaction.ItemID,
			fmt.Sprintf("%d", interaction.Timestamp),
		}); err != nil {
			return xerrors.Errorf("行動データの出力に失敗しました: %w", err)
		}
	}

	return nil
}

func items() ([]models.Item, error) {
	categories := []models.Category{
		models.Action,
		models.Adventure,
		models.Music,
		models.Puzzle,
		models.Race,
		models.RolePlaying,
		models.Shooting,
		models.Simulation,
		models.Supports,
	}
	rand.Seed(time.Now().Unix())

	items := make([]models.Item, 3000)
	for i, item := range items {
		item.ItemID = uuid.New().String()
		item.Category = &categories[rand.Intn(len(categories))]
		item.CreationTimestamp = time.Now().Unix()
		items[i] = item
	}
	return items, nil
}

func users(items []models.Item) ([]models.User, error) {
	prefectures := []models.Prefecture{
		models.Aichi,
		models.Akita,
		models.Aomori,
		models.Chiba,
		models.Ehime,
		models.Fukui,
		models.Fukuoka,
		models.Fukushima,
		models.Gifu,
		models.Gunma,
		models.Hiroshima,
		models.Hokkaido,
		models.Hyogo,
		models.Ibaraki,
		models.Ishikawa,
		models.Iwate,
		models.Kagawa,
		models.Kagoshima,
		models.Kanagawa,
		models.Kochi,
		models.Kumamoto,
		models.Kyoto,
		models.Mie,
		models.Miyagi,
		models.Miyazaki,
		models.Nagano,
		models.Nagasaki,
		models.Nara,
		models.Niigata,
		models.Oita,
		models.Okayama,
		models.Okinawa,
		models.Osaka,
		models.Saga,
		models.Saitama,
		models.Shiga,
		models.Shimane,
		models.Shizuoka,
		models.Tochigi,
		models.Tokushima,
		models.Tokyo,
		models.Tottori,
		models.Toyama,
		models.Wakayama,
		models.Yamagata,
		models.Yamaguchi,
		models.Yamanashi,
	}
	rand.Seed(time.Now().Unix())

	users := make([]models.User, 500)
	for i, user := range users {
		user.UserID = uuid.New().String()
		// 0 から 2 個の都道府県
		for n := 0; n <= rand.Intn(2); n++ {
			prefecture := prefectures[rand.Intn(len(prefectures))]
			user.Prefectures = append(user.Prefectures, prefecture)
		}

		// 0 から 5 個のお気に入り
		for n := 0; n <= rand.Intn(5); n++ {
			item := items[rand.Intn(len(items))]
			user.Bookmarks = append(user.Bookmarks, item.ItemID)
		}
		users[i] = user
	}
	return users, nil
}

func interactions(users []models.User, items []models.Item) ([]models.Interaction, error) {
	interactions := make([]models.Interaction, 20000)
	// ランダム
	for i := 0; i < 18000; i++ {
		interactions[i].UserID = users[rand.Intn(len(users))].UserID
		interactions[i].ItemID = items[rand.Intn(len(items))].ItemID
		interactions[i].Timestamp = time.Now().Unix()
	}
	// 近い数字で回す
	for i := 18000; i < 20000; i++ {
		interactions[i].UserID = users[rand.Intn(20)].UserID
		interactions[i].ItemID = items[rand.Intn(150)].ItemID
		interactions[i].Timestamp = time.Now().Unix()
	}
	return interactions, nil
}
