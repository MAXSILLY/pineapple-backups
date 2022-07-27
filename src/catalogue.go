package src

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"sf/cfg"
	"sf/src/boluobao"
	"sf/src/hbooker"
	"sf/structural/sfacg_structs"
	"strconv"
	"time"
)

type Catalogue struct {
	ChapterBar   *ProgressBar
	ChapterList  []string
	ConfigPath   string
	SaveTextPath string
	ChapterCfg   string
}

func (catalogue *Catalogue) ReadChapterConfig() {
	catalogue.ConfigPath = path.Join(cfg.Vars.ConfigFile, cfg.Vars.BookInfo.NovelName+".conf")
	catalogue.SaveTextPath = path.Join(cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName+".txt")
}
func (catalogue *Catalogue) SfacgCatalogue() bool {
	catalogue.ReadChapterConfig()
	response := boluobao.GetCatalogueDetailedById(cfg.Vars.BookInfo.NovelID)
	for divisionIndex, division := range response.Data.VolumeList {
		fmt.Printf("第%v卷\t\t%v\n", divisionIndex+1, division.Title)
		for _, chapter := range division.ChapterList {
			if chapter.OriginNeedFireMoney == 0 {
				catalogue.ChapterList = append(catalogue.ChapterList, strconv.Itoa(chapter.ChapID))
			} else {
				if chapter.OriginNeedFireMoney == 0 {
					fmt.Println(chapter.Title, "已经下载过了")
				}
			}

		}
	}
	catalogue.ChapterBar = New(len(catalogue.ChapterList))
	catalogue.ChapterBar.Describe("进度:")
	for _, ChapID := range catalogue.ChapterList {
		catalogue.SfacgContent(ChapID)
	}
	catalogue.ChapterList = nil
	return true
}

func (catalogue *Catalogue) SfacgContent(ChapterId string) {
	catalogue.DelayTime()
	response := boluobao.GetContentDetailedByCid(ChapterId)
	for i := 0; i < 5; i++ {
		if response.Status.HTTPCode == 200 {
			catalogue.makeContentInformation(response)
			break
		} else {
			if response.Status.Msg == "接口校验失败,请尽快把APP升级到最新版哦~" || i == 4 {
				fmt.Println(response.Status.Msg)
				os.Exit(0)
			} else {
				fmt.Println(response.Status.Msg)
			}
		}
	}
}

func (catalogue *Catalogue) makeContentInformation(response sfacg_structs.Content) {
	writeContent := fmt.Sprintf("%v:%v\n%v\n\n\n",
		response.Data.Title, response.Data.AddTime, response.Data.Expand.Content,
	)
	cfg.EncapsulationWrite(catalogue.SaveTextPath, writeContent, 5, "a")
	//catalogue.AddChapterConfig(response.Data.ChapID)

}

func (catalogue *Catalogue) CatCatalogue() bool {
	catalogue.ReadChapterConfig()
	for index, division := range hbooker.GetDivisionIdByBookId(cfg.Vars.BookInfo.NovelID) {
		fmt.Printf("第%v卷\t\t%v\n", index+1, division.DivisionName)
		for _, chapter := range hbooker.GetCatalogueByDivisionId(division.DivisionID) {
			if chapter.IsValid == "1" {
				catalogue.ChapterList = append(catalogue.ChapterList, chapter.ChapterID)
			}
		}
	}
	catalogue.ChapterBar = New(len(catalogue.ChapterList))
	for _, ChapterID := range catalogue.ChapterList {
		catalogue.CatContent(ChapterID)
	}
	catalogue.ChapterList = nil
	return true
}

func (catalogue *Catalogue) CatContent(ChapterId string) {
	catalogue.DelayTime()
	response := hbooker.GetContent(ChapterId)
	writeContent := fmt.Sprintf("%v:%v\n%v\n\n\n",
		response.ChapterTitle, response.Uptime, response.TxtContent,
	)
	cfg.EncapsulationWrite(catalogue.SaveTextPath, writeContent, 5, "a")
	//catalogue.AddChapterConfig(response.ChapterID)

}
func (catalogue *Catalogue) DelayTime() {
	if err := catalogue.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		time.Sleep(time.Second * time.Duration(rand.Intn(2)))
	}
}
