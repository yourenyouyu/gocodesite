{{define "comment_article_html"}}
<!-- 多说评论框 start -->
	<div class="ds-thread" data-thread-key="{{.article.Id}}" data-title="{{.article.ArticleTitle}}" data-url="{{.url}}"></div>
<!-- 多说评论框 end -->
{{end}}