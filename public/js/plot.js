var dpsTemperature = [{x: 1, y: 0}];
var dpsCPULoad = [{x: 1, y: 0}];
var dpsUsedMemory = [{x: 1, y: 0}];
var xValTemperature = 2;
var xValCPULoad = 2;
var xValUsedMemory = 2;
var chart = null;

function initCharts(historyDatapoints) {
    chartTemperature = new CanvasJS.Chart("chartContainerTemperature",{
        title :{
            text: "Temperature"
        },
        axisX: {
            title: "Timestep",
            interval: 1
        },
        axisY: {
            title: "Temperature [C]",
            minimum: 0,
            maximum: 100
        },
        data: [{
            type: "line",
            dataPoints : dpsTemperature
        }]
    });
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
    chartCPULoad = new CanvasJS.Chart("chartContainerCPULoad",{
        title :{
            text: "CPU Load"
        },
        axisX: {
            title: "Timestep",
            interval: 1
        },
        axisY: {
            title: "CPU Load [%]",
            minimum: 0,
            maximum: 100
        },
        data: [{
            type: "line",
            dataPoints : dpsCPULoad
        }]
    });
    chartUsedMemory = new CanvasJS.Chart("chartContainerUsedMemory",{
        title :{
            text: "Used Memory"
        },       axisX: {
            title: "Timestep",
            interval: 1
        },
        axisY: {
            // Hårdkodad Megabyte: lugnt?
            title: "Used memory [MB]",
            minimum: 0,
            maximum: 5000
        },
        data: [{
            type: "area",
            dataPoints : dpsUsedMemory
        }]
    });
    // Show empty initial chart
    chartTemperature.render();
    chartCPULoad.render();
    chartUsedMemory.render();
    chartTemperatureHistory.render();
}

function updateChart(chartType, newData) {
    // Very strange! newData MUST go through Math library
    if(chartType === "Temperature") {
        dpsTemperature.push({x: xValTemperature, y: Math.round(newData)});
        xValTemperature++;
        if (dpsTemperature.length > 10) {
            dpsTemperature.shift();
        }
        chartTemperature.render();

    }
    if(chartType === "Load") {
        dpsCPULoad.push({x: xValCPULoad, y: Math.round(newData)});
        xValCPULoad++;
        if (dpsCPULoad.length > 10) {
            dpsCPULoad.shift();
        }
        chartCPULoad.render();

    }
    if(chartType === "Used") {
        dpsUsedMemory.push({x: xValUsedMemory, y: Math.round(newData)});
        xValUsedMemory++;
        if (dpsUsedMemory.length > 10) {
            dpsUsedMemory.shift();
        }
        chartUsedMemory.render();

    }
}
