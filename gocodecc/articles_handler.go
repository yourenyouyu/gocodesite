package gocodecc

import (
	"strconv"
)

var articlesRenderTpls = []string{
	"template/articles.tpl",
	"template/component/article_detail_display.tpl",
}

func articlesHandler(ctx *RequestContext) {
	var err error
	//	which page
	ctx.r.ParseForm()
	pageStr := ctx.r.Form.Get("p")

	var page int
	if len(pageStr) == 0 {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if nil != err {
			ctx.RenderMessagePage("错误", err.Error(), false)
			return
		}
	}

	//	get total page
	articleCount, err := modelProjectArticleGetArticleCountAll()
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	articlePerPage := 10
	showPages := 5
	pages := (articleCount + articlePerPage - 1) / articlePerPage
	var articles []*ProjectArticleItem

	if 0 == pages {
		//	never post
		if page != 1 {
			ctx.RenderMessagePage("错误", kErrMsg_InternalError, false)
			return
		}
		articles = make([]*ProjectArticleItem, 0, 1)
	} else {
		if page <= 0 ||
			page > pages {
			ctx.RenderMessagePage("错误", kErrMsg_InternalError, false)
			return
		}

		//	get articles
		articles, err = modelProjectArticleGetRecentArticles(page-1, articlePerPage)
		if nil != err {
			ctx.RenderMessagePage("错误", kErrMsg_InternalError, false)
			return
		}
	}

	tplData := make(map[string]interface{})
	tplData["active"] = "articles"
	tplData["articles"] = articles
	tplData["pages"] = pages
	tplData["page"] = page
	tplData["showPages"] = showPages
	data := renderTemplate(ctx, articlesRenderTpls, tplData)
	ctx.w.Write(data)
}
