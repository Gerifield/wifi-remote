<!DOCTYPE html>
<html lang="en">
<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">

    <!-- Metro 4 -->
    <link rel="stylesheet" href="https://cdn.metroui.org.ua/v4/css/metro-all.min.css">
</head>
<body>

<div class="tiles-grid">
    {{range .Buttons}}
    <div data-role="tile" data-size="medium" data-button="{{.ID}}" {{if .ColorCode}}style="background-color: {{.ColorCode}}"{{end}}>
        <span class="icon {{.IconClass}}"></span>
    </div>
    {{end}}
</div>

<!-- Metro 4 -->
<script src="https://cdn.metroui.org.ua/v4/js/metro.min.js"></script>
<script>
    const socket = new WebSocket('ws://127.0.0.1:8080/connect');
    // Connection opened
    socket.addEventListener('open', function (event) {
        console.log(event);
    });

    // Listen for messages
    socket.addEventListener('message', function (event) {
        console.log('Message from server ', event.data);
        console.log(event);
    });

    $(".tiles-grid > div").on("click", function(e){
        let data = e.currentTarget.getAttribute('data-button');

        console.log('Button press:', data);
        socket.send(e.target.getAttribute('data-button'));
    });
</script>
</body>
</html>