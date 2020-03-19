package tmpl

import "text/template"

var Index = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html lang="en">
 <head>
  <title>{{ .Title }}</title>
  <link href="atom.xml" type="application/atom+xml" rel="alternate" title="{{.Title }} feed" />
  <meta name=viewport content="width=device-width, initial-scale=1"/>
  <style>
    h1 {
      font-family: sans-serif;
  }
  
  h2 {
      font-family: sans-serif;
  }
  
  pre {
      overflow-x: auto;
  }
  
  blockquote {
      margin: 10px;
      border-left: 3px solid #eee;
      padding-left: 15px;
      color: #333;
  }
  
  @media screen and (min-width: 15cm) {
      body {
          max-width: 18cm;
          padding-left: 8mm;
          padding-right: 8mm;
      }
  }
  
  hr {
      border: none;
      height: 1px;
      background-color: lightgray;
      margin-top: 5mm;
      margin-bottom: 8mm;
  }
  
  .entrylink {
      display: block;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
  }
  
  .datelink {
      padding-left: 1em;
      color: blue;
      font-style: italic;
      text-decoration: none;
      white-space: nowrap;
  }
  
  .vote {
      float: right;
  }
  </style>
  </head>

<body itemscope itemtype="http://schema.org/Blog">
 <h1>{{ .Title }}</h1>
 <main>
  <ul> 
  {{ range .Posts }}
  <li itemprop="headline">
    <i>{{ .Item.FrontMatter.Posted.Format "2006-01-02" }}</i> &mdash; 
    {{ if eq .Item.FrontMatter.Kind "thought" }}
        {{ .Item.Content }}
    {{ else }}
        <a href="{{ .URL }}">{{ .Item.FrontMatter.Title }}</a>
    {{ end }}
    {{ if eq .Item.FrontMatter.Kind "image" }}
        <a href="{{ .URL }}"><img src="{{ $t := index .Thumbs "100w" }}{{ $t.Path }}"/ ></a>
    {{ end }}
    <!-- {{ .Item.FrontMatter.Kind }}: {{ .Item.Path }}-->
    </li>
  {{ end }}
  </ul>
 <body>
 </body>
</html>
`))
