/* For Export Buttons available inside jquery-datatable "server side processing" - Start
 - due to "server side processing" jquery datatble doesn't support all data to be exported
 - below function makes the datatable to export all records when "server side processing" is on */

    var hidden_selected_language = $("#language").val();
    var searchErrorHeader, searchErrorContent, okButtonText, allTextLabel, exportDataLabel, stillExportDataLabel, dateRangeErrContent = "";
    if (hidden_selected_language != "" && hidden_selected_language === "English") {
        searchErrorHeader = "Search Error";
        searchErrorContent = "Please fill atleast one column!";
        dateRangeErrContent = "Please select the date not more than three months";
        okButtonText = "OK";
        allTextLabel = "All";
        exportDataLabel = "Exporting Data...";
        stillExportDataLabel = "Still Exporting Data...";
    } else {
        searchErrorHeader = "Erreur de recherche";
        searchErrorContent = "Veuillez remplir au moins une colonne!";
        dateRangeErrContent = "Veuillez entrer une date inférieure à 3 mois";
        okButtonText = "Valider";
        allTextLabel = "Tous";
        exportDataLabel = "Exportation de données ...";
        stillExportDataLabel = "Toujours en train d'exporter des données ...";
    }

$.fn.dataTable.customExportAction = function (e, dt, button, config) {
    /*$(`<div class="theme-loader">
     <div class="loader-track">
     <div class="loader-bar"></div>
     </div>
     </div>`).appendTo('body');*/
    /*$.LoadingOverlay("show", {
     image: "",
     progress    : true,
     text: "Exporting Data...",
     textAutoResize: true, // Boolean
     textResizeFactor: 0.2,
     //fontawesome: "fa fa-cog fa-spin"
     });*/
    $('body').loadingModal("show");
    //$('body').loadingModal({text: 'Exporting Data...', animation : "cirlce", backgroundColor : "black"});
    $('body').loadingModal({text: exportDataLabel, backgroundColor: 'gray'}).loadingModal('animation', 'circle');
    setTimeout(function () {
        //$.LoadingOverlay("text", "Still Exporting Data...");
        $('body').loadingModal({text: stillExportDataLabel, backgroundColor: 'gray'}).loadingModal('animation', 'circle');
    }, 1500);

    //console.log("customExportAction called");
    var self = this;
    var oldStart = dt.settings()[0]._iDisplayStart;
    dt.one('preXhr', function (e, s, data) {
        // Just this once, load all data from the server...
        data.start = 0;
        data.length = 2147483647;
        dt.one('preDraw', function (e, settings) {
            // Call the original action function
            if (button[0].className.indexOf('buttons-copy') >= 0) {
                $.fn.dataTable.ext.buttons.copyHtml5.action.call(self, e, dt, button, config);
            } else if (button[0].className.indexOf('buttons-excel') >= 0) {
                $.fn.dataTable.ext.buttons.excelHtml5.available(dt, config) ?
                        $.fn.dataTable.ext.buttons.excelHtml5.action.call(self, e, dt, button, config) :
                        $.fn.dataTable.ext.buttons.excelFlash.action.call(self, e, dt, button, config);
            } else if (button[0].className.indexOf('buttons-csv') >= 0) {
                $.fn.dataTable.ext.buttons.csvHtml5.available(dt, config) ?
                        $.fn.dataTable.ext.buttons.csvHtml5.action.call(self, e, dt, button, config) :
                        $.fn.dataTable.ext.buttons.csvFlash.action.call(self, e, dt, button, config);
            } else if (button[0].className.indexOf('buttons-pdf') >= 0) {
                $.fn.dataTable.ext.buttons.pdfHtml5.available(dt, config) ?
                        $.fn.dataTable.ext.buttons.pdfHtml5.action.call(self, e, dt, button, config) :
                        $.fn.dataTable.ext.buttons.pdfFlash.action.call(self, e, dt, button, config);
            } else if (button[0].className.indexOf('buttons-print') >= 0) {
                $.fn.dataTable.ext.buttons.print.action(e, dt, button, config);
            }
            dt.one('preXhr', function (e, s, data) {
                // DataTables thinks the first item displayed is index 0, but we're not drawing that.
                // Set the property to what it was before exporting.
                settings._iDisplayStart = oldStart;
                data.start = oldStart;
            });
            // Reload the grid with the original page. Otherwise, API functions like table.cell(this) don't work properly.
            setTimeout(dt.ajax.reload, 0);
            // Prevent rendering of the full data to the DOM
            return false;
        });
    });
    // Requery the server with the new one-time export settings

    dt.ajax.reload();
    setTimeout(function () {
        $('body').loadingModal('hide');
    }, 1500);
    //$(".theme-loader").remove();
};
//For Export Buttons available inside jquery-datatable "server side processing" - End



$(".dataTables_processing").text(''); //Removing Processing... text for loading

//localStorage.clear();
/*Thsi is for First time and only once clear localStorage */
$(document).one('ready', function() {
    //localStorage.clear();
});
var required_columns_orientation = 6;

var table;

$(window).bind("load", function() {
    $(".dataTables_processing").text(''); //Removing Processing... text for loading
    //console.log($(".dataTables_processing").text());
});

$(window).bind("unload", function() {
    $(".dataTables_processing").text(''); //Removing Processing... text for loading
    //console.log($(".dataTables_processing").text());
});

$(document).ready(function() {
    
    function parseDate(date){
    	 var md = date.split("/");
		return new Date(md[0], md[1], md[2]);
    }
    
    function datediff(from, to){
    	return Math.round((to-from)/(1000*60*60*24));
    }

    if ($('.listTable').length) {
        $(".dataTables_processing").text(''); //Removing Processing text for loading
        var headerlength = $(".listTable thead th").length;
        if (hidden_selected_language === "English") {
            $(".listTable tbody").html("<tr><td class='text-center' colspan='" + headerlength + "'>Processing....</td></tr>");
        } else {
            $(".listTable tbody").html("<tr><td class='text-center' colspan='" + headerlength + "'>Traitement en cours...</td></tr>");
        }
        var order_by = $(".listTable").data("orderby");
        var order_by_column = $(".listTable").data("orderby-column");
        var documentTitle, fileName = $(".listTable").data("document-title");
        //console.log(documentTitle);
        if (!order_by || order_by === '') {
            order_by = "desc";
        }
        if (!order_by_column || order_by_column === '') {
            order_by_column = "0";
        }
        if (!hidden_selected_language || hidden_selected_language === '') {
            hidden_selected_language = "French";
        }

        var get_url = $(".listTable").data('url');

        var columns_count = $('.listTable thead > tr > th').length;
        //console.log(columns_count);
        if (get_url) {
            table = $('.listTable').on('length.dt', function(e, settings, len) {
                $(".form-loader").show();
                //console.log('New page length: ' + len);
            }).on('page.dt', function() {
                $(".form-loader").show();
            }).on('draw.dt', function() {
                if (hidden_selected_language === "English") {
                    $(".dataTables_processing").text('Processing...');
                } else {
                    $(".dataTables_processing").text('Traitement en cours...');
                }
            }).DataTable({
                "lengthMenu": [
                    [10, 25, 50, 100, 250, -1],
                    [10, 25, 50, 100, 250, allTextLabel]
                ],
                rowReorder: {
                    selector: 'td:nth-child(2)'
                },
                responsive: false,
                "drawCallback": function() {
                    $(".form-loader").hide();
                    var info = table.page.info();
                    var tot = info.recordsTotal;
                    if (tot == 0) {
                        table.buttons().nodes().addClass('hidden');
                    } else {
                        table.buttons().nodes().removeClass('hidden');
                    }
                },
                createdRow: function(row, data, dataIndex) {

                    var tdHtml = $(row).html();
                    tdHtml = $.parseHTML(tdHtml);
                    // console.log(tdHtml);
                    $.each(tdHtml, function(indx, val) {
                        //console.log($("tr").find('th:eq(0)').attr('class'));
                        //console.log($("tr").find('th:eq(0)').attr('data-title'));
                        $(row).find('td:eq(' + indx + ')').attr('data-title', $("tr").find('th:eq(' + indx + ')').attr('data-title'));
                        if ($("tr").find('th:eq(' + indx + ')').hasClass("hidden")) {
                            $(row).find('td:eq(' + indx + ')').addClass("hidden");
                        } else {
                            $(row).find('td:eq(' + indx + ')').removeClass("hidden");
                        }
                    });
                },
                //dom: '',
                //"dom": '<"top"lBfrtip>', //To change Place of Selection records length and exports buttons
                "dom": '<"top col-sm-12 rm-pad"<"col-sm-6 pull-left text-left rm-pad"B><"col-sm-6 pull-right text-right rm-pad"f>>rt<"bottom col-sm-12 rm-pad"<"col-sm-4 pull-left text-left rm-pad"li><"col-sm-8 pull-right text-right rm-pad"p>>',
                buttons: [{
                        extend: 'copy',
                        text: '<i class="fa fa-copy export-copy"></i>'
                    },
                    {
                        extend: 'excelHtml5',
                        filename: fileName,
                        text: '<i class="fa fa-file-excel-o export-excel"></i>',
                        title: fileName,
                        "action": $.fn.dataTable.customExportAction
                    }, {
                        extend: 'pdfHtml5',
                        text: '<i class="fa fa-file-pdf-o export-pdf"></i>',
                        "action": $.fn.dataTable.customExportAction,
                        footer: true,
                        orientation: 'landscape',
                        pageSize: 'A4',
                        filename: fileName,
                        message: fileName,
                        header: true,
                        customize: function(doc) {

                            var body_length = doc.content[2].table.body[0].length;
                            var percentage = 100 / body_length;
                            var table_widths = [];

                            for (var i = 0; i < body_length; i++) {
                                table_widths.push(percentage + "%");
                            }
                            doc.content[2].table.widths = table_widths;
                            //Remove the title created by datatTables
                            doc.content.splice(0, 1);
                            //Create a date string that we use in the footer. Format is dd-mm-yyyy
                            var now = new Date();
                            var jsDate = now.getFullYear() + '/' + (now.getMonth() + 1) + '/' + now.getDate();
                            // Logo converted to base64
                            // var logo = getBase64FromImageUrl('https://datatables.net/media/images/logo.png');
                            // The above call should work, but not when called from codepen.io
                            // So we use a online converter and paste the string in.
                            // Done on http://codebeautify.org/image-to-base64-converter
                            // It's a LONG string scroll down to see the rest of the code !!!
                            var logo = 'data:image/png;base64,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';
                            var address = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAdUAAABYCAYAAACj4G1FAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH5AQNDTYFt0V7gAAAAB1pVFh0Q29tbWVudAAAAAAAQ3JlYXRlZCB3aXRoIEdJTVBkLmUHAAAgAElEQVR42uy9f3BT95kv/PECPiKQSISMj8JFSAkkUphgqYOncmB2pIZSqWSylrd3kLiT1Erm7UrtbSPvpIOcSZql+fHK4S0jk9xdm2ZvrHRzV4J5d6XcTCzx0mDRjrHJ2GvVcLEIEGljJpY9IRJgIsmEed4/ztFPS7INTtO05zOjAVk6R9/v832e7/PzfJ86IiIIECBAgAABAm4bfyWQQIAAAQIECBCUqgABAgQIEPAnheV/ioMaeP03aOlLAwCYxr/G2f2bIRbWSqB3LZw8iv/yUhxZAFihwOH3DNjxNQ1l+PXfwCjwrwABf5r701e8VyxMqc5ew/CJs/CGJzDy8TV8kpzFVSzDXavqsUG2FlsfUcL62CZsWS0wQDVcOfa/sfnAJLeQAJj1GoT+WYstFb6bCb+PjZ2XkBU2ZQE18Smeb3sPh6aqf4NZUQ+ZbC2ampSw7VYKMiqgsCdNxOD1f4Tj0c9wZiKNKzduIotluGvNSmyQ3YPtjzwI68778IDAM0urVK+Mj+KnnSMITt0s++Qmrl5P40z0Es5EL+Htd0ZgeeZ7+NXONRDd5qA26ZvwK9mX3Bv2XjB/hoTPXjqNl49txpGdd/5xBOjkUXzrpTiurnkQIe930PQXRu+/VGRvzOLCx5O48PEkfO+fxS9e2QXHQ8IK/4VzBU4fPoof9k5iYs5nN3E1OYMzyRmcGYvDE1Dgjf0G/C379Y32m7Y/1VSqmfFT2L03gpEbRX9cUY9NsjvB1gNXplI4k+SV7Y0UfAfeQ3ZFK97U356iYBs3w9r4587YN3H8rVMYeOS72L76qxeigWMTuPoXTe8/fzDsWmxfX5/fdLI3ZjE9USSj16fx8ku/x8Nvfhc7BO/jLxaZ4ZMlCpVZ04Bdj7B4mF0GJL/AmbEJ9H2cRhZAdiqOn7nG0NTViA1f03i/aftTdaU6O4X9ncUKdSUe3f3X+NUT92FDfe5vX2Jq7DR+fuBDBKcAIA3/GyfR2mTALkFo50fyIl46osHRp+/5an9nZgL+4ZsCvf/M0bB1G448s26OQTX81nvYfeQyZ1QlL+LQiW3Y8dgdAsH+IvElBt6/WFCo67fg3X/ahqb6Mp7pCaAlkOIUa/Q0PBcb8eJGgXq3pVSvnByGJ5+rWYatT+3Cb8z3lIV2l4Nt/Bb+ZT+DJ3/yewSvA3et/gLnP/4SaCzc+sr4ORwMnMXxsRQ+Sc4iu2IZGti1aGrajGeemJvnqZ2Y/hKfDJ/FofcvYiCawvnkLLJYhgb2TjysUsC6W4NdGysECGY+w78ficA3PJXPH2DVSjwgW4vtOxvxzE4Z2PoKVt1UDJ53zsI/No3zn88ii3rIZA3YsfNbeKZ1HW4lKiJrbAAzNo0LAEbePwmv6W+w5+5F3ODzT+F5ZxT+scs4PZXG1RvLcNeaO9HUqID1ia3YJcvRfgI//UEffNeLFflHMBo/ArAMrS9a8ea25QsoBPgSU2Mf4ZD/HI5HuTW8umIZGu6WYGvjRuzZvRm7ZKU0L84hM6pH8B9djRCNjeL5t87i+MUZXEE9ZBtlsD69DfbGChv8bBIfvB+B59inGOHXi1m1Eg9slKF1dxOsTXcuMs3wJabGzuJ1/0f4IJrCRJK734aN62AxaWHdtsj7zSbR986HOHTiU5yemkV21Wo0NW3GXtuWBYSnvsDpYxG8fmwCwxevYfo6wKy5E1saZbA9sXUOLW/Tf0XTE1q0HuvD20nuL8NjU8Bj9wH4FM/veQ+HkgBQj7b9T+LFFafx0wMRHL80i4ef2o2j5jWLpN8X8O79V/xs7Gaen/5j/+Y5cnLl2Pv41oFLnKJfcS/eeLtYBhbHb6X3asCv3m6FtVyepkZhaPsQI+DmannlSfyPpuWLlCdubB+4PDCf4Ob38O4fIPz0nTj9/km87I9j+NIssqtW4uHGTdhr/zZ2sMsrRo7OhyM42BfH8MVrmLh+E8yq1XhAtQ6trdV4e6n4N42JZMHAFj8kK1OoPM88/dd4FRPIyCTYwErwMFtpb5yA7/BYGd1WYotKhtbHvgVLVRld3Pzn3Z8WtX4Lwe3J5/Kq4cITk4Vw4RoF9rbeU33R2M149ZWVcNy9Dk1s6Q+ef/997H7jUmns/sZNTFyaxsSlafQNTuA37u9ix4KUShYDPe/hycDlslDmTUxPpXB8KoLjJy/A8uzf4H8Uh6CnzuLJ9t8jmCy7XVFO2H9Mg9/s15Yw2NTJfphdH+HMjZLdFBc+voQLhy7BP/htHH7lW9hSv8hgLLsZv3gkhb8bnAWuT2J/7wRan5UtSCgy46Mwv/AhBq6Xzv9qMoXjJyI4fjKOZ1414cXGpdqYv8AHr/fhR31lNL9xE9NTlxE8dhnB8Gm0PvM43ijKpzMrlhUZNGlcGO7Hs7/8CBduFNExehEv7J3CxCv/Fa82MSWe9fPtR3HoUql3nb2exsjYRxgZiyP41OM4bL5ngRtJZb7JXk/jzNhFvDAWh9+0C0fs6xZWEDb7GXr2vocXorNFvDSDgRMfYnd0GjbVMtSKAPW80IcXxmZLR5hMYeBECgMnY2h7rhUHti2hJ1m/GhvWAOD5Pzt7ExkAItSDWZX7+01cnZrAS299OFdOFkW/O/DozgbcNcbtH9nxGAZmNuNvV5fe74NwYX+5q1GJXXffOr+JdZuxq+cSZzzeuIy+wS9gLfPEp07GcaZoP9vDG/2Ll6flENUvA8Dx5pWZaxh+K1yIBPD7ysjgaZij1/CbfzQUzY0zxryuPvx8cCZftMjRcgYjIx9hZOQifLsMOPyMrMgQWUr+rcddqwrjnx4+i77PZaVjBID6dbDa11U1UD8Jh7H7wMUiec7RbQYDg+MYGPwInp16HH52U6lBdUvz/yPuh0sgn1WeU01h+OPChnaXSoHt8yiODQ/dN0ehYmoMz/fkFOoybN2lR+jN3Rjc/9do5SmWnbqIn/dOILOQCV+M4PkcY61ai7affQ/vuk3od+/Cr3ffCxkA3JiB743foW+mwJB9PafyG8WmrRr8+pXHEXKb8O6Lj6Dtfm4DnI5G8PN3Piv81sQY/i6vUOvx6O7vof/tJ/F//vF7+EXjSu6asQ/xo7c+XdjYi630mXrsenortq/gfyp8CocuLuTKz3DoQIGB7lJtwWH3bgz+4y68snUlv/mkcOj1EZzmAoLY+8rjeHd3Q8F7WrUev9r/OELux/GLxvktuPOH+/DD/AZXj607v41fv7gLhzsegU3FM8WNNPwH+vDy+Jf560T1y4sEZgIvHfgIV2UK2ExbYNsqwV0FDQrPO2fxSXF46q3+gkJdsx6vvPIDDL5pwm9MDfx1sxh45/fwTS2M3lPHfosf5fhmxVrYOkz4j399EoP7H+H58CZGAr/F8yezC7rfJ++H8XJeoS7Dw7pv49ev7MLhZzV4FBM4dGKmqnL/4I2jeYFl1m/EG/t34//8626Ent2Ih1dw/Pv2gd/C+/kSOquzKZwvolXDmjvyxoi4aCP65P0R+JJAAyvB1vslaKi/NfqxjyjzvI0bkwgOZ+ekI4LRm3n67dipyI/jlvitXoY9upX5eQyciGOqzDA8Png5v4lv0m3m97PFyhOPIoPxSnQEP/ensGHrg3hm90Novb9oo0zG8dKRqRJlNPzW/1dQKCtWo3W3Hodf+R7e2L0em/jxn+n7LX5+7IuviH8ZbN/WUJC/ZBw/bPtf2H3gJDzhT3H+8y/nV2RjJ7G7s6BQZaoH8auO7+Hwi3q8olvL3/smzhwL44eHP7vt+S/dfji/8b0U8rm8mncyUXTRBlZySxW9mak0mMZ78X1O1GC3KzlPULYGr+w+h743ppEFMDEWwxnISipSKyqjj6dwPscajRq8+th9+XFteehesKvC8F6qR8MaCcSzBQNh4OPcGwn2PK3F3+ZzAyy2yxiI3pkA1twB2f3gLfgvMXAkgoEbOUtaizeevo+znNj74HjxJs60fQD/deDCsVEcf2LdInPINwHZZvxi52kY+2aAG5dx8K1zsLyqrG2dfZ7CNHsvvs/ywvHENux4CADW4IGnlfCNRHAGQPbSpxiYALbIGGx4aB0apooEvf4OPNy4bl5ac5vfObx0pGgz2mXAu8+s42kuw45tLBp+EsDLlzjl6O39CM9UCPUheRnnGx9BeH+u2CGLXQcOo+UYF9LJXpzA8My3sGE1AFzDJzfW4vtbbwJYBrF+G+xNXAjyAbsWrSfe48KYN6bRN5aFdSczr+D5Dl/CNP9u6xM78KqeD2ne3Yg3nr2M4b0fYQJp+I+cwy+2Nc5jIX8Gz/sFmtzVuA2Hn8vNWYYdqnq0/OTDPO+UavdzeD2c5t+shvVZPfY8xIkgu/O7eHPiMvRHUshen8ShwGfYsyS59i8w8NaH8F8veCrbmxoqfnMkeg2P2n6A35REpW6BfqsVaG1chuDITU7JDX6KjL4gq5mxiwXPYpUMrbkoxS3z23Jsf2wTNvWdxgUA2fELOP755kI4eWYCfeM3C3vATvYW5Wkuza5+nAKz+3G8+zTLjfMJJWQ/CeD1Sznb/FN8Apbj+88v4uD7KX5+y7D96V14s5WnZZMMD8/+G74TSAGYxfHDp3F+pxYPLDn/Auxj2/Diiffw85w3dmMGx4+dxvFjnOq5a40ETY0y7NI9iNZt95R5v9fge+sjXMjtw/d/G0e6voUH+Pc7timwpT4n2zcxciSCD1q/ix31tzr/pdwP57O+l0Y+KyvV2dKwG7Nm2S2Js6hRi3+pUrXFrrkDDMAR+PM0rizExlrF5K/JDo/gZ4cB27Z1eFh2B0RYju3m72L7nKuWoWFFQcF6Dv0ODbs349HGe7gcqkyJV59Tll1zGR+MpfPvHtjaAPHslwWPtH4dHt0I+McAXJ9GX/RL7GpabNx+OZfrOsEp56sjp7B/bCMO1PIe796EV1/dVPmzNXeiodgomrn97TgzfKEorLIW1tZ1pcZVPYvWx9Zi/yFuI7w6HsPIzOYKBsZK7Hlqc1H1IIPtO2WQHfuIi2LcSGMiCWA1JxB7nn0MeyqOaGVJGPPK518A82UwpybQd6lwfdNDdyIzW2SN3y9D04qPMHGDU+4DM41locq5+eyRSwXe2r5zY+kmJlPC8tCHGBirYBSOxTGcU7Yr7sF2GUrGImu6Fw1HUpgAcH5sAp/gnkVVXE6PnMSTz68sCpnOYnriMkaKcmjM/Zvh2FaZx5j7N+PV8jTPLdGPwY6d9+KuES7POT18EQOz93EbK77EwIlP80qioelBPLp6Cfht40Ow3n8aL3zMGVz+kS+wZycXorsyfCFv5DD3P4hdG5dQnlbJ4HiCLYyznsWebRIcOsIX+SSvYRrABgCZsQtFxtZatOrWlOwHW57YgXeb0nla3/VV8C8A4B5Y95vxgP8UDgYu4njZ45Jc6DSF4ydO43n2Xux99rtw5OoePo/BXxRl2N6qLFN8DLabNmHTMc7AwfUJ9EWBHY23OP9aWOL9cKnks7J01S8r2aquJGdvfWOeiOHQkdPoG7uM85/P4uqNao73QpT0ZrSycbw9xbn1/t4P4O8FmFWr8bCKxfamjWjV3Yctd5cyUOtjDTh4aBpXAUyMjeNnY+MAlkG2fi2aGtfh+/rN2NVYlByfvYbzRZ76SO+/YX1v1dgaJibSQNMtPEZ09ybs3R1BX+9lZJGGr2cE9n/SoqZBNfMZ/v3ICN4+OYXzU2lMV6Enbty+Up34eKaQw1klwQMVBrZBJoEYl7lN8sYMF2YsF+oVa7F1YxmrsRI0AIVc+2ypdzX8/jAOHZvA8EQa09dvInvLk0gV5fPTOLT3f+JQte/eSOHMFGpvSlOpotBiPTasL1fqd+CB9SuBIqMsHzaeSBfmcSOOH/7X/1k9EDVxGRP8ZrzgXP3UZb4Kv4pRym7Em/9Q3fpvUMnmfnaL9BM3KbFj1SXOQ77+KfqiX2JH43JgdhL+4dn8xvn9nYVagtvjtzVofawBL78xjSxuYuBYHFM7N4PFlxg4MZ2/b9POTaVzvE15YmT3YmtZaqxBthJAKh+Vylac3514oDyPufoebG/6ivm3iE+3t34H21u/gysTn2JgbAIDw1MYjk7jTLIw5uzUJF5+/iiYf2qFXQbg0uWiVM1KPLy+Qm5RthYPrAAfHp7FhYkvgMY7bm3+80bTlm4/XCr5rOIW3YkGFkAuhDGRwhWwiz7VJzN+Ci3Fz7muqIeMrYd4xXLg+rXC83MLxWoZXt2/AyLXSXiiBQJwSe4ZjIxcxOtvrUar3YA3HitY3Btad+EIfoufv3WpqOioUCzl74ugQbUFb764DdvvBjA7i+wilNKVZBbArT2b+0DrNlj7uFNxsh+fxUvHtuDNFdUYaALP/vc+zqjIeeHsSjSsWA7cyOKTqXTVZ1FvBVdmilRZfX3l9V9VX2SAfYkrs5W9cnH93AhC5ZTCNfz7L/34u8GCUrprzUrIVjEQ4UtMT81UF5yKNs/sInLes7iSnOcrN4oVPAOmQq0BU185spOdWYRpMJvlaFl/e2vInajUgO07N8PxWPHjcBWUKrty6ei3WgFLUz38J2YBpHE8PA00ruO8lZw3ukaBPUURntvlN1a3GY/2TCN4oygEvHoSfbkw54p7YdHdubTyVL98Dh+XFOnVkCfm6+DfChDL1mGXbB12PcaPc2IC/ndO4qUTKW7+N6Zx6MgE7M/KgJnZoqjicogrKfD6+hKaZGZv3vr8ayrUpd0Pl0o+qyhVCbZvrMfrlzhmzEbjGJhR1swbnvcfxf6Je2FpVWKHjOFi74dO5xUqc/8WvLt/G5py9zh5FPe/FF+0EhCxm/Bq1ybsnfgUA8NxfDD2Gc58fBlnpma5ze7GDPxvHEWDzIxX86FUBk2tjyH82DWcHp7AwNinGIhexsjFVH6Dno6expOdEgzu3wy2vh7iFQXr5tGfPYkjX9VzffXr8MwT6+E/cAnTmEXwnVMYeKoQ5i7G6SMnC8U5qxrwyv7HYc95gJ+fxe7/9nscX8KhiVczANJ54a4Yor8+WzTOKkK2GENsbBgvDxZCQK3PtuKNnbkoQhIHf3IEL3+8GEOME/CrPF//4k0zHLLbGOCK4ihOFtkKRsTV65WNRaaYnmsexLve71RIV9w6ZLsex+gz65aWP2+ZfsuxfacMDScuYhrAxEgcp9GAqycn8p7XJp2yJLd/2/y2eiOs204ieGKWy7mPZLFnfQzDuUKWkirjP748MWXzy34d/LsgJSuD9bnHcdfMv+LvRvgq4YnP8Alk2LCaM3au5oyaSmHVmVJDQMwbmbc0/xpY6vVbKvn8q5oCkWfkOPa/U6PKdWIMz78Th79vEOYf/Qta/ElgdgoDFwubS9POLQWFCuCTqWu3RVSxbB12tW7DgX/4Gxx9+ymcfXMHbPfnLMQZfBC+XEGB3Ykt2zbDbv8u/qXLjLP/9iRCtvX5cOvV8Qs4PsN9T3Z3cVjg2lfKxOxOLWz358KLF7H/2BcVrPQvcCaaKpwdrNoMa3FINZkqCsss0SZ9f1GV7vUUzk9UCJlcTBU2vxUSPHCbx5lNDE/n821YtQ7WncVh+ZmSKtaFEVdSFO77AuenvrzNxZIU5VBn8cml7BxP+/xEuuKlG2SrCwp55ho+mcGfPm6DfqJGJb6fS5tNfYqBi9PoG0nnDfddO9kl5rfl2P7YRl6eb2Lg5ASGB3MnB9WXVBl/HfJUsv7Xr83l5c8/hffwKRw8fAoHD5/D6dml598rw6fwo73/L/R7/hn/5Qfv499nqoeH2buXFTvBHNavLQp5pnHmUoUq3YkpnC86NOgB2R23Pv+qWPr1Wyr5rNr6TdS0Fc8UPW93JnAU5p4LOF/yQ19iamwUT74wiOP5ar570cYnoKsq4dmpkgpK4Etk5k3bZnH62Cm8dOAontz7W3inypXsJth2SgpHtPHPYWFqAp63foef/vJ/48nyx1/q70BTayO2ryqE9jjmWYsdjYVQ2IXBcxgoGd9n8Ljex49c/Xi+ZwwDt/34wz2wPf1gfjMYGal8Jme5V5Qt+t8HgQv5ijzgZnV6zs5WCdFW4oGN2JGjDS7D4y+j3+ynONRXWMeGpo3zPno1r6eK6vn1qRNjBT4DkL2+gPQBK8Oj6ws7w/H346Ue0NQ5PPvLo/jpgd/hpcMX5hfEu1lsZQt0Hjh2ruSazMXT8ESrGIKNMjTlHzWZhv9EqbE2dfJ3+NEvf4tnX/8dDoY/W/SjWl+NUr0N+tXfi79tWpnnn2DgLAZ4uWXWK7Bn49Lzm6hxM6z8eK+OncXBYX7jXSWDZRuz9PK0GEdg632FR40wDW/4s9Jo37GT+HlvBC/3RvBy3xSy9UvPv+JVWZweu8zlTa9fwsuuMZyuoDwyF8dw8GRh0uKN93BO1t33oTWvF25iwH8u/0RGTtn1BeIF2q2RYZfqNuZfFUu/fkslnzVKTdfA/tw2DOQPTZjFQOADPBL4PTbdXzj793yyOMe0Eq0/+w7+9m7OEm1igSCflx0OnIJXpUETUvD2nsSh5EpsWpXGBf6B7eCJKTRtlWDD3dUi7cuQHb+A149xHDDwSyC7W4mHZSvBzM5i+lIch/Ll+PVoalrLW7PX0Ocfx/EbAAZ/ix9CC1uTBA2rliF7PYUzx0bQdz0n6OvQdDdv8e7egu1h/tGIqXH89AVg7+5N2IBrGPYPY/8I96wVs34LLE/fvsCJmrTYu/UifjZSTVGsxAOylUCUfwxl7DT2hyWwym7i9Psn8cIxYNOaZbiQvMmt1YkLOM+ug0x2B0SrGYgBzgO8PoH9b5xF9pF6MKwMOzbWyGysLi6kAi70HYX5xla06dZAnPwMfn+kEH5Z0YBnnt502910Nty/GneBz+Vcn8D+t2J4ZeeduDoWwfM9l8CwK4Epjgbnhz/CwE4Gm9g7azxGcA8s5vU4dIB7LGF68PfYfeAL7N15D0SfT8H7zgh8/DOxsp33wTG/lsEeXaG68+rYh/jhL2/imZ33QJSMwfPOOM6vWMYZaEWGghgAWCWe2RbBwAnucYPjPX346exW7Nl4BzIT57D/rY8wcp3j3+83Ni1gLH8M3A79lmNrUZX3wLHCw9gP73xwblHUkvDbPYUK4euTCPKpAtk2ZZkCvnV5umXc/SAc+mEc5x8lG3nnKJ683gRr4zJMD5/DwbyjsQzbTVv40PgS8+9DTXhRF8cPT3BjmBgZxHf2RLD1obXYsIYBgyymJ1IYLi4qWiGBzZQrKLsTlqcfRM/ece7xpY8/xO6917DXdB8acA1njp3GwXz6ph7ff7qpQPdbmj+Wfj+sKtpLI591REQ1PYeJC3je9Xu8/fE8qn7VWtieNeDVbYVCgKlj7+M7BwrPWBVJDywvGvDosQD+brCo3L/xEfzH/kZcqHYs1cyneGlvH17/uJaHsgybdHq8+1zhJI9Pjr2PlgOXKnh/xeNvwCuvPA77QwU745PwbyucGlIUg1+zHr/ab8CeBRyDVXJs3yM7cPEfNs1d3IunoP/vkcLJL+Xznxit8gzkMjxsMuANdhjGQ9MFI2fFevzmvcewa+YCfsQ/V1uMrbb/hqOtd85zDNgX+OD19/CjvlT1/PeqtXjmuV14salow5mvZ+HUGFraBjGQyxX9oxmOjZw38tJP3ss/51eMu1Tfxru2NH7296eLaLQSbft/iAMztX7vCwy83ocny0/pKUKDSoPfvKItSVFUxeyneKm9Mh8y6x/ELxqn8UJfig9LcUc05pX+zBQOvtBXdHjE3LXcusuAI8/IFmCglLZ+W3xO9TMc/L/+jX/uE9hacixhqfdx6/Sr0J5uRQN+9c+tsLKVw3q3xG8lYcSz2N32e86Q5vcbm9uMVx8qk9NblCfxPMfmFbdunJObm/0Mnhf68PxYukpEZhke3rkDh5+9r8hQXGr+/Qxe19E5pxpVxAoJ2p7bhQPb7iyJUH5y7APsfj1edW8EVuLRpwx401xW4HoL86+6P93qflhrb1oC+ZxXG4hkm3Dgn2SwnzwHz7GLGP44lX805q5VK7FBthbbdZsr9t1jdxrw7oqTeOmdOAYupZFdsRIPPCSD9SktrA/dAWz8awxMnYL/4zSyq1ai6f552satXocXu8zY8f4oPCcmMTxxjX/cYhka1tyJB1QNaH1MA0tT6X027HwM4fvPwXPkHPqiKZz/PI2rN+bvNblB/12ENypw6MhZ7pGgqfnPLb4tbNTgFzvPwXwsXSXJ+S38Zj+Dl3pPIziewvSNZZDdfy92tWqxd+c9EM+uxKtjH2D/cArTqMemXMhm9Sb86h8uA6+fRd8lbg4ydi2aZAt5/vgO7HjmBxjUf4TX/efy5y1jRT0aWAm2Nylh3b2Z9/CXAPXr8OL+XWg4NAzP8GVMXAfErATbdRr84olN2FCfxa+e+gw/PTKJC9eXoWH9WjStATBTew7bnzEhvO0sDuXOTr1+k6tG38if4/zYIs5xrl+HF/c/jgfeOoVDg9M4n7wJZo0ETU0PYa+tEbITfrycC0fNlj0OtJqFo8uM7WWPDM1/9uvXiduh3zq0PrIahwKFBWIeUvIP7H9F/Hb3RlgaT+J4Luqz/j5YHlq+ZPJ0WwU29dwzoluPjeDQsQkM5M+WrXVm7lLz7z3Y8w9mbB/7CP5jMRwfv1z0OAp3fu8DsnvQtHUjLBX7ZC/Hhp0GhBtj8Bw+i+DINE5/zp3BLr5bgocbZbCaqpzBfkvzx9LuhzV1zO3L57yeqgABAgR8s7BQD1yAgKXHXwkkECBAwJ8TpsKncCiXPlhxL6w7BYUq4I+H5QIJBAgQ8E3HlfFz8Ee/wJWJCXiPTebrOB5+TIvWuwX6CBCUqgABAgQsGNODEfz8SKrkb3fdvwW/epr9E8tPCxCUqgABAgT8iYNZU4+GFcD0Db7LyiNb8IunNy+617EAAbcLoVBJgAABAgQIWCIIhUoCBAgQIECAoFQFCBAgQIAAQakKECBAgAABglIVIECAADSB50YAACAASURBVAECBAhKVYAAAQIECBCUqgABAgQIECAo1a8KkQ6o6uqg6ogU/XEIHao61Kk6EFmq34l2QiOqg6azSnPM1BB67EZopBKIRCJIVc2wdIYQzzfZy2CoywiFRASRxIpApuz9HzqhEYnQXO3+S4SUzwRJXR3qJBYEUn/6yxvvaoaoTgpr+Kv7jYBFxNGjwmcJjxGiOgXahzi+alfUQWT0YHGkC8EqrUNdXeElkqpgbA8gjkXeNx5Cl2cIfxJLl4nAY9FAKhJBJG2GxRMp9JTMxBHq0ENaVweRyVdlvBXoIlFAY+pAKF68/kXfEUmgaLaga6jSHTOI9FjRrJBAJJJA0WxFT6TQ5TIV6YFFJUJdnQr57WKoHYq60jHU1Ylg9FShcLwHelEd6uoUsIczt8hH5eDXX9+DxDz3uXVB8sHerIBEJIJUYyqlXyKMTqMCojoRmrvighZcYgiHP9zS5jKEDr0er/1BArXZinYNkAgF4Hvu+xiKeDHks0CKKAI9R/GfUhu8ng7oRVF0Fr9XZtDlUwEqxVepUhHyhZHRGqCNhOALp2AySYT1WzBUsPb4oZdocEtUk7fA2a6HFBlEQz3wHLTAJI0g0rHwWwx1taMjbIXR2oyve+Ui+6ywH07B6OyEKtKF1+xWaDQRdGii6DHq0RGXQMIs4EZKM9wdRkgApCI+dB18DRaoEA1Y+S8w0No6YVEBSAzB03MYf28CFFEfitk3E26H6cc+wGBHpz0BT+fb+LFViuZIJxQhO5pNAUBRRjWFCfvcirzSz0R86Hw7DqlUVMWm8WFIooNOOoSQbwgZvX5BJzRJ9Pvg86egUs3/3eaOAPxWKZqXTqOix2rFoXgzHJ12JHr24e8tHdBEe6BP+WBptmJIpBBOmvqqQN80jDpJCZDSOVr0x0FyKkFQOon7a5L6XS2kZhliGJa0bb00mua+mex3UYuaJYZhSKzUkSM4yd8jTaPdLaQUM8TIdWRz2UjNgNSu8TlDmOw2EAOGtK5xShf+St4WlgA5OQZj5NYyBIB7MTvp6W8Xv28hf8RFaoa7BxFROuYnh05JYoYhsVxL5u5RSvJ3jgWdZFCyxDBikuts5I3lB0I6BsTa+ivTKumlFjFDht5R6jYwJDb7uXsmvdQiBonN/vz4x11qYqAk5ygRJUepu01LcjFDDKumFvcgP5Z+srEgsbmb/A4dycUMiZVm6h1PE9EkdesYgtxBg/w9/WaGwJjJP8+aFCPm1hIDltqqTanK+k1264hhtOTsdZFBLiaGVVObd5T6XQaSixli1TbyTxaNS9xCrm4zNx5WS208USd7DcRATo5Bjq8cchBj6KUkpcnfwhDkNurtbctfZ/PHKowySG0siNG5Kf9p2k9mMYjRuilWcl8imgySw6AklmGIYZXU4uqnJBH129gSnvGmidLjXrLl+ERpIGeOf9Pj5LVxa8IwLKnN3Tx951uXhWKcXOqiOcXcpGUYXj5GqdflpfEkP+8Wb55356VLmexWWv9xF/+3YJkcBt3kdHZTf5KTX38LP680UTLoJnf/JI27tQW+nstt1K0Tl42n/HNObgZdamLkNuov4tmF8xHHtwalmBixkgxOF5nlIEbXTZP5+/DrUXUdF8F//L6Q2yOT3pbCWGJ+cveOUnLQSUowpHXHSMDS4hubU02EO9He3s6/uhBOFH3ms8LyXBiS9hCGQu2QBOyw7osAmTA6rM8hBAsCQwFYJUM4aO1AKAMg2gV7+7tIaNrh87RDEQ4gWrFxYgZD4SFkmWZYraoia08Kk90IFgmEwxlYPR60yQFG60Bw6Nd48e3i913Ql1j0EXSZLDgYVWFfIIQeYwaBH5vQEc4A8R5YTa8hotqH8JAPpowHVmsPF0YUSaHSaKFRVfZhUiEfwplmmIwamEzNyIR8CKcASIwwNTO4Eg6Ai2jFEQpEkVWaYNFkEO4w4ce+DEyeIYT3qRDpsHBjgQgiEXAl1AWfogOeLguk5w6jfV9g3vBk1TVZVISgxvqJREA2Ap8vBXtXB5ozf8DbdhM6E3b07DMCfziEfV1Fv3clDN+QHp2eLpgkp/C2vQO+mpMQASIACR+6wnp0ejphFJ3CofZOLChSLZJAIspxUGk0wWe34OCQFO2hMDxG4N3n7OiMAM0dHjjUAJQ2+Ie6YBRFsM9ixaFEM3rCQ+hpjqPLYocvBcQ9dlgPRaHpDCPc1YzE4XbYe5YytBdHJA5AqoAU3L8KZBGPxgFoYO2wQLVQVzqTQSqRQCKRQDwcQDgOsM16VIvZSCQiAJkyugFSYzs6O+1oRgLxIR98kSzEej00IkBibEe7XlrTG0v4OrBvSAp7p73yb8dD8A2JYLTo0Wy0QJUIwTdUNooF8dEQ9lmfw9FUM/b5emDNhBBKVKFy1XVcBP/Fo4hnGShUUo5+UimkSCAaTQEKE9qtmjwvChA81bynikovpZNGKU3eFoaz7HMWrFlMjJr7LJ1MUpK3NgcdSgKjo+5J3tOBmMz+/IckRyVPNUm9hoJFXGlscsdg3gLPWaNz3o8XeaqjLlLnryOidIwGg0EajKX5cRUs98leAzH8mGsjSd4WMTFaFw1OTtLkoJO0jJjM/uRcS7rEsuW8UUbXzVvv/eSQ57xh3sNSu4ijyig51bn3tTyiWmuyGE+1xvr1GorWL7dGBu730n4yMyCG98y5cRVomPOEbP21PNXcdQbq5a/rt7El7+d4ZFonDcZiFIuNU7/bQCxAcls/pcvum04mKZmbVLCNxGCoxZsmIi7ikaf3qJPUxVGacRepGTH33TR3jzQRUbKXDPkIxiI81XE/OR0Ocjgc5OgeLPU2015qYUDivLsY5DzvomhHft7zeKrlcsvkox259ReT2RujWCxGsVEv2dQc3avxfM6jF6vbClGcEn6q4KmmB8mhBIkNvVRNlGLdOmIYA3WPT9JkLEg2OcOvXzEdF8BH4+XyzdOuoqdabR0XwX/9NmLBkMGbLNrLGNIVE3BU8FS/Knxjc6pK5yiinZqiQqVH8Bpv+adSALLvwioVgcvSZJEVxZFAAvF9VnT4hhBPAchmAUYKZIBUKgVAWsitSBWQMphjHQMi3nJOIMFdUjDAEwmkwEAhlQCLKS1JxZEAA0Uu/yNSoNnI2c6ReApZTMFnFMGXm0tWjmi89LcruKnwha8ge+U5PHLvc4X6K18YKZMJUqMFeuYphAIRJFQBDGWVaLdogEwAiRSQPdEOlagd/E8C8ThS4BNEUin/0xJIRQyymUp0KvXGqq/JouITCFVZP+TGw6+fRAJAIuX+hQiici9KJIVUmpuOhB9jZiHZMkj460QiSUUPKofsqdfwyH2v5d+L1Q549ukhwlBJ5CMeaId1XwCRRIajC6okJlMpJABMdTVD1MX/RhYQxeNIRUKw27sQiiaQ4f/OLjr8E4Dn4NuYAsDoVOiwN5d42iIRkMlw8xVluEUXiUSLz82pbfB2WiAVZZBJxRHo7MBTegswFICec/9weM99OJz7PiNHS1cXrFX4XWP3wN8cga9zH6xGEaSRHujnGVQq1AXfOTGMnaYqYhRHwDeEbDaLHz90L36cr7XyYSijL9x/IXyU4vYFaV6+Jflr5owr4plnHRfAf7k1yWTykQHuz4J7KhQq3RIk3EbKtKAr0gWjqEgZhjqgOjgEhSuCVIcKkXYFHunJhZgkACJIJLgwJ+JxJLKoUBwiQrOxGczhMDyeCKwdGn5TSSDQE8IUFLDoVQCGFjFkCaTIIhFP5SQLgcAQMioL9AoJGLCweMLY11wQCol0vk3Dh/AVOczdPbAqCqGlH4d8CKdMMEmNsOgZ2MM+eKJDyCrtsGh4xSQBGNU+hD2WwoYjkkIyT221SAQgk+KFPIFEYv41kS4m+hvaB3uV9Vt8sRk/PimQ4P8jkYjmt4UWsy8p29DbaYFUBIgkKmiaFXP5KeVDu/1tRPW9iMatkAYskLYGavAJILEHEGpXFes7eIwdOJyywJ/wwAQPjNKn8qtVfV3KoPcgQZ4qHyqgUTAIJBJIAFDE44iDgUqz+EI7RqJCs7EQ7m1OBeB5KoTAUIZXqmIYXB60a0QARFBomqGqUEgUDXTCMySFZZ8VJo0RqoQPgefCCEUAfc2qnwzCvhCmxHqY9FVi1vEAAkOA2uZFZ646KtoD699zIWB9TqsuhI8kUkiK5bvomnJF7mmvvo4L5j+FBiomi3iU+5FUPI4EpDCphCJFQaneEkTQW4wQvzsEnycCjSkD375ORFSd8Ol5qy6VQDQcQlcoBWQTiEQSsOqNUDEnEOraB5+oGZEuHxJAxYpLqWUf2rv0eO05I5ojFlg0IsTDPniOTkFu7kHHYsv4NBYY1V14zbcPnfoOKELtsB5KwdpvhUVvgob5e4R9IUQVGiQ8+9ATN6Ir0IHmVAB2UycSph4E2jUlnmHIF8YV1gK71VhkVRux77CPrwKWwmhqBtp70CnKQGm3gLtDM0xGOQ75AvAN6WGRhNG5LwBJuw8eS61JSKFQSIGjIXR1hZCR+NCT2w0y1dck0GWssE+kEPV1oTgFKlIYYUH19TMulk2yQ+jp8EBhycDT8wdkWTOMGgBLmIZkpBroTUYo5skvZsApvcRQAD09nDGWiAwhYVJwqeJoGIGQEVa9CUZlF7pCPoRNdigiXejwZWD1dCKTyXK5yugQfB4PIiIgFY0gktLXWJfFGAkqmCwq7NvXifaODFSRLpyCBm6jCkACQ4EhJDDEPVKWGEIoIIJEoYdRM1eCsokhBDweTrYyUQS6wsgyKjSrROBCFyJIm00w6ueR9EQIXa8NIRCNwNoMhHv+gKy4BRoFkImGEYqmkIikkEUG0XAAgbgUGmMzFKIIwkNXwKj00FTVqQEMZVVot1s4vuC0P5r3vZuvAl4wHyma0SwHDvn2oVPfDmmok6sDqKDsa63jgiE1wqIX4ylPO9qlRiR6wgWjORVBKBxHJh5FBlmkIiEEAlJIm01olgoK8S86p1q7+ncyX2kKRkxKnYP8MS6P4jYoScyISa5zkH/QS2Y5w1X1UZIG+UpRsVxHjl5nSQXd3JTlIHXbDKRmxcQwYmKVOmpzB+dUNS4op0plVZ1yLbV1j+ZzNzG/g3RKMTFgiFW3FCo+q1X/VqjuLa5AzVcBx7jrUZ5zSg6S26zmqlHF8qJK3bKqVT7nB7WLo3vMS21qlq9a7KVemzxftVp1TSrmVCvk3Qy9lKyxfv9WkgtNc/knto2Cxbk+s5/SuSpK1kxud67aW0sOvjR43pyquJCLLM7pzl/lSiX8WrhvkvqdOpLz1byu/kHqNrDEiNXkHCWK9baQnGGIkZvJmyRKj/dSm1ZOYgbEsGoy85XZk34baVmuWtvcPUj9Li2JGTHpumPzrMsikB6n3jYtN1a5lmzeXL0Bl7MuX7O5VekVcqqMmFi1gRxerpJ+vurv8upcv9PA8RQYYtWFauhxl3puzUVurXI59qq53xi5dVylbekwJrlcvdxG/enF8BHRZNBBOjnDVf+6usmmLtQtFPNVrXVcOP9xVb4OnZzEDEOs2kzduVL7fhvJ58gXU6glEXDbEPqpChAgQMDXhgwCFglaQyb4Uz6YBIJ84yEcUyhAgAABXxMSQz4EIlnukSuBHIJSFSBAgAABt+6lDnW14+24HLoO+xKeqCTg64QQ/hUgQIAAAQIET1WAAAECBAgQlKoAAQIECBAgKFUBAgQIECBAgKBUBQgQIECAAEGpChAgQIAAAYJSFfAVIgSrtA51dYWXSKqCsT3An442hHZFHURGz2KO9f+KkYDPJEWdwo5wBkgNdcGkkUIk4o6zs3ui/Nm0cfjseqgkIohEUmhMnfm2fiGrpGTOdXUimHwZxHv0kIg06IhUOdo+EUKHUQWpSFRGJwBxH+zNCkhEIkg1JnQNVaZYKtxZGK9KD7svN94MIh4rmhUSiEQSKJot6Klyj0SoA0YV9z2VsaPQ/iseQHtufBIV9MXjK6NhqMPI0UaiKP1eaghdJg13D0UzrHl6liHjg0kkgmlOf7Ia88hE4LHw95Y2w+KJVLl3FD67HgqJCCKJFBpLF9dqELXWuwhD7VDU1c1ZY6MnxY2vp3h8VvRUWe+qdM5E4cmNTySFxtRRtQWbAAG3DOFQqW8i+OPe5C3kdLvJ7XaRzSAnBrmm0eXHCX79SA86SJlrP5UOUpscBLmBHG4XmdViAqMld6zQ+kvb5iK3w0BygNi2IKVz7ePkLeT2+8nPvwYniSjNtaer3MZrkrwtYgKjJLPTTc4WJTEQU4s3mW9CDVbHjUOZO4KOKh77yCjN5HQ7qUXJt/yKFY7NZHU2crlspGWLj8ssvYdZDBJr28jlaiOtGMSa/ZSkSeo1iAliLTm8QfI6tCQGQ4a5/eQo1q0jMVjS2lzkatOSGLlWfknym1mCWE1mp4scOrb68XWjTlIy6rmt0GrMY9SpJgZyanG6yWmQExg1uUYr3VpNDFjS2VzksmkLre5qrHfpUvVTr9tNbv7FzZGltmCa0v02koMhucFBbpeZ1GIQ1IuhM/ENy1kyuPwU7G0jNcN9JhzQJ2ApISjVb7BSLTlblj/Xl9G6KZZXqi7yOnXcecZKM3XzPSspPU5eh4GULEMMw5JSZ8v3s5zs1hHDaMnpdVOLUkwMIyedI5hXziVnFCsN859DzPdBDbaxBDF3fi3F/OS02cjFX5vke6G2BdM02G0jm6OX3yz5eWjdFOP7gjJaN42n04X+o0UbOhjd3I2axqnXYSNb9yC3eY67SZs707mkjyxR0ttScl5r8TmqLoeTekeLaMT3zUz2u8lmc+bPMQ6axUV9Y4uPXTbnFUQxPfxpfo65/rUxbnzyOYMY53qr5s6OTqfzfWW5vptMvvdvOl3oOTtH5/QaiBHP7QVcfR7j5FIX8VrMTdqi3yru3zvY7SSnO8gbNoW+vJNV17v2ub7dOnH+dyeDbnI6u6k/14+3pXJP4+p0JvLn5pTmewEri3sDCxAgKFVBqZYc2M5vYlo3jfPKCKySDDY3dTu5BtniFi8lKU2DDiXnebiD1B/sJrMceQ+NOwicIVbnIG+/n5w6MQE5z2aUnGqGoGwj7+AoeduUxIh5BZL0k02rJYN7tOYB8pU9Sbbg+ZV4t05SM7nG3pxSAaskpZg7BJzVOflNlvKNmA29tX3zcbeOxDlPNX/NZNFh4zXukU5SbLyf3AaWwJrJX/61JO+RVfBUx91aYhg15XTRuEvNe2u8l8maqTeWpJjXTCxY3pOe2wyB1ZmpRSkmgCG5wUWDycLcteY2zoNjWNI5ghWbb/fb5MRo51EkJfMoa0bOH0YvrqERk5MxGvXbSM0wpJ7j0lZf75Jvec3EMkpyDKbLlmCSYoO9ZJaDxIbuOQ0LqtM55+kryRacpOS4mwxilDQdFyBAUKp/6UpV66TBWIxisXHqdxsK4bacUs2HxzirHGoXjVM/2eQo2Vi5kKucbP05pSrOd61Ie1uIyXWxGHWSurhD0LiL1IyYWuZrd5LsJUNF72uSgjY1F+51DZZsbulRNxlYEORm8k5yc3Dp5MQqDeTo9lK3Q0diFG2KVX+j4C2Pe82kZECsoZube7+NWDBkyCkwXjnpKsZNc7QAgdWSIzhZ7upxBgijJFswWTE0Whzy5EKRfBh10k9tSibfMURp9s7tbsPPD4yWbL1B8jp1JM6F+4NtJAYXnu72+8nVIidATo7+St4fQ+K2YHVFUj6PtJdaSpRomZKdO1OO18CQvKWbRtMLW+8ya4ocysrh/H4bSwBIrG4jb4wWR2eKUW+LPN8FSax2UDAp7CYCBKUqoFILrdwmMVnqGc5p0Zb2UwsDYgyFtlfJIsVZ3rKK27B5pdpvIxYgMAwx/AsotK+r4R6Slin/Xoy8ZjmXg3P1l+R+k4NO0vI5TG+s2tZbbhzwbe0qelBpGu02EFu+0S/WU50cpaC/l88ragt5xckg2dQMQawmW7CyQh53lXpQo3kPapwLc8rN1B3sp2C3mZRgSFvu4fHKLd+uLM0bVi1eSvLGgTanSfj8tbY8Fp4OUhvLcO3gKs6v0jwW66kmabzfT1435zWzZj+vGKuv95w7+M3E5vPFZZ+NBsnfy+W/GeXc/Hd1Oqdp0KkmhlFTW28/9ftdZGA5xR0TNhQBSwih+vebDGUbev1BBINB9A/GEI90wThfo2GRBFIpgFQcucLHRCIBQAqJZJ6O1RIJpACU9gCi0Sii0ShisSgCdtUiB55CuN0E62HA2BtGqENf6NAR98Bieg1RhQOhsA8WhShfuRnq6UJXKJ4rVuWbbEvm7bOdCFhh+nEYkjYfwgE7NLkLFBqomCziUY4SqXgcCUihUpX2C8lEfOjs6IAvpYHRZEVnhwnybASBcBzIRNBpsuBQohmuUBg9VRZAoVJAlE0gmuAGnojEAZECKoThG7oCid4Ou1EPo90OPZtFJBAurQAWqaBSAEgk+IruDDIZACIRRAoVVEwWqUQqTxvuozLKxCOIpiRQqRQVKnerzUMBjYIBEgmOX+JxxMFApVHMqUwOd3WgoycCqd4ES3sn2psZTIUCiNRa77kDQdgXwpRYD5O+8K1ooBMdHR7EVUaYrB3YZ1UB58IIRRZIZ2kUgdAfkFUY0W7VQ29qh7WZwZVwAEMpYSsRsHRYLpDgmwtGqoHeZIRiUVc1w2pSwnOwBx1dGrSroujqOgUoHbA2A1We5eCgMcGo7EJXyIewyQ5FpAsdvgysngDsCMBu6kTC1INAu6b0OqkEEmQRzW36kU609/wBWWULVKkQerpC3Nf0FqCzA0enxNCapIj4uhABAJEGJrsEkZ4OPBf3IdJhgXSoB54pBs0dem7+qRRSGUAiKduuM2Hsaz+M/2TUMKviCHR18buvEe0mIyx6MZ7ytKNdakSiJ4ys0g5L2fBFiMDX9RrioTgiVg0SPh/+ExLoVVLEPVbsO3UFrEEDDHnQNQQAEjRbrWguGopIb4GRPYxARzs6TRn4QlfAmqxoloqgkgInwl3oCgGqeA9CU4DEqEKpetbAYlWj67lOWNsBY8KDwBUx9MZmiBQaWJr34cc9dlglFogCHvwno4FdLy174iWCKDSoZAPFPe1V52GyqLBvXyfaOzJQRbpwChq4jXNvkgj34LV3fRiK2mEUDcETzoLRaKCKdMJUcb3tsGjKTaIIwkNXwKj00BTTLxFC12tDCEQjsDYD4Z4/ICtuQblur0pnkQIphRgIBdDl0cMqicAzlAUUCiiEnmsChEdqhPDv3EKlyoVBc8K/RHOrfw0O8uaqf2uFf4koPd5LbVo5iRkQw6rJ7B7kfqNm9W/ZePzmfF6r8GJI1/1bcqmZOWFtsDbq53/bpiv8tsFZVIxTrVBpks9Flt2TafFy4cyYnxw6OYkZhli1mbpH01XysTbSKVliwJBYriUzH8Lsd8jnjhcVHlkhosl+JxmUYmIYMSlbXPkiq+Sgm8w5morlpDW7CwVYJcMYJ6+Nq+ZmxEoyOPyF9Z8MktPAV2XLdWTzjs/JWY46lcQoHVQp61xzHulx6m3Tkpzh5m7zVgn3T/aTy6wmuZgh8FXl3vF0jfWuwL18eDkf5i5KF/idBlKzXMqBVRdVni+QzjQZJFdLbnxiUuraqqy3AAG3DqH1m4A/AjIIWRX4fkAPf9wH01fgGUQ6NPhWlwju6BDaFQLFBQgQ8PVAyKkK+CNABL3dAuWVAHp8X8ERNpkwenx/gFhvh0VQqAIECPgaIXiqAv5ISMBn0mBPxIT+aA/0oqW7c7xHD017CvbwEDqbRQKpBQgQIChVAQIECBAg4JsOIfwrQIAAAQIECEpVgAABAgQIEJSqAAECBAgQIChVAQIECBAgQICgVAUIECBAgIA/R6UaglVah7q6speoGV3RpfuVzFAHNCIJ9D1xAHH47HqoJCKIRFJoTJ0I5w++DaHDqIJUJIJIqoKxPcCf1DeEdlX5GPXoOReGXVEHqcmHak9cpoa6YNJIIRKJIFHoYfdEkeE+wVCXCSqpCCKJAnq7r/KpgJkofHY9FBIRRBIpNJYuhFPzjbcYAVgkc2kssYby47Pw45NqTOiqcvjp/N9LIWRVoK6uDpqOyDeE5YfQrqiDyOjBN+/I1wR8JinqFHaEM7X4LIOIx4pmhQQikQSKZgt6+LVLeYwQlfGF1B5GvEcPiUiDjkimym/XkKG4D/ZmBSTz8VO4szBelR52X268QCLUAaOKG6/K2IFQJeGKdqFZVMbXqg7uKMtaMlPGs5VlsIq8VxLQjA8mkQgmX/kPVKf7ksh+JgKPRcPLfjMsnggyVfgk1GHk1kqigD6/R2QQ6SkenxU9Vda71nrEfXboFYXPwsLZyaX42jqsyFvI6XaTO/fq9tLokrVhmqRegzjfI5RrbSYmbZuL3A4DyQFi24KUpknytnBtrsxONzlblMTkem1SP7WxILHOSX6/n3/1UyxNNNltIAZKqthlLM33opQbyOF2kVktzreiSvNdUJRmF7kdOmKrHNU26lRz3TxsLnLZtEUt3WqNt/RIt2AxbZ1cyzOlY5Drlcny43M5yCAHQd5GwfLT2hbwvVy/UwCkrnQu358kyo9w/OYgzXef0XVP1uQzGnWSEiBWZyOXy0ZattDjNdatIwZKauv25/k6OJokSveTQ1653VqhPWAlGeLayYHVceNQMnm5K+UnL7XwnYecbie1KJlCT9Wkl+swpG0jl6uNtHx3m+RcweC67zi8BZnsH6dkTZkpo2FVGawu71RpHEyFoyhr0P32ZT/3mZxanG6+U5KaXBXEjusby5LW5iJXm5bEfMef3BjkBge5XWau9656bu/fmusx6uJ6HBsc5OL3FaEn7dfe+m2ec2sng+QwKIllGGJYJbXwZ6ym+x2kZBhSO0cpTUST3haumXPv5NxzZ2Nu0jK5jT5Ng902sjl6eebhN1Wtm2I0Tr0OG9m6+d6O49x1Sudo/gxSFALUwQAADL5JREFUuWOQ0ukkJdNzmY6t1Jcy5ienzUYu/lzSJN+ftC3INYgGwzf1pkFyKis17k7SYLeTnO7cubZ883FdN03WGm/1rZhGXWpi2BauLynfJiwn0Nz4WGor16rzfm+c3Foxyc1m0jG3o1T5c4nl3Pm+/EG0XNNtd2zuOcU6G/WOp4vYxUktapYYhiFW3UKu3EGvVfiotlJNUr+rhdSseO79kqPU3ablzt1l1dSSO/OYXx+x2UVug5zEhv+H/u9a81nwfcrHl6ZgG0sQmzn+qcFnyX432WxO8vMCFjSL83zHNe02UG8yTclksoR/uV6kOnLH5vJQVRniZS/Hg1wbQflcgzPmJ5fDSb38WbuT3TpiwJKtnyjtN5M4z1uFefrLhavfRnKIyRxMUzpZLJO1ZKb8FlVksJa8l3NJr4GYCuOrRfelkH2XumjfjHHtFNVz2i6Ok1vLFOaeTufnMhl0k9PZzZ+FnCZ/C0Ng5s6j1noMOpTc/5PcvdNJQZ3+6SjVfINt/jWZpDQlOU9MrCNX/yB525QEKHmLME39DiUxYh11D3JWOtvC90JM+smm1ZLBPVr5UHia611Vsq7G3ToS5zw/vim0WK4klgEBYlK35Q4wn6ReA0OQVz6cvNhj9rawvEWeLLtmkrPwldXvkZyM0ajfRmqGIXUFk7RkvFV1VjfpxEW9TPkD8vMHzwfbKlvN83xvstdALGsm7zhHp9vxVDkvSJm3/EedykIPTIeSs87dQeoPdpNZjoInFOsmnZj3Lvr9XHNt1kz+ZC0+qqFUBx0kh5gMriANDgbJ3SInVuuiUUrzm6KaHP5RGuw2k5yRk60/nb8fxHIy2NzUGxynaI35LOY+6Soe9mRNPqsQcZAXPKZBB9fAXKkUc40F5AZyD6ZrNyWoJUOL7UmbTlJsvJ/cBpZfK76ReFEP1HFXaaPxAk+aSQyG5Eq+0TijJHNvKZ1qy0wNGawp73MVc6GPbzXhLab7Usj+Avvaprm+wqzOTC1KMdcs3uCiwWTxEkxSbLCXzHKQ2NA9Z47V14Pf95QGajNwayBWmkuMXAFfo1JFedcKxkx+Is4CzVk//Mbe4k3nTEGyKRlixGKC2EC9sSoBPoect8bL+G3UTQYWBLmZ89pKupBwoQzW0M0Jy2QvtchZkmvbyN3bS04DSyhqnFztN4o3uqBNzYXMXIOUniNIvKCxRR5NaSCInEqum0dJY+1q460smpyVyW9exV682OCm0dg4ec1yrtF4+Q5W63vJILXJxWTojeWNj4pKNc17FMWvdHqudz9ebHXzFrnWTbE5jchzClhOtv6ct8N5LkRE6clBCgb7aTxZi49qKNWiEF1vcJQm8wMteAwx/r1DnouM8MpQWRRGqzWfxdyn1D0iA+9J1eazEteJMzQYJdmCvD/sUBPLKqnF2cs1EmdATK57UdXfqCFDfETDkDPseCWr667cQSbtbeEUIqslB+9lcx5yQYmOu7XEYG5oM+1vIzkrJ7XZRb29rtIQ8rwyU0mZFcngPPJeHFnp1jEkrhSlqkr3JZB9vkl9QYmWKdkyPgGjJVtvkLxOHYlR6tH221gCQGJ1G3ljVDEEXXk9+KgSWNI5vRTstXH8UyVlICjVP7ZSVbZRrz9IwWDuNUqTlKbx3jbSyrnwG8NwjJVXqkQ06lJzDFEp55J3ssRciKKY1wadpOVzOt6SREmaRrsNxFYVxNxG6SI1ChvOuEtLDKMl93jlkKbXLOdyI/mw42Kt1SSN9/u5jY/PaUwuZry8YcBtksVeeZL6HWoS8+235Fp1qacx7/diNOhQkljrpMFkkpIxPvznGJwTMgu2iSu0ExOTOTiHuOTiQ1YxXiHp3DGitJ9aGBBjKLQB48KLXCu6UaeaUDG3XYuPaoV/JynoNJBSnDP05GRw9lOSHwfA8PdjiEHuHpU8yNrzWfh95hoe2pJwXyU+K6RRbGqGIFaTLThZNTXgN4uLjEPOy5nj/dSSocV6qpOjFPT38jlBLblG///2zhdOcWuL4z93cWEVuKQqcVAFVVAVniIOqhgH/TyRqQI5LlRB1VIFT4UqcFBFqjJPQdVkFakiT8EqMuo8kYThT8LM7s7bbfvu9/OJmM+Em5t7z+HmnvMLJ/Klp53RMmmnemHe0UPVcz7zgp3qM/5+opnIJJStSxz31/D9l+5Ug8X3UDpvH0YGj0rpbZczmgyD/DeTL/PfyfMR9lto0Czsp1k9/ptDRPTFXqmJCmxXKtGRR3Y3xm3rX3CUPhzfx25cg3AiQByh03UgyjIw7eDu3n/ZxdwR6tqPcCQdc2uMuvT0o+ve9Aba9xbSjTGsaQuHmsmuhUG/j6kTXsPfwQeQSj33g+07WLcabn4BKkML804ZQaWzNCQpC3gu3F2gpnTcRzApD+VMuWf1O+gMVsiWNdRvu7gtMvxnPsXqWn9jFXxjWI8ZVLQink5Lo9y/h/OwxHLtYdVRsEMWipI9+3TSeT6sexfv//0jvnnzBm+++h6/PQLvfvoGyq110kLxbg7bts+OObrF854q0DQFuJ9iPJ1ihSI0TQJSaWSzAHbuQWnteR6ALNLpFLJSGgwePC/UvDpzjAYjWO+esaNEsqh053B2W6yXM/QqwK/9DkZeGtk0wEp3sBwHjuPAWa/hjOo4rmKXesn9fFA7H2NngUq0q9Xxs1eEMbcwqGQPn1mNB+gPrHA8ffi+DyD1/DWTfEjKQ2GPcJ2gxZ3rwkMWinJa389fjdHtdDDe5VHRbtDtaBAfV5haLiRFQurRg+P5AHx4KxdISTg3Se9+hH5/hEhQ6/vBaKVT133m2KYTffCl/u6u4OzSUBQpRrmbNO6v4fsS8hIDPC+YO9eFCwblokq7AkUC4Hmhut0PximVwh/TLjqdEVylAu2mg7sbBXhnYX4m3E+ejzQkKQ34O+x2R22nU+BlLP4E6t9YoVIoehBKPbLtCbVVkQQwKrQXtNmHeSOxQbPNMtgJ5Npk78OcaqlA1TCnehqa3ZJZC0I5haZxpDZe0HofhBfBclQzjtSyk4cgX8citaJBjUIQzmmH+afE8O8yVMTKp+pmc7mnfShYydUM6jVLlIEQs0MMFb4QqaQbZLSrJLMwfHitvzG7kFlDuOzjdkK1DEgoNMno6aSKwTheKgCTztvTZmnTYrEIjplOBQYSG0Oy15+gp30wKMcEEgR2ZBvnOdUeVUUEhbajnCoDCSWdzENOtUqmc82OkneqG7NKGSFHTdOm5dImsykTYyoNN3uaNYKdlW7aZM8MqhYK1DDXyTvfhPv54HaSwr9X7CxQ+IIyqk7G4X9DsrdhTjdUl/Z0lUQGEqpmsKNLDP9e8aFIaX+s/o3mJ8YvhFzkTxlCJIY5tjUj3J01ZhfjsDVrQYRG1alnNKmUiRTpV3zm3CuSfPAZfz8R8TCVhps4+ULSuNOn+36Y22TsWP0bHylbGqGCWO+RUQtSA+pwE/aPkVzVyYgU/UL1LBVGV+fjSUHcpl67SiKLnyse/v2zLKq0pUW7RCJjJMgqGQub3qoZYkKO/mnUArVvaAH7RfB6Qa69pP2Z+vdUqPRARo5dhiHDPIrKLvO7rGrSnojWE51UOVCWCnKJmmYkirgiVJrUgpzRyRGFi7Zk96qUyzBigkwlfRYf6tssyKjlSBQYIVS8mg/7Qzg3qb+x+aOY0Mx60qSSKBBjAolR23FB7Ject/10odIhj5QLx+o47neu/lX1k36sZ+3DHGVytVCtm2xH7eWVxWv/QMNmiWSBEcAoI6ukR3O+talXywVqYkGkQmMYht6T2ku4nw9u51KotH3Gzha6GBN2D18B2dr0thaqjwWRcrUeRQLnZKHSFR8KDIX0kkhCOAdvY3MSe3owm1SSM8TASBALVDsKWW8WbVLlwNbkqvHUpzObDtTZjMAEEguNp2sl+UzMd0ySDyb7+/GzgRw81MUJmK6N+6f6fmSfjUJg12KBmmaCmmL/QGazFM6xTKoeCa7WNGmrwfiBUSanUjshNZA8H3tavm1QQQz8MZc4V/+//D1Lv7l9FJUf4N8tserkX7/93Rh16TtY2gzuqMJDH5zPgI/5jYR/TMuYuGNo6de/wqqTx9f9FHrOPW4lPuIczl/kF5U+A1IdrbKA3wcjWP7rN++NR5i+l1FvlfmCyvlMpFBu1SG/n2Iw9v4Ha7aFwfh3COUW6nxB5XA+mr9tkXL/voNieYB0fwWr9YrfEr6FlvItpnkTq2kdWW5DnM+Gh7GWx3crDQtngPIrPtG5gzLytzu0rHt0i/xRkcPhiyqHw+FwOF8YXqWGw+FwOBy+qHI4HA6HwxdVDofD4XD4osrhcDgcDocvqhwOh8Ph8EWVw+FwOJy/Cv8F8uGyTexSViEAAAAASUVORK5CYII=';
                            // A documentation reference can be found at
                            // https://github.com/bpampuch/pdfmake#getting-started
                            // Set page margins [left,top,right,bottom] or [horizontal,vertical]
                            // or one number for equal spread
                            // It's important to create enough space at the top for a header !!!
                            doc.pageMargins = [20, 80, 20, 15];
                            // Set the font size fot the entire document
                            doc.defaultStyle.fontSize = 11;
                            // Set the fontsize for the table header
                            doc.styles.tableHeader.fontSize = 11;
                            doc.styles.title = {
                                color: 'red',
                                fontSize: '12',
                                background: 'blue',
                                // alignment: 'center'
                            }
                            doc.styles.tableHeader = {
                                // fontSize: 14,
                                bold: true,
                                fillColor: "#003f8a",
                                alignment: "left",
                                color: "#ffffff"
                            }

                            // Create a header object with 3 columns
                            // Left side: Logo
                            // Middle: brandname
                            // Right side: A document title
                            doc['header'] = (function() {
                                return {
                                    columns: [{
                                            alignment: 'left',
                                            image: logo,
                                            width: 140,
                                            margin: [20, 15, 10, 30]
                                        },
                                        // {
                                        //     alignment: 'left',
                                        //     image: address,
                                        //     width: 285,
                                        //     margin: [173, 15, 100, 30]
                                        // },
                                    ],
                                    // margin: 20
                                }
                            });

                            // Create a footer object with 2 columns
                            // Left side: report creation date
                            // Right side: current page and total pages
                            doc['footer'] = (function(page, pages) {

                                return {
                                    columns: [{
                                            alignment: 'left',
                                            fontSize: [8],
                                            margin: [20, 0, 20, 0],
                                            text: ['© 2020 | Supernet ']
                                        },
                                        {
                                            alignment: 'center',
                                            fontSize: [8],
                                            margin: [20, 0, 20, 0],
                                            text: ['- Supernet Technologies - ']
                                        },
                                        {
                                            alignment: 'right',
                                            fontSize: [8],
                                            margin: [20, 0, 20, 0],
                                            text: [jsDate.toString()]
                                        }
                                    ],
                                    // margin: 20
                                }
                            });

                            // Change dataTable layout (Table styling)
                            // To use predefined layouts uncomment the line below and comment the custom lines below
                            // doc.content[0].layout = 'lightHorizontalLines'; // noBorders , headerLineOnly
                            var objLayout = {};
                            // objLayout['hLineWidth'] = function(i) { return .5; };
                            // objLayout['vLineWidth'] = function(i) { return .5; };
                            objLayout['hLineColor'] = function(i) { return '#aaa'; };
                            objLayout['vLineColor'] = function(i) { return '#aaa'; };
                            // objLayout['paddingLeft'] = function(i) { return 4; };
                            // objLayout['paddingRight'] = function(i) { return 4; };
                            doc.content[0].layout = objLayout;
                            doc.content[0] = [{
                                columnGap: 20,
                                columns: [{
                                        width: '100%',
                                        style: 'tableExample',
                                        color: '#444',
                                        alignment: "center",
                                        bold: true,
                                        fontSize: 18,
                                        margin: [0, 10, 0, 10],
                                        table: {
                                            widths: ['100%'],
                                            body: [
                                                [fileName],
                                            ]
                                        },
                                        layout: {
                                            paddingTop: function(i, node) { return 10; },
                                            paddingBottom: function(i, node) { return 10; },
                                        }
                                    }

                                ]
                            }];
                            console.log("download end");
                        }
                    }
                ],
                "bStateSave": true,
                //"stateDuration": 0,
                "pagingType": "full_numbers",
                "order": [
                    [order_by_column, order_by]
                ],
                "columnDefs": [{
                    "targets": 'no-sort',
                    "orderable": false,
                    "className": "div.dt-buttons",
                }],
                "language": {
                    "url": "/static/vendors/dataTable/i18n/" + hidden_selected_language + ".json"
                },
                "destroy": true,
                "bProcessing": true, //Feature control the processing indicator.
                "searching": true,
                "serverSide": true, //Feature control DataTables' server-side processing mode.
                "search": {
                    "regex": true
                },
                // Load data for the table's content from an Ajax source
                "ajax": {
                    "url": get_url,
                    "type": "POST",
                    "data": function(data) {
                        data.order_by_column = order_by_column;
                        data.order_by = order_by;

                        /*form Search Fields */
                        var search_form = $("#search_form");
                        var data1 = {};
                        $.each($('input, select ,textarea', '#search_form'), function(k) {
                            //console.log(k + ' ' + $(this).attr('name'));
                            var attrName = $(this).attr('name');
                            var attrValue = $(this).val();
                            /*if(attrValue && attrValue!=''){
                             data1[attrName] = attrValue;
                             }*/

                            data1[attrName] = attrValue;
                            //console.log(data1);

                        });
                        $.extend(data, data1);
                        //console.log(data);
                    },
                    "dataSrc": function(jsonData) {
                        $(".form-loader").hide();
                        //console.log(jsonData);
                        return jsonData.data;
                    },
                    "error": function(xhr, resp, text) {
                        if(resp == "parsererror"){
							window.location.href = "/";
		                }
                    },
                },
            });
        } else {
            $('.listTable').DataTable({
                "lengthMenu": [
                    [10, 25, 50, 100, 250, -1],
                    [10, 25, 50, 100, 250, allTextLabel]
                ],
                "bStateSave": true,
                "destroy": true,
                "pagingType": "full_numbers",
                "order": [
                    [order_by_column, order_by]
                ],
                "columnDefs": [{
                    "targets": 'no-sort',
                    "orderable": false,
                }],
                "language": {
                    "url": "/static/vendors/dataTable/i18n/" + hidden_selected_language + ".json"
                }
            });
        }

        $('#Search').click(function() { //button filter event click

            var data1 = [];
            $.each($('input, select ,textarea', '#search_form'), function(k) {
                //console.log(k + ' ' + $(this).attr('name'));
                var attrName = $(this).attr('name');
                //var attrValue = $(this).val();
                var attrValue = $.trim($(this).val());
                if (attrValue && attrValue != '') {
                    data1.push(attrValue);
                }
            });
            //console.log(data1.length);
            if (data1.length && data1.length > 0) {
            	var daterange = $('.simple-date-range-picker').val();
				if($('.simple-date-range-picker').length > 0){
			var daterangearr = daterange.split("-");
			var nodays = datediff(parseDate(daterangearr[0]), parseDate(daterangearr[1]));
			if (nodays > 95){
				$.alert({
                    title: searchErrorHeader,
                    type: 'blue',
                    content: dateRangeErrContent,
                    closeIcon: true,
                    closeIconClass: 'fa fa-close',
                    buttons: {
                        Ok: {
                            text: okButtonText,
                            btnClass: 'btn-red btn-round',
                            action: function() {

                            }
                        },
                    }
                });
                return false;
			}
			}
                table.clearPipeline();
                table.ajax.reload(); //just reload table
                
            } else {
                $.alert({
                    title: searchErrorHeader,
                    type: 'blue',
                    content: searchErrorContent,
                    closeIcon: true,
                    closeIconClass: 'fa fa-close',
                    buttons: {
                        Ok: {
                            text: okButtonText,
                            btnClass: 'btn-red btn-round',
                            action: function() {

                            }
                        },
                    }
                });
                return false;
            }
        });

        //$("span.select2-selection__clear").remove(); //Select2 clear button remove
        $('#Reset').click(function() { //button reset event click
            localStorage.clear();

            var data1 = [];
            $.each($('input, select ,textarea', '#search_form'), function(k) {
                var attrId = $(this).attr('id');
                var attrType = $(this).attr('type');
                var attrValue = $(this).val();
                if (attrValue && attrValue != '') {
                    data1.push(attrValue);
                }
                if ($("#" + attrId).css('display')) {
                    if ($("#" + attrId).css('display').toLowerCase() != 'none') {
                        $("#" + attrId).val(null).trigger("change");
                        $("span.select2-selection__clear").remove();
                    }
                }
            });
            // $("select2").val('').trigger();
            //console.log(data1);
            $('#search_form')[0].reset();
            table.clearPipeline();
            
            // if (data1.length && data1.length > 0) {
            //     table.clearPipeline();
            //     table.ajax.reload(); //just reload table
            // }
            // $('#search_form')[0].reset();
            //table.ajax.reload(); //just reload table
            
            //Datepicker reinitialize after reseted value
            var language = $("#language").val();
    var TodayLabel, YesterdayLabel, thisMonthLabel, lastMonthLabel, last7DaysLabel, last30DaysLabel, customRangeLabel, applyLabel, cancelLabel, daysOfWeek, monthNames, last3MonthLabel;
    if (language == "English") {
        TodayLabel = "Today";
        YesterdayLabel = "Yesterday";
        thisMonthLabel = "This Month";
        lastMonthLabel = "Last Month";
        last3MonthLabel = "Last 3 Month";
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
        last3MonthLabel = "Le 3 mois dernier";
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
            $('.simple-date-range-picker, .date-range-picker').daterangepicker({
        autoUpdateInput: true,
        showDropdowns: true,
        linkedCalendars: false,
        startDate: moment().subtract(2, 'month').startOf('month').format('YYYY/MM/DD'),
        ranges: {
            [TodayLabel]: [moment(), moment()],
            [YesterdayLabel]: [moment().subtract(1, 'days'), moment().subtract(1, 'days')],
            [last7DaysLabel]: [moment().subtract(6, 'days'), moment()],
            [last30DaysLabel]: [moment().subtract(29, 'days'), moment()],
            [thisMonthLabel]: [moment().startOf('month'), moment().endOf('month')],
            [lastMonthLabel]: [moment().subtract(1, 'month').startOf('month'), moment().subtract(1, 'month').endOf('month')],
            [last3MonthLabel]: [moment().subtract(2, 'month').startOf('month'), moment().endOf('month')]
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
    table.ajax.reload(); //just reload table
        });

        /*When hit Enter key, Search action should work */
        $("#search_form").keypress(function(e) {
            if (e.which == 13) {
                //            alert('You pressed enter!');
                $("#Search").click();
                e.preventDefault();
            }

        });
    }

});