/*
 * @created_on    : Friday 06 December 2019 01:28:06+0530
 * @package       : SUNU DIGITASS
 * @author1       : Supernet Dev Team (RAMAKRISHNA NADIMPALLY)
 * @author2       : Supernet Dev Team (JITHIN KR)
 * @author3       : Supernet Dev Team (RAJENDIRAN SELVAM)
 * @copyright     : Copyright (c) Supernet Technologies India Pvt Ltd (http://supernet-india.com/)
 * @license       : Supernet Technologies India Pvt Ltd
 * @version       : 3.0
 * @since         : APRIL 2017 to PRESENT
 * @project       : SUNU DIGITASS
 * @file          : jquery-confirm_1.js
 */

(function ($) {

    $.confirm = function (params) {

        if ($('#confirmOverlay').length) {
            // A confirm is already shown on the page:
            return false;
        }

        var buttonHTML = '';
        $.each(params.buttons, function (name, obj) {

            // Generating the markup for the buttons:

            buttonHTML += '<a href="#" class="button ' + obj['class'] + '">' + name + '<span></span></a>';

            if (!obj.action) {
                obj.action = function () {};
            }
        });

        var markup = [
            '<div id="confirmOverlay">',
            '<div id="confirmBox">',
            '<h1>', params.title, '</h1>',
            '<p>', params.message, '</p>',
            '<div id="confirmButtons">',
            buttonHTML,
            '</div></div></div>'
        ].join('');

        $(markup).hide().appendTo('body').fadeIn();

        var buttons = $('#confirmBox .button'),
                i = 0;

        $.each(params.buttons, function (name, obj) {
            buttons.eq(i++).click(function () {

                // Calling the action attribute when a
                // click occurs, and hiding the confirm.

                obj.action();
                $.confirm.hide();
                return false;
            });
        });
    }

    $.confirm.hide = function () {
        $('#confirmOverlay').fadeOut(function () {
            $(this).remove();
        });
    }

})(jQuery);