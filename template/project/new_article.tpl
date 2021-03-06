{{define "Title"}}新的文章{{end}}
{{define "importcss"}}
<link href="/static/css/editormd.min.css" rel="stylesheet" />
<link href="/static/css/new_article.css" rel="stylesheet" />
{{end}}
{{define "importjs"}}
<script src="/static/js/editormd.min.js"></script>
<script src="/static/js/editor.js"></script>
<script src="/static/js/new_article.js"></script>
<script type="text/javascript">
	var editor = editormd("editormd", {
		height: 400,
		markdown: null,
		autoFocus: false,
		path: "/static/js/editor.md-1.5.0/lib/",
		//path: "../../../static/js/editor.md-1.5.0/lib/",
		placeholder: "采用markdown语法",
		toolbarIcons: function() {
		  return ["undo", "redo", "|", "bold", "italic", "quote", "|", "h1", "h2", "h3", "h4", "h5", "h6", "|", "list-ul", "list-ol", "hr", "|", "link", "reference-link", "image", "code", "preformatted-text", "code-block", "|", "goto-line", "watch", "preview", "fullscreen", "|", "help", "info"]
		},
		saveHTMLToTextarea: true,
		imageUpload: false,
		imageFormats: ["jpg", "jpeg", "gif", "png"],
		imageUploadURL: "/ajax/article_image_upload?projectId={{.project.Id}}",
		onchange: function() {
		  $("#article-submit").attr('disabled', this.getMarkdown().trim() == "");
		}
	});
</script>
{{end}}
{{define "content"}}
<div id="id-content" class="container">
	<div class="row">
		<div class="col-md-8 col-md-offset-2">
			<div class="breadcrumb">
				<li>
					<a href="/">
						<i class="fa fa-home"></i>首页
					</a>
				</li>
				<li>
					<a href="/project">
						项目
					</a>
				</li>
				<li>
					<a href="/project/{{.project.Id}}/page/1">
						{{.project.ProjectName}}
					</a>
				</li>
			</div>
			<div id="article-tip" class="alert alert-danger hide" role="alert">
				<span id="article-tip-text">ERROR</span>
				<a class="close" data-dismiss="modal" onclick="$('#article-tip').addClass('hide');">×</a>
			</div>
			<div class="reply-container">
				<form id="postarticle-form" action="/ajax/article_submit" method="post" role="form">
					<fieldset>
						<p><span style="color:red;">创建文章暂不支持上传图片，请在编辑文章页面进行图片上传</p>
						<div class="from-group">
							<label for="title">文章标题</label>
							<textarea id="text-title" name="title" class="form-control" rows="1" placeholder="请输入文章标题" 
								onchange="this.value=this.value.substring(0, 64)" 
								onkeydown="this.value=this.value.substring(0, 64)" 
								onkeyup="this.value=this.value.substring(0, 64)"></textarea>
							<input type="hidden" name="projectid" value="{{.project.Id}}" />
						</div>
						<hr/>
						<div class="form-group">
							<label for="title">文章封面</label>
							<input type="cover" class="form-control input-md" placeholder="封面路径" name="coverImage" />
						</div>
						<hr/>
						<div class="form-group">
							<label for="title">文章内容</label>
						</div>
						<div class="form-group">
							<div id="editormd">
								<textarea style="display:none;"></textarea>
							</div>
						</div>
						<hr/>
						<div class="form-group">
						  <div class="input-group">
							<input type="text" id="captchaSolution" name="captchaSolution" placeholder="请输入右侧验证码" />
							<img id="id-article-captchaimg" src="/captcha/{{.captchaid}}.png" alt="验证码" title="看不清，点击" />
							<input type="hidden" id="id-article-captchaIdHolder" name="captchaid" value="{{.captchaid}}">
						  </div>
						</div>
						<div style="text-align:center">
							<!--input type="submit" id="article-submit" class="btn btn-primary" value="提交"/-->
							<a id="article-submit" href="javascript:void(0);" onclick="submitPostArticle(this)" class="btn btn-success">提交</a>
						</div>
					</fieldset>
				</form>
			</div>
		</div>
	</div>
</div>
{{end}}