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
 * @file          : jquery-confirm-methods.js
 */

$.confirm({
    buttons: {
        hey: function () {
            // here the button key 'hey' will be used as the text.
            $.alert('You clicked on "hey".');
        },
        heyThere: {
            text: 'hey there!', // With spaces and symbols
            action: function () {
                $.alert('You clicked on "heyThere"');
            }
        }
    }
});