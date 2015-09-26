<?xml version="1.0"?>
<xsl:stylesheet version="1.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns="http://www.w3.org/1999/xhtml">
  <xsl:output method="xml" doctype-public="-//W3C//DTD XHTML 1.0 Strict//EN"
              media-type="application/html+xml" encoding="utf-8" omit-xml-declaration="yes" indent="yes"/>

  <xsl:template match="information">
    <html>
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <!-- Bootstrap -->
        <link href="public/css/bootstrap.css" rel="stylesheet" media="screen"/>
        <script type="text/javascript" src="http://code.jquery.com/jquery.js"></script>
        <script type="text/javascript" src="public/js/bootstrap.min.js"></script>
        <script type="text/javascript" src="public/js/jquery.canvasjs.min.js"></script>
        <script type="text/javascript" src="public/js/plot.js"></script>
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
            <div id="changingTable">
            </div>
            <div id="horizontalContainer">
              <div id="chartContainerTemperature" style="height: 300px; width: 50%; float:left;"> </div>
              <div id="chartContainerCPULoad" style="height: 300px; width: 50%;"> </div>
            </div>
            <div id="chartContainerUsedMemory" style="height: 300px; width: 100%;">
            </div>

            <script type="text/javascript">
              window.onload = function () {
              initCharts();
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
          // New XML data received!
          xml = e.data;

          // Remove old data from html
          $("#changingTable").empty();

          // Update with new data:
          $("#changingTable").append("<table></table>");
          $(xml).find("ServerData").each(function() {
          descr = $(this).find("description").text();
          value = $(this).find("value").text();
          unit = $(this).attr("unit");
          var table = $("#changingTable").children();
          table.append("<tr><td>"+ descr + "</td><td>"+ value + "</td><td> " + unit + " </td></tr>");
          // console.log(descr + " " + value + " " + unit);
          updateChart(descr, value);
          });
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
