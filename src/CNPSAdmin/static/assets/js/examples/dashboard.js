'use strict';
var current_url2 = window.location.origin;

$('.chart_filter .filter-ico').click(function() {

    var p = $(this).parent();
    if (p.find('.filter-box').is(":visible") == true) {
        p.find('.filter-box').hide();
        $(this).find('i').removeClass('fa-minus');
        $(this).find('i').addClass('fa-plus');
    } else {
        p.find('.filter-box').fadeIn(1);
        $(this).find('i').removeClass('fa-plus');
        $(this).find('i').addClass('fa-minus');
    }
});

$('.close').click(function() { //button reset event click
    $(this).parents(".filter-box").hide();
    $(this).parents('.chart_filter').find('.filter-ico i').removeClass('fa-minus');
    $(this).parents('.chart_filter').find('.filter-ico i').addClass('fa-plus');

});
$("#transaction_status_chart_reset, #transaction_chart_reset").click(function() {
    var language = $("#language").val();
    // console.log(language);
    var TodayLabel, YesterdayLabel, thisMonthLabel, lastMonthLabel, last7DaysLabel, last30DaysLabel, customRangeLabel, applyLabel, cancelLabel, daysOfWeek, monthNames;
    if (language == "English") {
        TodayLabel = "Today";
        YesterdayLabel = "Yesterday";
        thisMonthLabel = "This Month";
        lastMonthLabel = "Last Month";
        last7DaysLabel = "Last 7 Days";
        last30DaysLabel = "Last 30 Days";
        customRangeLabel = "Custom Range";
        applyLabel = "Apply";
        cancelLabel = "Cancel";
        daysOfWeek = ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"];
        monthNames = [
            "January",
            "February",
            "March",
            "April",
            "May",
            "June",
            "July",
            "August",
            "September",
            "October",
            "November",
            "December"
        ];
    } else {
        TodayLabel = "Aujourd'hui";
        YesterdayLabel = "Hier";
        thisMonthLabel = "Ce mois-ci";
        lastMonthLabel = "Le mois dernier";
        last7DaysLabel = "Les 7 derniers jours";
        last30DaysLabel = "Les 30 derniers jours";
        customRangeLabel = "Période";
        applyLabel = "Valider";
        cancelLabel = "Annuler";
        daysOfWeek = ["Do", "Lu", "Ma", "Me", "Je", "Ve", "Sa"];
        monthNames = [
            "Janvier",
            "Février",
            "Mars",
            "Avril",
            "Mai",
            "Juin",
            "Juillet",
            "Août",
            "Septembre",
            "Octobre",
            "Novembre",
            "Décembre"
        ];
    }

    $('.simple-date-range-picker').daterangepicker("refresh");
    $('.simple-date-range-picker').daterangepicker({
        autoUpdateInput: true,
        startDate: moment().subtract(4, 'months'),
        ranges: {
            [TodayLabel]: [moment(), moment()],
            [YesterdayLabel]: [moment().subtract(1, 'days'), moment().subtract(1, 'days')],
            [last7DaysLabel]: [moment().subtract(6, 'days'), moment()],
            [last30DaysLabel]: [moment().subtract(29, 'days'), moment()],
            [thisMonthLabel]: [moment().startOf('month'), moment().endOf('month')],
            [lastMonthLabel]: [moment().subtract(1, 'month').startOf('month'), moment().subtract(1, 'month').endOf('month')]
        },
        locale: {
            customRangeLabel: customRangeLabel,
            applyLabel: applyLabel,
            cancelLabel: cancelLabel,
            format: "YYYY/MM/DD",
            daysOfWeek: daysOfWeek,
            monthNames: monthNames
        },
    });
});

$(document).ready(function() {

    var language = $("#language").val();

    var apexDefaultLocale, apexLocales;
    var amountLabel = "";
    if (language == "English") {
        amountLabel = "Amount(XOF)";
        apexDefaultLocale = "en"
        apexLocales = {
            name: 'en',
            "options": {
                "months": [
                    "January",
                    "February",
                    "March",
                    "April",
                    "May",
                    "June",
                    "July",
                    "August",
                    "September",
                    "October",
                    "November",
                    "December"
                ],
                "shortMonths": [
                    "Jan",
                    "Feb",
                    "Mar",
                    "Apr",
                    "May",
                    "Jun",
                    "Jul",
                    "Aug",
                    "Sep",
                    "Oct",
                    "Nov",
                    "Dec"
                ],
                "days": [
                    "Sunday",
                    "Monday",
                    "Tuesday",
                    "Wednesday",
                    "Thursday",
                    "Friday",
                    "Saturday"
                ],
                "shortDays": ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
                "toolbar": {
                    "exportToSVG": "Download SVG",
                    "exportToPNG": "Download PNG",
                    "exportToCSV": "Download CSV",
                    "menu": "Menu",
                    "selection": "Selection",
                    "selectionZoom": "Selection Zoom",
                    "zoomIn": "Zoom In",
                    "zoomOut": "Zoom Out",
                    "pan": "Panning",
                    "reset": "Reset Zoom"
                }
            }
        }
    } else {
        amountLabel = "Montant(XOF)";
        apexDefaultLocale = "fr"
        apexLocales = {
            name: 'fr',
            "options": {
                "months": [
                    "janvier",
                    "février",
                    "mars",
                    "avril",
                    "mai",
                    "juin",
                    "juillet",
                    "août",
                    "septembre",
                    "octobre",
                    "novembre",
                    "décembre"
                ],
                "shortMonths": [
                    "janv.",
                    "févr.",
                    "mars",
                    "avr.",
                    "mai",
                    "juin",
                    "juill.",
                    "août",
                    "sept.",
                    "oct.",
                    "nov.",
                    "déc."
                ],
                "days": [
                    "dimanche",
                    "lundi",
                    "mardi",
                    "mercredi",
                    "jeudi",
                    "vendredi",
                    "samedi"
                ],
                "shortDays": ["dim.", "lun.", "mar.", "mer.", "jeu.", "ven.", "sam."],
                "toolbar": {
                    "exportToSVG": "Télécharger au format SVG",
                    "exportToPNG": "Télécharger au format PNG",
                    "exportToCSV": "Télécharger au format CSV",
                    "menu": "Menu",
                    "selection": "Sélection",
                    "selectionZoom": "Sélection et zoom",
                    "zoomIn": "Zoomer",
                    "zoomOut": "Dézoomer",
                    "pan": "Navigation",
                    "reset": "Réinitialiser le zoom"
                }
            }
        }
    }



    var colors = {
        primary: $('.colors .bg-primary').css('background-color'),
        primaryLight: $('.colors .bg-primary-bright').css('background-color'),
        secondary: $('.colors .bg-secondary').css('background-color'),
        secondaryLight: $('.colors .bg-secondary-bright').css('background-color'),
        info: $('.colors .bg-info').css('background-color'),
        infoLight: $('.colors .bg-info-bright').css('background-color'),
        success: $('.colors .bg-success').css('background-color'),
        successLight: $('.colors .bg-success-bright').css('background-color'),
        danger: $('.colors .bg-danger').css('background-color'),
        dangerLight: $('.colors .bg-danger-bright').css('background-color'),
        warning: $('.colors .bg-warning').css('background-color'),
        warningLight: $('.colors .bg-warning-bright').css('background-color'),
    };

    /**
     *  Slick slide example
     **/

    if ($('.slick-single-item').length) {
        $('.slick-single-item').slick({
            autoplay: true,
            autoplaySpeed: 3000,
            infinite: true,
            slidesToShow: 4,
            slidesToScroll: 4,
            prevArrow: '.slick-single-arrows a:eq(0)',
            nextArrow: '.slick-single-arrows a:eq(1)',
            responsive: [{
                    breakpoint: 1300,
                    settings: {
                        slidesToShow: 3,
                        slidesToScroll: 3,
                    }
                },
                {
                    breakpoint: 992,
                    settings: {
                        slidesToShow: 3,
                        slidesToScroll: 3,
                    }
                },
                {
                    breakpoint: 768,
                    settings: {
                        slidesToShow: 2,
                        slidesToScroll: 2
                    }
                },
                {
                    breakpoint: 540,
                    settings: {
                        slidesToShow: 1,
                        slidesToScroll: 1
                    }
                }
            ]
        });
    }

    if ($('.reportrange').length > 0) {
        var start = moment().subtract(29, 'days');
        var end = moment();

        function cb(start, end) {
            $('.reportrange span').html(start.format('MMMM D, YYYY') + ' - ' + end.format('MMMM D, YYYY'));
        }

        $('.reportrange').daterangepicker({
            startDate: start,
            endDate: end,
            ranges: {
                'Today': [moment(), moment()],
                'Yesterday': [moment().subtract(1, 'days'), moment().subtract(1, 'days')],
                'Last 7 Days': [moment().subtract(6, 'days'), moment()],
                'Last 30 Days': [moment().subtract(29, 'days'), moment()],
                'This Month': [moment().startOf('month'), moment().endOf('month')],
                'Last Month': [moment().subtract(1, 'month').startOf('month'), moment().subtract(1, 'month').endOf('month')]
            }
        }, cb);

        cb(start, end);
    }

    var chartColors = {
        primary: {
            base: '#3f51b5',
            light: '#c0c5e4'
        },
        danger: {
            base: '#f2125e',
            light: '#fcd0df'
        },
        success: {
            base: '#0acf97',
            light: '#cef5ea'
        },
        warning: {
            base: '#ff8300',
            light: '#ffe6cc'
        },
        info: {
            base: '#00bcd4',
            light: '#e1efff'
        },
        dark: '#37474f',
        facebook: '#3b5998',
        twitter: '#55acee',
        linkedin: '#0077b5',
        instagram: '#517fa4',
        whatsapp: '#25D366',
        dribbble: '#ea4c89',
        google: '#DB4437',
        borderColor: '#e8e8e8',
        fontColor: '#999'
    };

    if ($('body').hasClass('dark')) {
        chartColors.borderColor = 'rgba(255, 255, 255, .1)';
        chartColors.fontColor = 'rgba(255, 255, 255, .4)';
    }

    /// Chartssssss

    chart_demo_1();

    chart_demo_2();

    chart_demo_3();

    chart_demo_4();

    chart_demo_5();

    chart_demo_6();

    chart_demo_7();

    chart_demo_8();

    chart_demo_9();

    chart_demo_10();

    function chart_demo_1() {
        if ($('#chart_demo_1').length) {
            var element = document.getElementById("chart_demo_1");
            element.height = 146;
            new Chart(element, {
                type: 'bar',
                data: {
                    labels: ["2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019"],
                    datasets: [{
                        label: "Total Sales",
                        backgroundColor: colors.primary,
                        data: [133, 221, 783, 978, 214, 421, 211, 577]
                    }, {
                        label: "Average",
                        backgroundColor: colors.info,
                        data: [408, 947, 675, 734, 325, 672, 632, 213]
                    }]
                },
                options: {
                    legend: {
                        display: false
                    },
                    scales: {
                        xAxes: [{
                            ticks: {
                                fontSize: 11,
                                fontColor: chartColors.fontColor
                            },
                            gridLines: {
                                display: false,
                            }
                        }],
                        yAxes: [{
                            ticks: {
                                fontSize: 11,
                                fontColor: chartColors.fontColor
                            },
                            gridLines: {
                                color: chartColors.borderColor
                            }
                        }],
                    }
                }
            })
        }
    }

    function chart_demo_2() {
        if ($('#chart_demo_2').length) {
            var ctx = document.getElementById('chart_demo_2').getContext('2d');
            new Chart(ctx, {
                type: 'line',
                data: {
                    labels: ["Jun 2016", "Jul 2016", "Aug 2016", "Sep 2016", "Oct 2016", "Nov 2016", "Dec 2016", "Jan 2017", "Feb 2017", "Mar 2017", "Apr 2017", "May 2017"],
                    datasets: [{
                        label: "Rainfall",
                        backgroundColor: chartColors.primary.light,
                        borderColor: chartColors.primary.base,
                        data: [26.4, 39.8, 66.8, 66.4, 40.6, 55.2, 77.4, 69.8, 57.8, 76, 110.8, 142.6],
                    }]
                },
                options: {
                    legend: {
                        display: false,
                        labels: {
                            fontColor: chartColors.fontColor
                        }
                    },
                    title: {
                        display: true,
                        text: 'Precipitation in Toronto',
                        fontColor: chartColors.fontColor,
                    },
                    scales: {
                        yAxes: [{
                            gridLines: {
                                color: chartColors.borderColor
                            },
                            ticks: {
                                fontColor: chartColors.fontColor,
                                beginAtZero: true
                            },
                            scaleLabel: {
                                display: true,
                                labelString: 'Precipitation in mm',
                                fontColor: chartColors.fontColor,
                            }
                        }],
                        xAxes: [{
                            gridLines: {
                                color: chartColors.borderColor
                            },
                            ticks: {
                                fontColor: chartColors.fontColor,
                                beginAtZero: true
                            }
                        }]
                    }
                }
            });
        }
    }

    function chart_demo_3() {
        if ($('#chart_demo_3').length) {
            var element = document.getElementById("chart_demo_3"),
                ctx = element.getContext("2d");


            new Chart(ctx, {
                type: 'line',
                data: {
                    labels: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
                    datasets: [{
                        label: 'Success',
                        borderColor: colors.success,
                        data: [-10, 30, -20, 0, 25, 44, 30, 15, 20, 10, 5, -5],
                        pointRadius: 5,
                        pointHoverRadius: 7,
                        borderDash: [2, 2],
                        fill: false
                    }, {
                        label: 'Return',
                        fill: false,
                        borderDash: [2, 2],
                        borderColor: colors.danger,
                        data: [20, 0, 22, 39, -10, 19, -7, 0, 15, 0, -10, 5],
                        pointRadius: 5,
                        pointHoverRadius: 7
                    }]
                },
                options: {
                    responsive: true,
                    legend: {
                        display: false,
                        labels: {
                            fontColor: chartColors.fontColor
                        }
                    },
                    title: {
                        display: false,
                        fontColor: chartColors.fontColor
                    },
                    scales: {
                        xAxes: [{
                            gridLines: {
                                display: false,
                                color: chartColors.borderColor
                            },
                            ticks: {
                                fontColor: chartColors.fontColor,
                                display: false
                            }
                        }],
                        yAxes: [{
                            gridLines: {
                                color: chartColors.borderColor
                            },
                            ticks: {
                                fontColor: chartColors.fontColor,
                                min: -50,
                                max: 50
                            }
                        }],
                    }
                }
            });

        }
    }

    function chart_demo_4() {
        if ($('#chart_demo_4').length) {
            var ctx = document.getElementById("chart_demo_4").getContext("2d");
            var densityData = {
                backgroundColor: chartColors.primary.light,
                data: [10, 20, 40, 60, 80, 40, 60, 80, 40, 80, 20, 59]
            };
            new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
                    datasets: [densityData]
                },
                options: {
                    scaleFontColor: "#FFFFFF",
                    legend: {
                        display: false,
                        labels: {
                            fontColor: chartColors.fontColor
                        }
                    },
                    scales: {
                        xAxes: [{
                            gridLines: {
                                color: chartColors.borderColor
                            },
                            ticks: {
                                fontColor: chartColors.fontColor
                            }
                        }],
                        yAxes: [{
                            gridLines: {
                                color: chartColors.borderColor
                            },
                            ticks: {
                                fontColor: chartColors.fontColor,
                                min: 0,
                                max: 100,
                                beginAtZero: true
                            }
                        }]
                    }
                }
            });
        }
    }

    function chart_demo_5() {
        if ($('#chart_demo_5').length) {
            var ctx = document.getElementById('chart_demo_5').getContext('2d');
            window.myBar = new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: ['January', 'February', 'March', 'April', 'May'],
                    datasets: [{
                            label: 'Dataset 1',
                            backgroundColor: [
                                chartColors.info.base,
                                chartColors.success.base,
                                chartColors.danger.base,
                                chartColors.dark,
                                chartColors.warning.base,
                            ],
                            yAxisID: 'y-axis-1',
                            data: [33, 56, -40, 25, 45]
                        },
                        {
                            label: 'Dataset 2',
                            backgroundColor: chartColors.info.base,
                            yAxisID: 'y-axis-2',
                            data: [23, 86, -40, 5, 45]
                        }
                    ]
                },
                options: {
                    legend: {
                        labels: {
                            fontColor: chartColors.fontColor
                        }
                    },
                    responsive: true,
                    title: {
                        display: true,
                        text: 'Chart.js Bar Chart - Multi Axis',
                        fontColor: chartColors.fontColor
                    },
                    tooltips: {
                        mode: 'index',
                        intersect: true
                    },
                    scales: {
                        xAxes: [{
                            gridLines: {
                                color: chartColors.borderColor
                            },
                            ticks: {
                                fontColor: chartColors.fontColor
                            }
                        }],
                        yAxes: [{
                                type: 'linear',
                                display: true,
                                position: 'left',
                                id: 'y-axis-1',
                            },
                            {
                                gridLines: {
                                    color: chartColors.borderColor
                                },
                                ticks: {
                                    fontColor: chartColors.fontColor
                                }
                            },
                            {
                                type: 'linear',
                                display: true,
                                position: 'right',
                                id: 'y-axis-2',
                                gridLines: {
                                    drawOnChartArea: false
                                },
                                ticks: {
                                    fontColor: chartColors.fontColor
                                }
                            }
                        ],
                    }
                }
            });
        }
    }

    function chart_demo_6() {
        if ($('#chart_demo_6').length) {
            var ctx = document.getElementById("chart_demo_6").getContext("2d");
            var speedData = {
                labels: ["0s", "10s", "20s", "30s", "40s", "50s", "60s"],
                datasets: [{
                    label: "Car Speed (mph)",
                    borderColor: chartColors.primary.base,
                    backgroundColor: 'rgba(0, 0, 0, 0',
                    data: [0, 59, 75, 20, 20, 55, 40]
                }]
            };
            var chartOptions = {
                legend: {
                    scaleFontColor: "#FFFFFF",
                    position: 'top',
                    labels: {
                        fontColor: chartColors.fontColor
                    }
                },
                scales: {
                    xAxes: [{
                        gridLines: {
                            color: chartColors.borderColor
                        },
                        ticks: {
                            fontColor: chartColors.fontColor
                        }
                    }],
                    yAxes: [{
                        gridLines: {
                            color: chartColors.borderColor
                        },
                        ticks: {
                            fontColor: chartColors.fontColor
                        }
                    }]
                }
            };
            new Chart(ctx, {
                type: 'line',
                data: speedData,
                options: chartOptions
            });
        }
    }

    function chart_demo_7() {
        if ($('#chart_demo_7').length) {
            var ctx = document.getElementById("chart_demo_7").getContext("2d");
            new Chart(ctx, {
                type: 'doughnut',
                data: {
                    datasets: [{
                        data: [15, 25, 10, 30],
                        backgroundColor: [
                            colors.success,
                            colors.danger,
                            colors.warning,
                            colors.info
                        ],
                        label: 'Dataset 1'
                    }],
                    labels: [
                        'Social Media',
                        'Organic Search',
                        'Referrral',
                        'Email'
                    ]
                },
                options: {
                    elements: {
                        arc: {
                            borderWidth: 0
                        }
                    },
                    responsive: true,
                    legend: {
                        display: false
                    },
                    title: {
                        display: false
                    },
                    animation: {
                        animateScale: true,
                        animateRotate: true
                    }
                }
            });
        }
    }

    function chart_demo_8() {
        if ($('#chart_demo_8').length) {
            new Chart(document.getElementById("chart_demo_8"), {
                type: 'radar',
                data: {
                    labels: ["Africa", "Asia", "Europe", "Latin America", "North America"],
                    datasets: [{
                        label: "1950",
                        fill: true,
                        backgroundColor: "rgba(179,181,198,0.2)",
                        borderColor: "rgba(179,181,198,1)",
                        pointBorderColor: "#fff",
                        pointBackgroundColor: "rgba(179,181,198,1)",
                        data: [-8.77, -55.61, 21.69, 6.62, 6.82]
                    }, {
                        label: "2050",
                        fill: true,
                        backgroundColor: "rgba(255,99,132,0.2)",
                        borderColor: "rgba(255,99,132,1)",
                        pointBorderColor: "#fff",
                        pointBackgroundColor: "rgba(255,99,132,1)",
                        data: [-25.48, 54.16, 7.61, 8.06, 4.45]
                    }]
                },
                options: {
                    legend: {
                        labels: {
                            fontColor: chartColors.fontColor
                        }
                    },
                    scale: {
                        gridLines: {
                            color: chartColors.borderColor
                        }
                    },
                    title: {
                        display: true,
                        text: 'Distribution in % of world population',
                        fontColor: chartColors.fontColor
                    }
                }
            });
        }
    }

    function chart_demo_9() {
        if ($('#chart_demo_9').length) {
            new Chart(document.getElementById("chart_demo_9"), {
                type: 'horizontalBar',
                data: {
                    labels: ["Africa", "Asia", "Europe", "Latin America", "North America"],
                    datasets: [{
                        label: "Population (millions)",
                        backgroundColor: colors.primary,
                        data: [2478, 2267, 734, 1284, 1933]
                    }]
                },
                options: {
                    legend: {
                        display: false
                    },
                    scales: {
                        xAxes: [{
                            gridLines: {
                                color: chartColors.borderColor
                            },
                            ticks: {
                                fontColor: chartColors.fontColor,
                                display: false
                            }
                        }],
                        yAxes: [{
                            gridLines: {
                                color: chartColors.borderColor,
                                display: false
                            },
                            ticks: {
                                fontColor: chartColors.fontColor
                            },
                            barPercentage: 0.5
                        }]
                    }
                }
            });
        }
    }

    function chart_demo_10() {
        if ($('#chart_demo_10').length) {
            var element = document.getElementById("chart_demo_10");
            new Chart(element, {
                type: 'bar',
                data: {
                    labels: ["1900", "1950", "1999", "2050"],
                    datasets: [{
                            label: "Europe",
                            type: "line",
                            borderColor: "#8e5ea2",
                            data: [408, 547, 675, 734],
                            fill: false
                        },
                        {
                            label: "Africa",
                            type: "line",
                            borderColor: "#3e95cd",
                            data: [133, 221, 783, 2478],
                            fill: false
                        },
                        {
                            label: "Europe",
                            type: "bar",
                            backgroundColor: chartColors.primary.base,
                            data: [408, 547, 675, 734],
                        },
                        {
                            label: "Africa",
                            type: "bar",
                            backgroundColor: chartColors.primary.light,
                            data: [133, 221, 783, 2478]
                        }
                    ]
                },
                options: {
                    title: {
                        display: true,
                        text: 'Population growth (millions): Europe & Africa',
                        fontColor: chartColors.fontColor
                    },
                    legend: {
                        display: true,
                        labels: {
                            fontColor: chartColors.fontColor
                        }
                    },
                    scales: {
                        xAxes: [{
                            gridLines: {
                                color: chartColors.borderColor
                            },
                            ticks: {
                                fontColor: chartColors.fontColor
                            }
                        }],
                        yAxes: [{
                            gridLines: {
                                color: chartColors.borderColor
                            },
                            ticks: {
                                fontColor: chartColors.fontColor
                            }
                        }]
                    }
                }
            });
        }
    }

    if ($('#circle-1').length) {
        $('#circle-1').circleProgress({
            startAngle: 1.55,
            value: 0.65,
            size: 110,
            thickness: 10,
            fill: {
                color: colors.primary
            }
        });
    }

    if ($('#sales-circle-graphic').length) {
        $('#sales-circle-graphic').circleProgress({
            startAngle: 1.55,
            value: 0.65,
            size: 180,
            thickness: 30,
            fill: {
                color: colors.primary
            }
        });
    }

    if ($('#circle-2').length) {
        $('#circle-2').circleProgress({
            startAngle: 1.55,
            value: 0.35,
            size: 110,
            thickness: 10,
            fill: {
                color: colors.success
            }
        });
    }

    ////////////////////////////////////////////

    if ($(".dashboard-pie-1").length) {
        $(".dashboard-pie-1").peity("pie", {
            fill: [colors.primaryLight, colors.primary],
            radius: 30
        });
    }

    if ($(".dashboard-pie-2").length) {
        $(".dashboard-pie-2").peity("pie", {
            fill: [colors.successLight, colors.success],
            radius: 30
        });
    }

    if ($(".dashboard-pie-3").length) {
        $(".dashboard-pie-3").peity("pie", {
            fill: [colors.warningLight, colors.warning],
            radius: 30
        });
    }

    if ($(".dashboard-pie-4").length) {
        $(".dashboard-pie-4").peity("pie", {
            fill: [colors.infoLight, colors.info],
            radius: 30
        });
    }

    ////////////////////////////////////////////

    function bar_chart() {
        if ($('#chart-ticket-status').length > 0) {
            var dataSource = [
                { country: "USA", hydro: 59.8, oil: 937.6, gas: 582, coal: 564.3, nuclear: 187.9 },
                { country: "China", hydro: 74.2, oil: 308.6, gas: 35.1, coal: 956.9, nuclear: 11.3 },
                { country: "Russia", hydro: 40, oil: 128.5, gas: 361.8, coal: 105, nuclear: 32.4 },
                { country: "Japan", hydro: 22.6, oil: 241.5, gas: 64.9, coal: 120.8, nuclear: 64.8 },
                { country: "India", hydro: 19, oil: 119.3, gas: 28.9, coal: 204.8, nuclear: 3.8 },
                { country: "Germany", hydro: 6.1, oil: 123.6, gas: 77.3, coal: 85.7, nuclear: 37.8 }
            ];

            // Return with commas in between
            var numberWithCommas = function(x) {
                return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
            };

            var dataPack1 = [40, 47, 44, 38, 27, 40, 47, 44, 38, 27, 40, 27];
            var dataPack2 = [10, 12, 7, 5, 4, 10, 12, 7, 5, 4, 10, 12];
            var dataPack3 = [17, 11, 22, 18, 12, 17, 11, 22, 18, 12, 17, 11];
            var dates = ["Jan", "Jan", "Jan", "Apr", "May", "Jun", "Jul", "Aug", "Sept", "Oct", "Nov", "Dec"];

            var bar_ctx = document.getElementById('chart-ticket-status');

            bar_ctx.height = 115;

            new Chart(bar_ctx, {
                type: 'bar',
                data: {
                    labels: dates,
                    datasets: [{
                            label: 'Pending Tickets',
                            data: dataPack1,
                            backgroundColor: colors.primaryLight,
                            hoverBorderWidth: 0
                        },
                        {
                            label: 'Solved Tickets',
                            data: dataPack2,
                            backgroundColor: colors.successLight,
                            hoverBorderWidth: 0
                        },
                        {
                            label: 'Open Tickets',
                            data: dataPack3,
                            backgroundColor: colors.dangerLight,
                            hoverBorderWidth: 0
                        },
                    ]
                },
                options: {
                    legend: {
                        display: false
                    },
                    animation: {
                        duration: 10,
                    },
                    tooltips: {
                        mode: 'label',
                        callbacks: {
                            label: function(tooltipItem, data) {
                                return data.datasets[tooltipItem.datasetIndex].label + ": " + numberWithCommas(tooltipItem.yLabel);
                            }
                        }
                    },
                    scales: {
                        xAxes: [{
                            stacked: true,
                            gridLines: { display: false },
                            ticks: {
                                fontSize: 11,
                                fontColor: chartColors.fontColor
                            }
                        }],
                        yAxes: [{
                            stacked: true,
                            ticks: {
                                callback: function(value) {
                                    return numberWithCommas(value);
                                },
                                fontSize: 11,
                                fontColor: chartColors.fontColor
                            },
                        }],
                    }
                },
                plugins: [{
                    beforeInit: function(chart) {
                        chart.data.labels.forEach(function(value, index, array) {
                            var a = [];
                            a.push(value.slice(0, 5));
                            var i = 1;
                            while (value.length > (i * 5)) {
                                a.push(value.slice(i * 5, (i + 1) * 5));
                                i++;
                            }
                            array[index] = a;
                        })
                    }
                }]
            });
        }
    }

    bar_chart();

    function users_chart() {
        if ($('#transaction_chart_canvas').length > 0) {
            $("#transaction_chart_canvas").remove();
            $(".transaction_chart_div").html('<canvas id="transaction_chart_canvas"></canvas>');
            var element = document.getElementById("transaction_chart_canvas");
            var daterange = $('#transaction_chart_form #daterange').val();
            var userChart = new Chart(element, {});
            $.post(current_url2 + "/transaction_amount_status", { daterange: daterange }, function(data, status) {
                element.height = 110;
                userChart.destroy();
                userChart = new Chart(element, {
                    type: 'bar',
                    data: {
                        labels: data.months,
                        datasets: [{
                                label: data.approved_title,
                                type: "bar",
                                backgroundColor: colors.info,
                                data: data.approved_transaction,
                            }, {
                                label: data.pending_title,
                                type: "bar",
                                backgroundColor: colors.infoLight,
                                data: data.pending_transaction
                            },
                            {
                                label: data.declined_title,
                                type: "bar",
                                backgroundColor: colors.danger,
                                data: data.declined_transaction
                            }
                        ]
                    },
                    options: {
                        title: {
                            display: true,
                            text: data.transaction_amount_detail_label
                        },
                        legend: { display: true },
                        scales: {
                            yAxes: [{
                                fontSize: 11,
                                fontColor: chartColors.fontColor,
                            }],
                            xAxes: [{
                                gridLines: { display: true },
                                ticks: {
                                    display: true
                                }
                            }]
                        },
                        tooltips: {
                            callbacks: {
                                label: function(tooltipItem, datas) {
                                    var amount = tooltipItem.yLabel.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
                                    return "Total : " + amount + " XOF";
                                }
                            }
                        }
                    }
                });
            });
        }
    }

    users_chart();

    $("#transaction_chart_search").click(function() {
        users_chart();
    });
    $("#transaction_chart_reset").click(function() {
        // $("#transaction_chart_form")[0].reset();
        users_chart();
    });

    function transaction_status_chart() {
        if ($('#transaction_status_chart_canvas').length > 0) {
            $("#transaction_status_chart_canvas").remove();
            $(".transaction_status_chart_div").html('<canvas id="transaction_status_chart_canvas"></canvas>');
            var element = document.getElementById("transaction_status_chart_canvas");
            var daterange = $('#transaction_status_chart_form #daterange').val();
            var userChart = new Chart(element, {});
            $.post(current_url2 + "/transaction_count_status", { daterange: daterange }, function(data, status) {
                element.height = 110;
                // console.log(data);
                userChart.destroy();
                userChart = new Chart(element, {
                    type: 'line',
                    data: {
                        labels: data.months,
                        datasets: [{
                                label: data.approved_title,
                                type: "line",
                                borderColor: colors.success,
                                data: data.approved_transaction,
                            }, {
                                label: data.pending_title,
                                type: "line",
                                borderColor: colors.info,
                                data: data.pending_transaction
                            },
                            {
                                label: data.declined_title,
                                type: "line",
                                borderColor: colors.danger,
                                data: data.declined_transaction
                            }
                        ]
                    },
                    options: {
                        title: {
                            display: true,
                            text: data.title
                        },
                        legend: { display: true },
                        scales: {
                            yAxes: [{
                                fontSize: 11,
                                fontColor: chartColors.fontColor,
                            }, ],
                            xAxes: [{
                                gridLines: { display: true },
                                ticks: {
                                    display: true
                                }
                            }]
                        }
                    }
                });
            });
        }
    }

    transaction_status_chart();

    $("#transaction_status_chart_search").click(function() {
        transaction_status_chart();
    });
    $("#transaction_status_chart_reset").click(function() {
        // $("#transaction_status_chart_form")[0].reset();
        transaction_status_chart();
    });


    function transaction_by_operator_chart() {
        if ($('#transaction_by_operator_chart_canvas').length > 0) {
            $("#transaction_by_operator_chart_canvas").remove();
            $(".transaction_by_operator_chart_div").html('<canvas id="transaction_by_operator_chart_canvas"></canvas>');
            var element = document.getElementById("transaction_by_operator_chart_canvas");
            var status = $('#transaction_by_operator_chart_form #status').val();
            var userChart = new Chart(element, {});
            $.post(current_url2 + "/transaction_by_operator", { status: status }, function(data, status) {
                element.height = 110;
                // console.log(data);
                userChart.destroy();
                userChart = new Chart(element, {
                    type: 'bar',
                    data: {
                        labels: data.operator_list,
                        datasets: [{
                            label: data.chart_label,
                            backgroundColor: [
                                colors.primary,
                                colors.secondary,
                                colors.success,
                                colors.warning,
                                colors.info
                            ],
                            data: data.amount
                        }]
                    },
                    options: {
                        legend: { display: false },
                        title: {
                            display: true,
                            text: data.title
                        },
                        scales: {
                            yAxes: [{
                                scaleLabel: {
                                    display: true,
                                    labelString: amountLabel
                                }
                            }]
                        },
                        tooltips: {
                            callbacks: {
                                label: function(tooltipItem, datas) {
                                    var amount = tooltipItem.yLabel.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
                                    return "Total : " + amount + " XOF";
                                }
                            }
                        }
                    }
                });
            });
        }
    }

    transaction_by_operator_chart();

    $("#transaction_by_operator_chart_search").click(function() {
        transaction_by_operator_chart();
    });
    $("#transaction_by_operator_chart_reset").click(function() {
        $("#transaction_by_operator_chart_form")[0].reset();
        transaction_by_operator_chart();
    });


    function transaction_by_channel_chart() {
        if ($('#transaction_by_channel_chart_canvas').length > 0) {
            $("#transaction_by_channel_chart_canvas").remove();
            $(".transaction_by_channel_chart_div").html('<canvas id="transaction_by_channel_chart_canvas"></canvas>');
            var element = document.getElementById("transaction_by_channel_chart_canvas");
            var status = $('#transaction_by_channel_chart_form #status').val();
            var userChart = new Chart(element, {});
            $.post(current_url2 + "/transaction_by_channel", { status: status }, function(data, status) {
                element.height = 110;
                // console.log(data);
                userChart.destroy();
                userChart = new Chart(element, {
                    type: 'bar',
                    data: {
                        labels: data.channel_list,
                        datasets: [{
                            label: data.chart_label,
                            backgroundColor: [
                                colors.secondary,
                                colors.success,
                                colors.warning,
                                colors.info,
                                colors.primaryLight,
                                colors.secondaryLight,
                                colors.successLight,
                                colors.warningLight,
                                colors.infoLight,
                                colors.danger,
                                colors.dangerLight
                            ],
                            data: data.amount
                        }]
                    },
                    options: {
                        legend: { display: false },
                        title: {
                            display: true,
                            text: data.title
                        },
                        scales: {
                            yAxes: [{
                                scaleLabel: {
                                    display: true,
                                    labelString: amountLabel
                                }
                            }]
                        },
                        tooltips: {
                            callbacks: {
                                label: function(tooltipItem, datas) {
                                    var amount = tooltipItem.yLabel.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
                                    return "Total : " + amount + " XOF";
                                }
                            }
                        }
                    },

                });
            });
        }
    }

    transaction_by_channel_chart();

    $("#transaction_by_channel_chart_search").click(function() {
        transaction_by_channel_chart();
    });
    $("#transaction_by_channel_chart_reset").click(function() {
        $("#transaction_by_channel_chart_form")[0].reset();
        transaction_by_channel_chart();
    });


    function users_chartALiveBackup() {
        if ($('#users-chartA').length > 0) {
            var lastDate = 0;
            var yValue = []
            var xValue = []
            var TICKINTERVAL = 86400000
            let XAXISRANGE = 777600000

            // function getDayWiseTimeSeries(baseval, count, yrange) {
            //     var i = 0;
            //     while (i < count) {
            //         var x = baseval;
            //         var y = Math.floor(Math.random() * (yrange.max - yrange.min + 1)) + yrange.min;

            //         data.push({
            //             x,
            //             y
            //         });
            //         lastDate = baseval
            //         baseval += TICKINTERVAL;
            //         i++;
            //     }
            // }

            // getDayWiseTimeSeries(new Date('11 Feb 2020 GMT').getTime(), 10, {
            //     min: 10,
            //     max: 90
            // });

            function getNewSeries(baseval, yrange) {
                var newDate = baseval + TICKINTERVAL;
                lastDate = newDate;

                for (var i = 0; i < data.length - 10; i++) {
                    // IMPORTANT
                    // we reset the x and y of the data which is out of drawing area
                    // to prevent memory leaks
                    data[i].x = newDate - XAXISRANGE - TICKINTERVAL
                    data[i].y = 0
                }

                data.push({
                    x: newDate,
                    y: Math.floor(Math.random() * (yrange.max - yrange.min + 1)) + yrange.min
                })

            }

            function resetData() {
                // Alternatively, you can also reset the data at certain intervals to prevent creating a huge series
                data = data.slice(data.length - 10, data.length);
            }
            var chart;
            var options = {
                // chart: {
                //     height: 350,
                //     type: 'area',
                // },
                chart: {
                    locales: [apexLocales],
                    defaultLocale: apexDefaultLocale,
                    height: 270,
                    type: 'line',
                    animations: {
                        enabled: true,
                        easing: 'linear',
                        dynamicAnimation: {
                            speed: 100
                        }
                    },
                    toolbar: {
                        show: false
                    },
                    zoom: {
                        enabled: false
                    }
                },
                dataLabels: {
                    enabled: false
                },
                stroke: {
                    curve: 'smooth'
                },
                series: [{
                    name: 'series1',
                    // data: [31, 40, 28, 51, 42, 109, 100]
                    data: []
                }],

                xaxis: {
                    type: 'datetime',
                    // categories: ["2018-09-19 15:04", "2018-09-19 01:30", "2018-09-19 02:30", "2018-09-19 03:30", "2018-09-19 04:30", "2018-09-19 05:30", "2018-09-19 06:30"],
                    categories: [],
                },
                tooltip: {
                    x: {
                        format: 'HH:mm'
                    },
                }
            }

            chart = new ApexCharts(
                document.querySelector("#users-chartA"),
                options
            );

            chart.render();

            var date = [];
            var amount = [];
            $.post(current_url2 + "/transaction_live", {}, function(ddata, status) {
                amount = ddata.approved_transaction;
                date = ddata.months;
            });

            var dateinner = [];
            var amountinner = [];
            var i = 0;
            window.setInterval(function() {
                // console.log(i);

                if (date.length > 0 && amount.length > 0 && i < date.length) {
                    dateinner.push(date[i]);
                    amountinner.push(amount[i]);
                    if (i > 2) {
                        chart.updateOptions({
                            xaxis: {
                                categories: dateinner
                            },
                        });
                        chart.updateSeries([{
                            name: 'Amount',
                            data: amountinner
                        }]);
                    }
                    i++;
                } else {
                    $.post(current_url2 + "/transaction_live", {}, function(ddata, status) {
                        console.log(ddata);
                        amount = ddata.approved_transaction;
                        date = ddata.months;
                    });
                    dateinner = [];
                    amountinner = [];
                    i = 0;
                }

                if (i > 20) {
                    dateinner.shift();
                    amountinner.shift();
                }

            }, 1000);
        }
    }

    function users_chartA() {
        if ($('#users-chartA').length > 0) {

            var chart;
            var xAxis = [];
            var yAxis = [];
            $.post(current_url2 + "/transaction_live", {}, function(ddata, status) {
                if (ddata.transaction) {
                    $.each(ddata.transaction, function(hour, amount) {
                        xAxis.push(hour);
                        yAxis.push(amount);

                    });
                }
                var options = {
                    chart: {
                        locales: [apexLocales],
                        defaultLocale: apexDefaultLocale,
                        height: 320,
                        type: 'line',
                        width: "100%",
                        animations: {
                            enabled: true,
                            easing: 'linear',
                            dynamicAnimation: {
                                speed: 100
                            }
                        },
                        toolbar: {
                            show: true,
                            tools: {
                                download: true,
                                selection: true,
                                zoom: true,
                                zoomin: true,
                                zoomout: true,
                                pan: false,
                                reset: true
                            },
                        },
                        zoom: {
                            enabled: true
                        },
                        showToolTip: 1
                    },
                    dataLabels: {
                        enabled: true
                    },
                    stroke: {
                        curve: 'smooth'
                    },
                    series: [{
                        name: amountLabel,
                        // data: [31, 40, 28, 51, 42, 109, 100]
                        data: yAxis
                    }],
                    xaxis: {
                        type: 'datetime',
                        // categories: ["2018-09-19 15:04", "2018-09-19 01:30", "2018-09-19 02:30", "2018-09-19 03:30", "2018-09-19 04:30", "2018-09-19 05:30", "2018-09-19 06:30"],
                        categories: xAxis,
                        labels: {
                            show: true,
                            rotate: -45,
                            rotateAlways: false,
                            hideOverlappingLabels: false,
                            showDuplicates: false,
                            trim: false,
                            minHeight: undefined,
                            maxHeight: 120,
                            style: {
                                colors: [],
                                fontSize: '12px',
                                fontWeight: 400,
                            },
                            offsetX: 0,
                            offsetY: 0,
                            format: undefined,
                            formatter: undefined,
                            datetimeUTC: true,
                            datetimeFormatter: {
                                year: 'yyyy',
                                month: "MMM 'yy",
                                day: 'dd MMM',
                                hour: 'HH:mm',
                            },
                        },
                        axisBorder: {
                            show: true,
                            color: '#78909C',
                            height: 1,
                            width: '100%',
                            offsetX: 0,
                            offsetY: 0
                        },
                        axisTicks: {
                            show: true,
                            borderType: 'solid',
                            color: '#78909C',
                            height: 6,
                            offsetX: 0,
                            offsetY: 0
                        },
                        tickAmount: undefined,
                        tickPlacement: 'on',
                        min: undefined,
                        max: undefined,
                        range: undefined,
                        floating: false,
                        position: 'bottom',
                        title: {
                            text: "Time(HH:MM)",
                            offsetX: 10,
                            offsetY: 10,
                            style: {
                                color: undefined,
                                fontSize: '12px',
                                fontWeight: 600,
                            },
                        },
                        crosshairs: {
                            show: true,
                            width: 1,
                            position: 'back',
                            opacity: 0.9,
                            stroke: {
                                color: '#b6b6b6',
                                width: 0,
                                dashArray: 0,
                            },
                            fill: {
                                type: 'solid',
                                color: '#B1B9C4',
                                gradient: {
                                    colorFrom: '#D8E3F0',
                                    colorTo: '#BED1E6',
                                    stops: [0, 100],
                                    opacityFrom: 0.4,
                                    opacityTo: 0.5,
                                },
                            },
                            dropShadow: {
                                enabled: true,
                                top: 0,
                                left: 0,
                                blur: 1,
                                opacity: 0.4,
                            },
                        }
                    },
                    yaxis: {
                        title: {
                            text: amountLabel,
                            offsetX: 10,
                            offsetY: 10,
                            style: {
                                color: undefined,
                                fontSize: '12px',
                                fontWeight: 600,
                            },
                        },
                    },
                    tooltip: {
                        y: {
                            show: true,
                            format: 'HH:mm',
                            formatter: function(value, { series, seriesIndex, dataPointIndex, w }) {
                                var amount = value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
                                return amount + " XOF";
                            }
                        },
                        x: {
                            show: true,
                            format: 'HH:mm',
                        }
                    }
                }

                chart = new ApexCharts(
                    document.querySelector("#users-chartA"),
                    options
                );

                chart.render();
            });

            window.setInterval(function() {
                var xAxis = [];
                var yAxis = [];
                $.post(current_url2 + "/transaction_live", {}, function(ddata, status) {
                    if (ddata.transaction) {
                        $.each(ddata.transaction, function(hour, amount) {
                            xAxis.push(hour);
                            yAxis.push(amount);

                        });
                    }
                    chart.updateOptions({
                        xaxis: {
                            categories: xAxis
                        },
                    });
                    chart.updateSeries([{
                        name: 'Amount',
                        data: yAxis
                    }]);
                });
            }, 60000);
        }
    }

    users_chartA();


    function device_session_chartA() {
        if ($('#device_session_chartA').length) {
            var options = {
                chart: {
                    type: 'area',
                    stacked: true,
                    events: {
                        selection: function(chart, e) {
                            // console.log(new Date(e.xaxis.min))
                        }
                    },
                    toolbar: {
                        show: false,
                    }

                },
                colors: ['#008FFB', '#00E396', '#CED4DC'],
                dataLabels: {
                    enabled: false
                },
                stroke: {
                    curve: 'smooth',
                    width: 1
                },
                series: [{
                        name: 'South',
                        data: generateDayWiseTimeSeries(new Date('11 Feb 2017 GMT').getTime(), 20, {
                            min: 10,
                            max: 60
                        })
                    },
                    {
                        name: 'North',
                        data: generateDayWiseTimeSeries(new Date('11 Feb 2017 GMT').getTime(), 20, {
                            min: 10,
                            max: 20
                        })
                    },

                    {
                        name: 'Central',
                        data: generateDayWiseTimeSeries(new Date('11 Feb 2017 GMT').getTime(), 20, {
                            min: 10,
                            max: 15
                        })
                    }
                ],
                fill: {
                    type: 'gradient',
                    gradient: {
                        opacityFrom: 0.6,
                        opacityTo: 0,
                    }
                },
                legend: {
                    show: false,
                    position: 'top',
                    horizontalAlign: 'left'
                },
                xaxis: {
                    type: 'datetime'
                },
            };

            var chart = new ApexCharts(
                document.querySelector("#device_session_chartA"),
                options
            );

            chart.render();

            /*
              // this function will generate output in this format
              // data = [
                  [timestamp, 23],
                  [timestamp, 33],
                  [timestamp, 12]
                  ...
              ]
              */
            function generateDayWiseTimeSeries(baseval, count, yrange) {
                var i = 0;
                var series = [];
                while (i < count) {
                    var x = baseval;
                    var y = Math.floor(Math.random() * (yrange.max - yrange.min + 1)) + yrange.min;

                    series.push([x, y]);
                    baseval += 86400000;
                    i++;
                }
                return series;
            }
        }
    }

    device_session_chartA();

    function device_session_chart() {
        if ($('#device_session_chart').length > 0) {
            var element = document.getElementById("device_session_chart");
            element.height = 155;
            new Chart(element, {
                type: 'line',
                data: {
                    labels: [1500, 1600, 1700, 1750, 1800, 1850, 1900, 1950, 1999, 2050],
                    datasets: [{
                        data: [2186, 2000, 1900, 2300, 2150, 2100, 2350, 2500, 2400, 2390],
                        label: "Mobile",
                        backgroundColor: colors.primary,
                        borderColor: colors.primary,
                        fill: false
                    }, {
                        data: [1282, 1000, 1290, 1302, 1400, 1250, 1350, 1402, 1700, 1967],
                        label: "Desktop",
                        backgroundColor: colors.success,
                        borderColor: colors.success,
                        fill: false
                    }, {
                        data: [500, 700, 900, 800, 600, 850, 900, 550, 750, 690],
                        label: "Other",
                        backgroundColor: colors.warning,
                        borderColor: colors.warning,
                        fill: false
                    }]
                },
                options: {
                    title: {
                        display: false
                    },
                    legend: { display: false },
                    scales: {
                        yAxes: [{
                            gridLines: { display: false },
                            ticks: {
                                display: false
                            }
                        }],
                        xAxes: [{
                            gridLines: { display: false },
                            ticks: {
                                fontSize: 11,
                                fontColor: chartColors.fontColor
                            }
                        }],
                    }
                }
            });
        }
    }

    device_session_chart();

});