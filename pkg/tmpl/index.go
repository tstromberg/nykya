package tmpl

import "text/template"

var Index = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html lang="en">
 <head>
  <title>{{ .Title }}</title>
 </head>

<body itemscope itemtype="http://schema.org/Blog">
 <main>
  {{ range .Posts }}
  <article itemprop="blogPosts" itemscope itemtype="http://schema.org/BlogPosting">
   <header>
    <h1 itemprop="headline">Example Image</h1>
   </header>
   <div itemprop="articleBody">
    <img src="{{ .Image }}" />
   </div>
   <footer>
    <p>Posted <time itemprop="datePublished" datetime="2009-10-10">Thursday</time>.</p>
   </footer>
  </article>
  {{ end }}
 <body>
 </body>
</html>
`))
