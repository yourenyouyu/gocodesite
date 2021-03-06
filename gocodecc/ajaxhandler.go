package gocodecc

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
)

type AjaxResult struct {
	Result    int    `json:"Result"`
	Msg       string `json:"Msg"`
	CaptchaId string `json:"CaptchaId"`
}

type ArticleImageUploadResult struct {
	Success int    `json:"success"`
	Url     string `json:"url"`
	Message string `json:"message"`
}

func ajaxHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	action := vars["action"]
	var result AjaxResult
	result.Result = -1

	//	for article upload image
	var uploadResult ArticleImageUploadResult

	defer func() {
		if action == "upload" {
			//	need present result
			redirectUrl := ""
			if 0 == result.Result {
				redirectUrl = fmt.Sprintf("/common/message?text=&result=&title=上传成功")
			} else {
				redirectUrl = fmt.Sprintf("/common/message?text=%s&result=1&title=上传失败", result.Msg)
			}
			ctx.Redirect(redirectUrl, http.StatusFound)
		} else if action == "article_submit" ||
			action == "article_edit" {
			if 0 != result.Result {
				//	new captcha
				result.CaptchaId = captcha.NewLen(4)
			}
			renderJson(ctx, &result)
		} else if action == "article_image_upload" {
			ctx.RenderJson(&uploadResult)
		} else {
			renderJson(ctx, &result)
		}
	}()

	switch action {
	case "project_create":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			//	must be superadmin
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "permission denied"
				return
			}

			//	check project name and project describe
			ctx.r.ParseForm()
			defer ctx.r.Body.Close()
			projectName := ctx.r.Form.Get("project[name]")
			projectDescribe := ctx.r.Form.Get("project[describe]")
			projectImage := ctx.r.Form.Get("project[image]")
			//	check with auth
			auth, err := strconv.Atoi(ctx.r.Form.Get("dst"))
			if nil != err {
				result.Msg = "Invalid auth select"
				return
			}
			if auth != kPermission_User &&
				auth != kPermission_SuperAdmin {
				result.Msg = "Invalid auth select"
				return
			}

			if len(projectName) == 0 ||
				len(projectDescribe) == 0 ||
				len(projectName) >= kCategoryNameLimit ||
				len(projectDescribe) >= kCategoryDescribeLimit {
				result.Msg = "invalid project name or project describe"
				return
			}

			var project ProjectCategoryItem
			project.Author = ctx.user.NickName
			project.Image = projectImage
			project.ProjectName = projectName
			project.ProjectDescribe = projectDescribe
			project.PostPriv = uint32(auth)
			err = modelProjectCategoryAdd(&project)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			//	done
			result.Result = 0
		}
	case "project_edit":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			//	must be superadmin
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "permission denied"
				return
			}

			//	check project name and project describe
			ctx.r.ParseForm()
			defer ctx.r.Body.Close()

			var err error
			projectName := ctx.r.Form.Get("project[name]")
			projectDescribe := ctx.r.Form.Get("project[describe]")
			projectImage := ctx.r.Form.Get("project[image]")
			projectId, err := strconv.Atoi(ctx.r.Form.Get("project[id]"))
			//	check with auth
			auth, err := strconv.Atoi(ctx.r.Form.Get("dst"))
			if nil != err {
				result.Msg = "Invalid auth select"
				return
			}
			if auth != kPermission_User &&
				auth != kPermission_SuperAdmin {
				result.Msg = "Invalid auth select"
				return
			}

			if len(projectName) == 0 ||
				len(projectDescribe) == 0 ||
				len(projectName) >= kCategoryNameLimit ||
				len(projectDescribe) >= kCategoryDescribeLimit ||
				nil != err ||
				0 == projectId {
				result.Msg = "invalid project name or project describe"
				return
			}

			//	get the original item
			var originPrj ProjectCategoryItem
			if err := modelProjectCategoryGetByProjectId(projectId, &originPrj); nil != err {
				result.Msg = err.Error()
				return
			}

			if originPrj.ProjectName == projectName &&
				originPrj.ProjectDescribe == projectDescribe &&
				originPrj.Image == projectImage &&
				originPrj.PostPriv == uint32(auth) {
				return
			}

			var newPrj ProjectCategoryItem
			newPrj = originPrj
			newPrj.ProjectName = projectName
			newPrj.ProjectDescribe = projectDescribe
			newPrj.Image = projectImage
			newPrj.PostPriv = uint32(auth)
			err = modelProjectCategoryUpdateProject(&originPrj, &newPrj)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			//	done
			result.Result = 0
		}
	case "project_delete":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			//	must be superadmin
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "permission denied"
				return
			}

			ctx.r.ParseForm()
			projectId, err := strconv.Atoi(ctx.r.Form.Get("project[id]"))
			ctx.r.Body.Close()

			if projectId == 0 ||
				nil != err {
				result.Msg = "invalid project name"
				return
			}

			err = modelProjectCategoryRemove(projectId)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			//	done
			result.Result = 0
		}
	case "article_submit":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			defer ctx.r.Body.Close()
			projectId, err := strconv.Atoi(ctx.r.Form.Get("projectid"))
			if projectId == 0 ||
				nil != err {
				result.Msg = "invalid project"
				return
			}
			//	check captcha
			if !captcha.VerifyString(ctx.r.Form.Get("captchaid"), ctx.r.Form.Get("captchaSolution")) {
				result.Msg = "验证码错误"
				return
			}
			//	check auth
			var prj ProjectCategoryItem
			if err := modelProjectCategoryGetByProjectId(projectId, &prj); nil != err {
				result.Msg = err.Error()
				return
			}
			//	check auth
			if ctx.user.Permission < prj.PostPriv &&
				ctx.user.NickName != prj.Author {
				result.Msg = "permission denied"
				return
			}
			//	check post time
			if ctx.user.Permission < kPermission_Admin {
				lastPostTime := modelProjectArticleGetLastPostTime(ctx.user.UserName)
				tmNow := time.Now().Unix()
				if tmNow-lastPostTime < kMemberPostInterval {
					nextPostTime := lastPostTime + kMemberPostInterval - tmNow
					result.Msg = "离下一次发帖时间还有" + strconv.FormatInt(nextPostTime, 10) + "秒"
					return
				}
			}
			//	check valid
			title := ctx.r.Form.Get("title")
			if len(title) >= kArticleTitleLimit {
				result.Msg = "标题长度太长了"
				return
			}
			if len(title) == 0 {
				result.Msg = "请输入标题"
				return
			}
			contentHtml := ctx.r.Form.Get("editormd-html-code")
			//contentHtml = strings.Replace(contentHtml, "<pre>", `<pre class="prettyprint linenums">`, -1)
			if len(contentHtml) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentHtml) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}
			contentMarkdown := ctx.r.Form.Get("editormd-markdown-doc")
			if len(contentMarkdown) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentMarkdown) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}
			coverImage := ctx.r.Form.Get("coverImage")

			//	do post
			var postArticle ProjectArticleItem
			postArticle.ActiveTime = time.Now().Unix()
			postArticle.PostTime = time.Now().Unix()
			postArticle.ArticleTitle = title
			postArticle.ArticleAuthor = ctx.user.NickName
			postArticle.ArticleContentHtml = contentHtml
			postArticle.ArticleContentMarkdown = contentMarkdown
			postArticle.ProjectName = prj.ProjectName
			postArticle.ProjectId = prj.Id
			postArticle.CoverImage = coverImage
			articleId, err := modelProjectArticleNewArticle(&postArticle)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			//	done
			result.Result = 0
			result.Msg = fmt.Sprintf("/project/%d/article/%d", projectId, articleId)
		}
	case "article_edit":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			defer ctx.r.Body.Close()

			//	check captcha
			if !captcha.VerifyString(ctx.r.Form.Get("captchaid"), ctx.r.Form.Get("captchaSolution")) {
				result.Msg = "验证码错误"
				return
			}
			projectId, _ := strconv.Atoi(ctx.r.Form.Get("projectId"))
			if projectId == 0 {
				result.Msg = "invalid project"
				return
			}
			articleId, _ := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if 0 == articleId {
				result.Msg = "invalid articleId"
				return
			}
			coverImage := ctx.r.Form.Get("coverImage")

			//	check auth
			article, err := modelProjectArticleGet(articleId)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			if article.ArticleAuthor != ctx.user.NickName {
				if ctx.user.Permission < kPermission_SuperAdmin {
					result.Msg = "access denied"
					return
				}
			}
			//	check valid
			title := ctx.r.Form.Get("title")
			if len(title) >= kArticleTitleLimit {
				result.Msg = "标题长度太长了"
				return
			}
			if len(title) == 0 {
				result.Msg = "请输入标题"
				return
			}
			contentHtml := ctx.r.Form.Get("editormd-html-code")
			//contentHtml = strings.Replace(contentHtml, "<pre>", `<pre class="prettyprint linenums">`, -1)
			if len(contentHtml) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentHtml) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}
			contentMarkdown := ctx.r.Form.Get("editormd-markdown-doc")
			if len(contentMarkdown) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentMarkdown) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}

			if len(coverImage) == 0 &&
				len(article.CoverImage) == 0 {
				//	find the first image label and use it
				coverImage = getOneImageFromHtml(contentHtml)
				if len(coverImage) != 0 {
					coverImage = filepath.Base(coverImage)
					extType := strings.ToLower(filepath.Ext(coverImage))
					switch extType {
					case ".jpg":
						fallthrough
					case ".jpeg":
						fallthrough
					case ".png":
						fallthrough
					case ".gif":
						fallthrough
					case ".webp":
						{
							//	nothing
						}
					default:
						{
							extType = ""
						}
					}
					if len(extType) == 0 {
						//	invalid image extension
						coverImage = ""
					}
				}
			}

			//	do post
			colsEdit := []string{"active_time", "edit_time"}
			article.ActiveTime = time.Now().Unix()
			article.EditTime = time.Now().Unix()
			if article.ArticleTitle != title {
				article.ArticleTitle = title
				colsEdit = append(colsEdit, "article_title")
			}
			if article.ArticleContentHtml != contentHtml {
				article.ArticleContentHtml = contentHtml
				article.ArticleContentMarkdown = contentMarkdown
				colsEdit = append(colsEdit, "article_content_html")
				colsEdit = append(colsEdit, "article_content_markdown")
			}
			if article.CoverImage != coverImage {
				article.CoverImage = coverImage
				colsEdit = append(colsEdit, "cover_image")
			}
			_, err = modelProjectArticleEditArticle(article, colsEdit)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			//	done
			result.Result = 0
			result.Msg = fmt.Sprintf("/project/%d/article/%d", projectId, articleId)
		}
	case "article_delete":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			ctx.r.Body.Close()
			if err != nil ||
				0 == articleId {
				result.Msg = "invalid articleId"
				return
			}

			//	get article
			article, err := modelProjectArticleGet(articleId)
			if nil != err {
				result.Msg = "invalid article"
				return
			}

			//	must be superadmin
			if ctx.user.Permission <= kPermission_Admin {
				result.Msg = "access denied"
				return
			}

			err = modelProjectArticleDelete(articleId, article.ProjectId)
			if nil != err {
				result.Msg = "delete article failed"
				return
			}

			//	done
			result.Result = 0
			result.Msg = fmt.Sprintf("/project/%d/page/1", article.ProjectId)
		}
	case "article_top":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			fmt.Println(ctx.r.Form)
			defer ctx.r.Body.Close()
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if err != nil ||
				0 == articleId {
				result.Msg = "invalid articleId"
				return
			}
			top, err := strconv.Atoi(ctx.r.Form.Get("top"))
			if err != nil {
				result.Msg = "invalid top"
				return
			}

			doTop := true
			if 0 == top {
				doTop = false
			}

			err = modelProjectArticleSetTop(articleId, doTop)
			if nil != err {
				result.Msg = "set top failed"
				return
			}

			//	done
			result.Result = 0
		}
	case "article_image_upload":
		{
			ctx.r.ParseForm()
			projectId, err := strconv.Atoi(ctx.r.Form.Get("projectId"))
			if err != nil ||
				0 == projectId {
				uploadResult.Message = "非法的参数"
				return
			}
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if nil != err ||
				0 == articleId {
				uploadResult.Message = "非法的参数"
				return
			}

			//	create directory
			articleImagePath := "." + kPrefixImagePath + "/article-images/" + strconv.Itoa(projectId) + "/" + strconv.Itoa(articleId)
			err = os.MkdirAll(articleImagePath, 0777)
			if nil != err {
				uploadResult.Message = err.Error()
				return
			}

			file, header, err := ctx.r.FormFile("editormd-image-file")
			if nil != err {
				panic(err)
				return
			}
			defer file.Close()

			// 检查是否是jpg或png文件
			uploadFileType := header.Header["Content-Type"][0]

			filenameExtension := ""
			if uploadFileType == "image/jpeg" {
				filenameExtension = ".jpg"
			} else if uploadFileType == "image/png" {
				filenameExtension = ".png"
			} else if uploadFileType == "image/gif" {
				filenameExtension = ".gif"
			}

			if filenameExtension == "" {
				uploadResult.Message = "不支持的文件格式，请上传 jpg/png/gif 图片"
				return
			}

			//	copy to dest directory
			uploadImagePath := articleImagePath + "/" + header.Filename
			f, err := os.OpenFile(uploadImagePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				uploadResult.Message = err.Error()
				return
			}
			defer f.Close()
			io.Copy(f, file)
			uploadResult.Success = 1
			uploadResult.Url = strings.Trim(uploadImagePath, ".")
		}
	case "account_verify":
		{
			if ctx.r.Method != "GET" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			account := ctx.r.Form.Get("account")
			password := ctx.r.Form.Get("password")

			if len(account) == 0 ||
				len(password) == 0 ||
				len(account) > 20 ||
				len(password) > 100 {
				result.Msg = "invalid input"
				return
			}

			user := modelWebUserGetUserByUserName(account)
			if nil == user {
				result.Msg = "user not exists"
				result.Result = -2
				return
			}

			if password != user.PassToken {
				result.Msg = "invalid password"
				result.Result = -3
				return
			}

			//	done
			result.Result = 0
		}
	case "upload":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "Invalid method"
				return
			}
			if ctx.user.Permission < kPermission_SuperAdmin {
				//result.Msg = "access denied"
				result.Msg = kErrMsg_AccessDenied
				return
			}

			//	1 MB
			var fileSizeLimit int64 = 1 * 1024 * 1024
			ctx.r.ParseMultipartForm(fileSizeLimit)
			file, handler, err := ctx.r.FormFile("uploadfile")
			path := strings.Trim(ctx.r.Form.Get("path"), "/")
			path = strings.Trim(path, "\\")
			if len(path) != 0 {
				path += "/"
			}
			if nil != err {
				result.Msg = err.Error()
				return
			}
			defer file.Close()

			fileSize := int64(0)
			if statInterface, ok := file.(FileStat); ok {
				fileInfo, _ := statInterface.Stat()
				fileSize = fileInfo.Size()
			}
			if 0 == fileSize {
				if sizeInterface, ok := file.(FileSize); ok {
					fileSize = sizeInterface.Size()
				}
			}

			if fileSize > fileSizeLimit {
				result.Msg = "文件大小超过限制"
				return
			}

			//	check with path
			pathSel := ctx.r.Form.Get("dst")
			if pathSel == "static" {
				pathSel = "./static/"
			} else if pathSel == "tpl" {
				pathSel = "./template/"
			} else {
				result.Msg = "Invalid file type"
				return
			}

			f, err := os.OpenFile(pathSel+path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				result.Msg = err.Error()
				return
			}
			defer f.Close()
			io.Copy(f, file)
			result.Result = 0
		}
	default:
		{
			result.Msg = "invalid ajax request"
		}
	}
}
