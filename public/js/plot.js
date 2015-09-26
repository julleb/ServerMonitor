var dps = [{x: 1, y: 0}];
var xVal = dps.length + 1;
var chart = null;
function initChart() {
    chart = new CanvasJS.Chart("chartContainer",{
        title :{
            text: "Live Data"
        },
        axisX: {
            title: "Timestep",
            interval: 1
        },
        axisY: {
            title: "Temperature",
            minimum: 0,
            maximum: 100
        },
        data: [{
            type: "line",
            dataPoints : dps
        }]
    });
    // Show empty initial chart
    chart.render();
}

function updateChart(newData) {
    // Very strange! newData MUST go through Math library
    dps.push({x: xVal, y: Math.round(newData)});
    xVal++;
    if (dps.length > 10) {
        dps.shift();
    }
    chart.render();
}
