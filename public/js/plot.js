function namn() {
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
