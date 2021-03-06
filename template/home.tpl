{{define "Title"}}sryan的个人小驿站 分享开发的过程与成果{{end}}
{{define "importcss"}}
<link href="/static/css/home.css" rel="stylesheet" />
<link href="/static/css/articles.css" rel="stylesheet" />
{{end}}
{{define "importjs"}}
<script src="/static/js/home.js"></script>
{{end}}
{{define "content"}}
<!--banner-->
<div class="banner">
	<section class="banner-box">
		<ul class="banner-text">
			<p>To be stronger!!!</p>
		</ul>
	</section>
</div>
<div id="id-content" class="container theme-showcase" role="main">
	<div class="row">
		<div class="col-md-9 col-md-offset-0">
			{{$topArticlesCount := len .topArticles}}
			{{if gt $topArticlesCount 0}}
			<h2 class="section-title-s2 lbcolor-box">置顶的文章</h2>
			<div id="recent-articles" class="articles-container" articleCount="{{len .topArticles}}">
				{{range .topArticles}}
				{{template "article_detail_display" .}}
				{{end}}
			</div>
			{{end}}
			
			<h2 class="section-title-s2 lbcolor-box">最近的文章</h2>
			<div id="articles" class="articles-container" articleCount="{{len .recentArticles}}">
				{{range .recentArticles}}
				{{template "article_detail_display" .}}
				{{end}}
			</div>
			{{$recentArticlesCount := len .recentArticles}}
			{{if gt .articleCount $recentArticlesCount}}
			<div class="white-box shadow-box" style="text-align:center;padding-top:10px;padding-bottom:10px;">
				<button type="button" class="btn btn-link" style="background-color:white; border:1px solid #DDDDDD;"><a href="/articles">...查看更多...</a></button>
			</div>
			{{end}}
		</div>
		<div class="col-md-3 col-md-offset-0">
			<h2 class="section-title-s2 lbcolor-box"><a href="/project">主题目录</a></h2>
			<div class="section-category shadow-box">
				<ul class="posts" style="list-style:none;">
					{{range .category}}
					<li class="post-item">
						<label class="post-item-label-1"><a href="/project/{{.Id}}/page/1">{{.ProjectName}}</a></label>
						<span style="float:right;"><label class="post-item-label-1 badge">{{.ItemCount}}</label></span>
					</li>
					{{end}}
				</ul>
			</div>
			<div style="height:25px;"></div>
			<h2 class="section-title-s2 lbcolor-box">统计</h2>
			<div class="section-statistics shadow-box">
				<ul class="posts" style="list-style:none;">
					<li class="post-item">
						<label class="post-item-label-2">主题数：</label>
						<span style="float:right;"><label class="post-item-label-2 badge">{{.articleCount}}</label></span>
					</li>
					<li class="post-item">
						<label class="post-item-label-2">会员数：</label>
						<span style="float:right;"><label class="post-item-label-2 badge">{{.memberCount}}</label></span>
					</li>
					<li class="post-item">
						<label class="post-item-label-2">建站时间：</label>
						<span style="float:right;"><label class="post-item-label-2 badge">{{formatDate .createSiteTime}}</label></span>
					</li>
				</ul>
			</div>
		</div>
	</div>
</div>
{{end}}