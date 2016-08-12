package gocodecc

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/cihub/seelog"

	"github.com/gorilla/mux"
)

var projectCategoryRenderTpls = []string{
	"template/project/category.tpl",
}

var projectArticlesRenderTpls = []string{
	"template/project/articles.tpl",
}

var projectArticleNewArticleTpls = []string{
	"template/project/new_article.tpl",
}

var projectArticleRenderTpls = []string{
	"template/project/article.tpl",
	"template/component/comment.tpl",
}

var projectArticleEditArticleRenderTpls = []string{
	"template/project/edit_article.tpl",
}

func projectHandler(ctx *RequestContext) {
	ctx.Redirect("/project/category", http.StatusFound)
}

func projectCategoryHandler(ctx *RequestContext) {
	//	search all project
	projects, err := modelProjectCategoryGetAll()
	if nil != err {
		panic(err)
	}

	tplData := make(map[string]interface{})
	tplData["category"] = projects
	tplData["active"] = "project"
	data := renderTemplate(ctx, projectCategoryRenderTpls, tplData)
	ctx.w.Write(data)
}

func projectArticlesHandler(ctx *RequestContext) {
	//	search all project
	vars := mux.Vars(ctx.r)
	projectName := vars["projectname"]
	page, err := strconv.Atoi(vars["page"])
	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}
	if page <= 0 {
		page = 1
	}

	pageItems := 10
	showPages := 5
	articles, pages, err := modelProjectArticleGetArticles(projectName, page-1, pageItems)
	if nil != err {
		panic(err)
	}

	tplData := make(map[string]interface{})
	tplData["articles"] = articles
	tplData["active"] = "project"
	tplData["project"] = projectName
	tplData["pages"] = pages
	tplData["page"] = page
	tplData["pageItems"] = pageItems
	tplData["showPages"] = showPages
	data := renderTemplate(ctx, projectArticlesRenderTpls, tplData)
	ctx.w.Write(data)
}

func projectArticleHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	articleId, err := strconv.Atoi(vars["articleid"])

	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	article, err := modelProjectArticleGet(articleId)
	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	//	get author
	author := modelWebUserGetUserByUserName(article.ArticleAuthor)
	if nil == author {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	//	increase click count
	if err = modelProjectArticleIncClick(articleId); nil != err {
		seelog.Error(err)
		return
	}
	article.Click = article.Click + 1

	tplData := make(map[string]interface{})
	tplData["active"] = "project"
	tplData["article"] = article
	tplData["author"] = author
	data := renderTemplate(ctx, projectArticleRenderTpls, tplData)
	ctx.w.Write(data)
}

func projectArticleCmdHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	cmd := vars["cmd"]
	project := vars["projectname"]
	cmd = strings.ToLower(cmd)

	switch cmd {
	case "new_article":
		{
			_newProjectArticle(ctx, project)
		}
	case "edit_article":
		{
			ctx.r.ParseForm()
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if nil != err {
				ctx.Redirect("/", http.StatusNotFound)
				return
			}

			_editProjectArticle(ctx, articleId)
		}
	default:
		{
			ctx.RenderString("invalid cmd")
		}
	}
}

func _newProjectArticle(ctx *RequestContext, project string) {
	tplData := make(map[string]interface{})
	tplData["active"] = "project"
	tplData["project"] = project
	data := renderTemplate(ctx, projectArticleNewArticleTpls, tplData)
	ctx.w.Write(data)
}

func _editProjectArticle(ctx *RequestContext, articleId int) {
	article, err := modelProjectArticleGet(articleId)
	if err != nil {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	tplData := make(map[string]interface{})
	tplData["active"] = "project"
	tplData["article"] = article
	data := renderTemplate(ctx, projectArticleEditArticleRenderTpls, tplData)
	ctx.w.Write(data)
}