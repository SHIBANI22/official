/*
 * @created_on    : Tuesday 23 June 2019 01:28:06+0530
 * @package       : datatable pipeline
 * @author3       : Supernet Dev Team (RAJENDIRAN SELVAM)
 * @copyright     : Copyright (c) Supernet Technologies India Pvt Ltd (http://supernet-india.com/)
 * @license       : Supernet Technologies India Pvt Ltd
 * @since         : APRIL 2017 to PRESENT
 * @project       : Supernet
 * @file          : datatable-pipeline.js
 */

//
// Pipelining function for DataTables. To be used to the `ajax` option of DataTables
//
$.fn.dataTable.pipeline = function(opts) {

    var hidden_selected_language = $("#language").val();
    if (!hidden_selected_language || hidden_selected_language === '') {
        hidden_selected_language = "French";
    }
    var headerlength = $(".listTable thead th").length;

    if (hidden_selected_language === "English") {
        $(".listTable tbody").html("<tr><td class='text-center' colspan='" + headerlength + "'>Processing....</td></tr>");
    } else {
        $(".listTable tbody").html("<tr><td class='text-center' colspan='" + headerlength + "'>Traitement en cours...</td></tr>");
    }
    // Configuration options
    var conf = $.extend({
        pages: 5, // number of pages to cache
        url: '', // script url
        data: null, // function or object with parameters to send to the server
        // matching how `ajax.data` works in DataTables
        method: 'POST' // Ajax HTTP method
    }, opts);

    // Private variables for storing the cache
    var cacheLower = -1;
    var cacheUpper = null;
    var cacheLastRequest = null;
    var cacheLastJson = null;

    return function(request, drawCallback, settings) {
        var ajax = false;
        var requestStart = request.start;
        var drawStart = request.start;
        var requestLength = request.length;
        var requestEnd = requestStart + requestLength;
        
        if (settings.clearCache) {
            // API requested that the cache be cleared
            ajax = true;
            settings.clearCache = false;
        } else if (cacheLower < 0 || requestStart < cacheLower || requestEnd > cacheUpper || requestLength == -1) {
            // outside cached data - need to make a request
            ajax = true;
        } else if (JSON.stringify(request.order) !== JSON.stringify(cacheLastRequest.order) ||
            JSON.stringify(request.columns) !== JSON.stringify(cacheLastRequest.columns) ||
            JSON.stringify(request.search) !== JSON.stringify(cacheLastRequest.search)
        ) {
            // properties changed (ordering, columns, searching)
            ajax = true;
        }

        // Store the request for checking next time around
        cacheLastRequest = $.extend(true, {}, request);

        if (ajax) {
            // Need data from the server
            if (requestStart < cacheLower) {
                requestStart = requestStart - (requestLength * (conf.pages - 1));

                if (requestStart < 0) {
                    requestStart = 0;
                }
            }

            cacheLower = requestStart;
            cacheUpper = requestStart + (requestLength * conf.pages);

            request.start = requestStart;
            request.length = requestLength * conf.pages;

            // Provide the same `data` options as DataTables.
            if (typeof conf.data === 'function') {
                // As a function it is executed with the data object as an arg
                // for manipulation. If an object is returned, it is used as the
                // data object to submit
                var d = conf.data(request);
                if (d) {
                    $.extend(request, d);
                }
            } else if ($.isPlainObject(conf.data)) {
                // As an object, the data given extends the default
                $.extend(request, conf.data);
            }

            settings.jqXHR = $.ajax({
                "type": conf.method,
                "url": conf.url,
                "data": request,
                "dataType": "json",
                "cache": false,
                "success": function(json) {
                    if (json.ErrorMessage && json.ErrorMessage != "") {
                        $.alert({
                            title: json.alerttitle,
                            type: 'blue',
                            content: json.ErrorMessage,
                            closeIcon: true,
                            closeIconClass: 'fa fa-close',
                            buttons: {
                                Ok: {
                                    text: json.okbtntext,
                                    btnClass: 'btn-red btn-round',
                                    action: function() {

                                    }
                                },
                            }
                        });
                    }
                    cacheLastJson = $.extend(true, {}, json);

                    if (cacheLower != drawStart) {
                        json.data.splice(0, drawStart - cacheLower);
                    }
                    if (requestLength > -1) {
                        json.data.splice(requestLength, json.data.length);
                    } else if (requestLength == -1) {
                        json.data.splice(json.data.length, json.data.length);
                    }

                    drawCallback(json);
                }
            });
        } else {
            json = $.extend(true, {}, cacheLastJson);
            json.draw = request.draw; // Update the echo for each response
            json.data.splice(0, requestStart - cacheLower);
            json.data.splice(requestLength, json.data.length);
            drawCallback(json);
        }
    }
};

// Register an API method that will empty the pipelined data, forcing an Ajax
// fetch on the next draw (i.e. `table.clearPipeline().draw()`)
$.fn.dataTable.Api.register('clearPipeline()', function() {
    return this.iterator('table', function(settings) {
        settings.clearCache = true;
    });
});