
function initCharts(historyDatapoints) {
    // This chart is built once (it does not receive more datapoints during runtime).
    chartTemperatureHistory = new CanvasJS.Chart("chartContainerTemperatureHistory",{
        title :{
            text: "Temperature History"
        },
        axisX: {
            title: "Date",
            interval:2,
            valueFormatString: "DD/MM hh:mm",
            labelAngle: -20
        },
        axisY: {
            title: "Temperature [C]",
            minimum: 0,
            maximum: 100
        },
        data: [{
            type: "line",
            dataPoints : historyDatapoints
        }]
    });

    chartTemperatureHistory.render();
}

