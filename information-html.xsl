<?xml version="1.0"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform" xmlns="http://www.w3.org/1999/xhtml">

  <xsl:output doctype-public="-//W3C//DTD XHTML 1.0 Strict//EN" encoding="utf-8" indent="yes" media-type="application/html+xml" method="xml" omit-xml-declaration="yes"/>

  <xsl:template match="informations">
    <html>
      <head>
        <meta content="width=device-width, initial-scale=1.0" name="viewport"/>
        <!-- Bootstrap -->
        <link href="public/css/bootstrap.css" media="screen" rel="stylesheet"/>
        <script src="http://code.jquery.com/jquery.js" type="text/javascript"></script>
        <script src="public/js/bootstrap.min.js" type="text/javascript"></script>
        <script src="public/js/jquery.canvasjs.min.js" type="text/javascript"></script>
        <script src="public/js/plot.js" type="text/javascript"></script>
        <title>ServerMonitor</title>
      </head>

      <body>

        <div id="wrap">
          <div class="container">
            <div class="page-header">
              <h1>DM2517 ServerMonitor</h1>
            </div>
            <h3 class="mainheader">Live feed:</h3>
            <div id="changingTable"></div>
            <div id="horizontalContainer">
              <div id="chartContainerTemperature" style="height: 300px; width: 50%; float:left;"></div>
              <div id="chartContainerCPULoad" style="height: 300px; width: 50%;"></div>
            </div>
            <div id="chartContainerUsedMemory" style="height: 300px; width: 100%;"></div>

            <script type="text/javascript">
              window.onload = function () {
                // The chart with history needs some data...
                data = [];
                i = 0;
                // We need {date,temperature} pairs which are in different parts in the xml
                dates = [];
                $("#historyTable").find("tr").find("td#date").each(function () {
                  dates.push(new Date($(this).text()));
                });
                $("#historyTable").find("tr").find("td#tempValue").each(function () {
                  data.push({
                    x: dates[i],
                    y: Math.round($(this).text())
                  });
                  i++;
                });
                initCharts(data);
                // hide
                $("#historyTable").empty();
              }
            </script>

            <div>

              <div class="historyTable" style="float: left;display:none; visibility:hidden;">
                <table class="table table-condensed" id="historyTable">
                  <xsl:apply-templates select="information"/>
                </table>
              </div>
              <div id="chartContainerTemperatureHistory" style="width: 70%; float: left;"></div>
              <table class="table table-bordered" id="funFacts" style="margin-top:40px;float:right; width:30%;">
                <thead>
                  <tr>
                    <th>Min</th>
                    <th>Max</th>
                    <th>Avg</th>
                    <th>Unit</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>

                      <xsl:value-of select="Funfacts/Min"></xsl:value-of>
                    </td>
                    <td>

                      <xsl:value-of select="Funfacts/Max"></xsl:value-of>
                    </td>
                    <td>

                      <xsl:value-of select="Funfacts/Avg"></xsl:value-of>
                    </td>
                    <td>

                      <xsl:value-of select="Funfacts/Unit"></xsl:value-of>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

          </div>
        </div>

        <script>
          var ip = window.location.href.split("/").pop();
          var serverSocket = new WebSocket("ws://localhost:8080/requestdata/" + ip);
          window.setInterval(function () {

            serverSocket.send(ip);
            serverSocket.onmessage = function (e) {
              // New XML data received!
              xml = e.data;

              if (xml === "-1") {console.log("heeej guys");
                $(".mainheader").append(" WARNING! Cannot connect to the ip...")}

              // Remove old data from html
              $("#changingTable").empty();

              // Update with new data:
              $("#changingTable").append("<table></table>");
              $(xml).find("ServerData").each(function () {
                descr = $(this).find("description").text();
                value = $(this).find("value").text();
                unit = $(this).attr("unit");
                var table = $("#changingTable").children();
                table.append("<tr><td>" + descr + "</td><td>" + value + "</td><td> " + unit + " </td></tr>");
                // console.log(descr + " " + value + " " + unit);
                updateChart(descr, value);
              });

              // Add the uptime info
              $(xml).find("Uptime").each(function () {
                console.log("hello its me");
                value = $(this).text();
                var table = $("#changingTable").children();
                table.append("<tr><td>" + "Uptime" + "</td><td>" + value + "</td><td> " + "" + " </td></tr>");
              });
            };
          }, 5000);
        </script>
      </body>
    </html>
  </xsl:template>

  <xsl:template match="information">
    <tr>
      <td id="date">
        <xsl:value-of select="Date"></xsl:value-of>
      </td>
      <xsl:apply-templates select="CPU"/>
    </tr>
  </xsl:template>

  <xsl:template match="CPU">
    <!-- We only care about history of temperature, right now. -->
    <xsl:apply-templates select="ServerData[Description[text() = 'Temperature']]"/>
  </xsl:template>

  <xsl:template match="ServerData[Description[text() = 'Temperature']]">
    <td id="tempValue">
      <xsl:value-of select="value"></xsl:value-of>
    </td>
    <td>
      <xsl:value-of select="Unit"></xsl:value-of>
    </td>
    <td>
      <xsl:value-of select="Description"></xsl:value-of>
    </td>
  </xsl:template>

</xsl:stylesheet>
