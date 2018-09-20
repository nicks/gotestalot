<!DOCTYPE html>
<html>
  <head>
    <title>Gotestalot &mdash; {{.Package}}</title>
    <link rel="stylesheet" href="/css/main.css">
  </head>

  <body>
    <h1 class="header unknown">
      🕒 Gotestalot: {{.Package}}
    </h1>

    <div id="output" class="output">
      Waiting...
    </div>

    <script src="/js/output.js"></script>
    <script>
      let el = document.querySelector("#output")
      let view = new SummaryView(el)
      view.load("/api/all")
    </script>
  </body>

</html>
