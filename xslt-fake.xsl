<?xml version="1.0"?>
<xsl:stylesheet version="1.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns="http://www.w3.org/1999/xhtml">

    <xsl:template match="information">
    <html>
        <head>
            <script src="http://code.jquery.com/jquery-latest.min.js"></script>
            <script src="public/js/socket.js"></script>
            <title>ServerMonitor</title>
        </head>

        <body>
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
