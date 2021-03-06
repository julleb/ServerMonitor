<?xml version="1.0"?>
<xsl:stylesheet version="1.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns="http://www.w3.org/1999/xhtml">
  <xsl:output method="xml" doctype-public="-//W3C//DTD XHTML 1.0 Strict//EN"
              media-type="application/html+xml" encoding="utf-8" omit-xml-declaration="yes" indent="yes"/>

  <xsl:template match="informations">
    <html>
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <!-- Bootstrap -->
        <link href="public/css/bootstrap.css" rel="stylesheet" media="screen"/>
        <script type="text/javascript" src="http://code.jquery.com/jquery.js"></script>
        <script type="text/javascript" src="public/js/bootstrap.min.js"></script>
        <script type="text/javascript" src="public/js/jquery.canvasjs.min.js"></script>
        <script type="text/javascript" src="public/js/android-plot.js"></script>
        <title>ServerMonitor</title>
      </head>

      <body>

        <div id="wrap">
          <div class="container-fluid">
            <div class="page-header">
              <h1>DM2517 ServerMonitor</h1>
            </div>
            <h3 class="mainheader">Live feed: Since android doesnt support Websockets, We cannot give you live update</h3>
            

            <script type="text/javascript">
              window.onload = function () {
              // The chart with history needs some data...
              data = [];
              
              i = 0;
              // We need {date,temperature} pairs which are in different parts in the xml
              dates = [];
              $("#historyTable").find("tr").find("td#date").each(function() {
              dates.push(new Date($(this).text()));
              });
              $("#historyTable").find("tr").find("td#tempValue").each(function() {
              data.push( {x:dates[i], y: Math.round($(this).text())} );
              i++;
              });
              initCharts(data);
              // hide
              $("#historyTable").empty();
              }
            </script>

            <div>

            <table id="funFacts" class="table table-bordered">
                <thead>
                  <tr>
                    <th>Min</th>
                    <th>Max</th>
                    <th>Avg</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td> <xsl:value-of select="Funfacts/Min"></xsl:value-of> </td>
                    <td> <xsl:value-of select="Funfacts/Max"></xsl:value-of> </td>
                    <td> <xsl:value-of select="Funfacts/Avg"></xsl:value-of> </td>
                  </tr>
                </tbody>
              </table>

              <div class="historyTable" style="float: left; display:none; visibility: hidden;">
                <table id="historyTable" class="table table-condensed">
                  <xsl:apply-templates select="information"/>
                </table>
              </div>
              <div id="chartContainerTemperatureHistory"> </div>
              
            </div>

          </div>
        </div>
      </body>
    </html>
  </xsl:template>


  <xsl:template match="information">
    <tr>
      <td id="date"> <xsl:value-of select="Date"></xsl:value-of> </td>
      <xsl:apply-templates select="CPU"/>
    </tr>
  </xsl:template>

  <xsl:template match="CPU">
    <!-- We only care about history of temperature, right now. -->
    <xsl:apply-templates select="ServerData[Description[text() = 'Temperature']]"/>
  </xsl:template>

  <xsl:template match="ServerData[Description[text() = 'Temperature']]">
    <td id="tempValue"> <xsl:value-of select="value"></xsl:value-of> </td>
    <td> <xsl:value-of select="Unit"></xsl:value-of> </td>
    <td> <xsl:value-of select="Description"></xsl:value-of> </td>
  </xsl:template>

</xsl:stylesheet>
