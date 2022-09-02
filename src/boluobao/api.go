package boluobao

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	req "github.com/VeronicaAlexia/pineapple-backups/src/https"
	_struct "github.com/VeronicaAlexia/pineapple-backups/struct"
	"github.com/VeronicaAlexia/pineapple-backups/struct/sfacg_structs"
	"github.com/VeronicaAlexia/pineapple-backups/struct/sfacg_structs/bookshelf"
	url_ "net/url"
	"strconv"
)

func GET_BOOK_INFORMATION(NovelId string) (_struct.Books, error) {
	s := new(sfacg_structs.BookInfo)
	req.Get(new(req.Context).Init("novels/"+NovelId).
		Query("expand", "intro,tags,sysTags,originTotalNeedFireMoney").QueryToString(), s, nil)
	if s.Status.HTTPCode == 200 && s.Data.NovelName != "" {
		return _struct.Books{
			NovelName:  config.RegexpName(s.Data.NovelName),
			NovelID:    strconv.Itoa(s.Data.NovelID),
			NovelCover: s.Data.NovelCover,
			AuthorName: s.Data.AuthorName,
			CharCount:  strconv.Itoa(s.Data.CharCount),
			MarkCount:  strconv.Itoa(s.Data.MarkCount),
			SignStatus: s.Data.SignStatus,
		}, nil
	} else {
		if s.Status.Msg != nil {
			return _struct.Books{}, fmt.Errorf(s.Status.Msg.(string))
		} else {
			return _struct.Books{}, fmt.Errorf("book is not found")
		}
	}

}

func GET_ACCOUNT_INFORMATION() *sfacg_structs.Account {
	return req.Get("user", &sfacg_structs.Account{}, nil).(*sfacg_structs.Account)
}

func GET_BOOK_SHELF_INFORMATION() (map[int][]map[string]string, error) {
	s := new(bookshelf.InfoData)
	bookshelf_info := make(map[int][]map[string]string)
	req.Get(new(req.Context).Init("user/Pockets").Query("expand", "novels").QueryToString(), s, nil)
	if s.Status.HTTPCode != 200 {
		return nil, fmt.Errorf(s.Status.Msg.(string))
	}
	for index, value := range s.Data {
		fmt.Println("bookshelf index:", index, "\t\t\tbookshelf name:", value.Name)
		var bookshelf_info_list []map[string]string
		for _, book := range value.Expand.Novels {
			bookshelf_info_list = append(bookshelf_info_list,
				map[string]string{"novel_name": book.NovelName, "novel_id": strconv.Itoa(book.NovelID)},
			)
		}
		bookshelf_info[index] = bookshelf_info_list
	}
	return bookshelf_info, nil
}

func GET_CATALOGUE(NovelID string) []map[string]string {
	var chapter_index int
	var division_info []map[string]string
	s := new(sfacg_structs.Catalogue)
	req.Get(new(req.Context).Init(fmt.Sprintf("novels/%v/dirs", NovelID)).
		Query("expand", "originNeedFireMoney").QueryToString(), s, nil)
	for division_index, division := range s.Data.VolumeList {
		fmt.Printf("第%v卷\t\t%v\n", division_index+1, division.Title)
		for _, chapter := range division.ChapterList {
			chapter_index += 1
			division_info = append(division_info, map[string]string{
				"division_name":  division.Title,
				"division_id":    strconv.Itoa(division.VolumeID),
				"division_index": strconv.Itoa(division_index),
				"chapter_name":   chapter.Title,
				"chapter_id":     strconv.Itoa(chapter.ChapID),
				"chapter_index":  strconv.Itoa(chapter_index),
				"money":          strconv.Itoa(chapter.OriginNeedFireMoney),
				"file_name":      config.FileCacheName(division_index, chapter_index, strconv.Itoa(chapter.ChapID)),
			})
		}
	}
	return division_info

}

func GET_CHAPTER_CONTENT(chapter_id string) string {
	s := new(sfacg_structs.Content)
	req.Get(new(req.Context).Init("Chaps/"+chapter_id).Query("expand", "content").QueryToString(), s, nil)
	if s != nil && s.Status.HTTPCode == 200 {
		content_title := fmt.Sprintf("%v: %v", s.Data.Title, s.Data.AddTime)
		return content_title + "\n" + config.StandardContent(s.Data.Expand.Content)

	} else {
		fmt.Println("download failed! chapterId:", chapter_id, "error:", s.Status.Msg)
	}
	return ""
}

func GET_SEARCH(keyword string, page int) *sfacg_structs.Search {
	s := new(sfacg_structs.Search)
	req.Get(new(req.Context).Init("search/novels/result").Query("q", url_.QueryEscape(keyword)).
		Query("size", "20").Query("page", strconv.Itoa(page)).QueryToString(), s, nil)
	return s

}

func LOGIN_ACCOUNT(username, password string) *sfacg_structs.Login {
	//s := new(sfacg_structs.Login)
	//req.Get(new(req.Context).Init("sessions").Query("username", username).
	//	Query("password", password).QueryToString(), s, nil)
	params := fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password)
	response, Cookie := req.Login(req.SET_URL("sessions", nil), []byte(params))
	for _, cookie := range Cookie {
		response.Cookie += cookie.Name + "=" + cookie.Value + ";"
	}
	return response
}
