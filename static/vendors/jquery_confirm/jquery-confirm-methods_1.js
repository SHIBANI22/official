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
 * @file          : jquery-confirm-methods_1.js
 */

$(document).ready(function () {

    $('.item .delete').click(function () {

        var elem = $(this).closest('.item');

        $.confirm({
            'title': 'Delete Confirmation',
            'message': 'You are about to delete this item. <br />It cannot be restored at a later time! Continue?',
            'buttons': {
                'Yes': {
                    'class': 'blue',
                    'action': function () {
                        elem.slideUp();
                    }
                },
                'No': {
                    'class': 'gray',
                    'action': function () {}	// Nothing to do in this case. You can as well omit the action property.
                }
            }
        });

    });

});