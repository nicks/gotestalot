<!DOCTYPE html>
<html>
  <head>
    <title>Gotestalot &mdash; {{.Package}}</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.0/css/bootstrap.min.css" integrity="sha384-9gVQ4dYFwwWSjIDZnLEWnxCjeSWFphJiwGPXr1jddIhOegiu1FwO5qRGvFXOdJZ4" crossorigin="anonymous">
    <link rel="stylesheet" href="/css/main.css">
  </head>

  <body>
    <h1 class="header unknown">
      🕒 Gotestalot: {{.Package}}
    </h1>

    <div id="output" class="output" data-url="/api/all">
      Waiting...
    </div>

    <script src="http://localhost:8001/index.js"></script>
  </body>

</html>
