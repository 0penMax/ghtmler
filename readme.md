# GHTMLER

Tool for building html pages from parts(components).

How to use:

- create files in /ghtml/ folder with same name what your want html file and extension .ghtml
- use special word (@include) in ghtml file for include html file in output file
- put html code that your want include in /components/path/what/you/want.html
- put all your static files in /static/ folder
- run ghtmler
- enjoy your site in dist folder

example index.ghtml: 
```$xslt
    <!DOCTYPE html>
    <html lang="en">
    
        <head>
            <meta charset="utf-8" />
            <title>example</title>
            <link rel="stylesheet" type="text/css" href="static/css/style.css" />
        </head>
    
        <body>
    
        @include ./components/shared/header.html
    
        @include ./components/index/home.html
    
        @include ./components/index/about.html
    
         @include ./components/index/feature2.html
    
        @include ./components/index/feature.html
    
        @include ./components/index/steps.html
    
        @include ./components/index/faq.html
        
        @include ./components/shared/footer.html
 
        </body>
    </html>
```
project structure:

```
├── ghtml
│   ├── **/*.ghtml
├── component
│   ├── index
│   │   ├── *.html
│   ├── shared
│   │   ├── *.html
├── dist
│   ├── static
│   │   ├── **/*.css
│   │   ├── **/*.js
│   ├── result.html
├── static
│   ├── style.css
├── liveReload
│   ├── temp folder for serve flag
├── runGhtml
└── .gitignore
```

## Flags

 - serve -  realtime rebuild on changes in files.
 - minifycss or mcss - minify css, only for build, ignoring for serve
 - optimizecss or ocss - optimize css, only for build, ignoring for serve
 - omcss - mcss && ocss


