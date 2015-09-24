<?xml version="1.0"?>
<xsl:stylesheet version="1.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns="http://www.w3.org/1999/xhtml">

  <xsl:template match="information">
    <html>
      <head>
        <script src="public/js/socket.js"></script>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <!-- Bootstrap -->
        <link href="css/bootstrap.css" rel="stylesheet" media="screen"/>
        <script type="text/javascript" src="http://code.jquery.com/jquery.js"></script>
        <script type="text/javascript" src="js/bootstrap.min.js"></script>
        <script type="text/javascript" src="js/canvasjs.min.js"></script>
        <title>ServerMonitor</title>
      </head>
      <body>
        <div id="wrap">
          <!-- Begin page content -->
          <div class="container">
            <div class="page-header">
              <h1>DM2517 ServerMonitor</h1>
            </div>
            <p class="lead">Foo</p>
            <div id="chartContainer" style="height: 300px; width: 100%;">
            </div>

            <script type="text/javascript">
              window.onload = function () {

              var dps = [{x: 1, y: 10}, {x: 2, y: 13}, {x: 3, y: 18}, {x: 4, y: 20}, {x: 5, y: 17},{x: 6, y: 10}, {x: 7, y: 13}, {x: 8, y: 18}, {x: 9, y: 20}, {x: 10, y: 17}];   //dataPoints.

              var chart = new CanvasJS.Chart("chartContainer",{
              title :{
              text: "Live Data"
              },
              axisX: {
              title: "Timestep"
              },
              axisY: {
              title: "Temperature"
              },
              data: [{
              type: "line",
              dataPoints : dps
              }]
              });

              chart.render();
              var xVal = dps.length + 1;
              var yVal = 15;
              var updateInterval = 1000;

              var updateChart = function () {


              yVal = yVal +  Math.round(5 + Math.random() *(-5-5));
              dps.push({x: xVal,y: yVal});

              xVal++;
              if (dps.length >  10) {
              dps.shift();
              }

              chart.render();

              // update chart after specified time.

              };

              setInterval(function(){updateChart()}, updateInterval);
              }
            </script>

          </div>
        </div>



        <xsl:apply-templates select="CPU"/>

        <script>
          var ip = window.location.href.split("/").pop();
          var serverSocket = new WebSocket("ws://localhost:8080/requestdata/" +ip);
          window.setInterval(function() {

          serverSocket.send(ip);
          serverSocket.onmessage = function(e) {
          console.log(e.data)
          };
          },5000);
        </script>

      </body>

    </html>
  </xsl:template>



  <xsl:template match="CPU">
    <h2>
      <xsl:apply-templates select="temp"/>
    </h2>
  </xsl:template>

  <xsl:template match="temp">
    <h2>
      <xsl:value-of select="."></xsl:value-of>

      <xsl:value-of select="@unit"></xsl:value-of>
    </h2>
  </xsl:template>

</xsl:stylesheet>
