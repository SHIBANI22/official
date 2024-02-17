var current_url = window.location.origin;
//console.clear();
//block special characters with number
function blockDisabledChar(e) {
    var k = e.keyCode;
    if (k == 60 || k == 62 || k == 39 || k == 34 || k == 40 || k == 41 || k == 37 || k == 96 || k == 33 || k == 47 || k == 59) {
        return false
    }
}

function blockSpecialCharwithNumber(e) {
    var k = e.keyCode;
    return ((k > 64 && k < 91) || (k > 96 && k < 123) || k == 8 || k == 32 || k == 39);
}
//block special characters only
function blockSpecialChar(e) {
    var k = e.keyCode;
    return ((k > 64 && k < 91) || (k > 96 && k < 123) || k == 8 || (k >= 48 && k <= 57));
}

$(window).on("load", function() {
    var language = $("#language").val();
    if (language == "French" || language == "FR") {
        //  alert("Page is loaded");
        $("#signin-title").text("Se connecter");
        $("#signin-button").text("Se connecter");
        $("#lang").text("Choisir la Langue");
        $("#uname").text("Nom Utilisateur");
        $("#username").attr("placeholder", "Nom Utilisateur");
        $("#pass").text("Mot de passe");
        $("#password").attr("placeholder", "Mot de passe");
        $("#forgot").text("Mot de Passe Oublié ?");
        $("#eng").text("Anglais");
        $("#fr").text("Français");

        $("#title").text("Mot de Passe Oublié");
        $("#uname").text("Nom Utilisateur");
        $("#submit-button").text("Envoyer");
        $("#back").text("Retour à la Connection");

        /*Begin:: Common Change Password*/
        $("#oldPass").text("Ancien mot de passe");
        $("#newPass").text("Nouveau mot de passe");
        $("#input_old_password").attr("placeholder", "Ancien mot de passe");
        $("#input_new_password").attr("placeholder", "Nouveau mot de passe");
        $("#input_confirm_password").attr("placeholder", "Confirmer mot de passe");
        $("#confirmPass").text("Confirmer mot de passe");
        $("#submit").text("Valider");
        $("#msg").text("Le Mot de passe doit contenir les éléments suivants:");
        $("#capital").text("Une lettre Majuscule");
        $("#letter").text("Une lettre minuscule");
        $("#number").text("Un chiffre");
        $("#length").text("6 Caractères au minimun");
        $("#special-char").text("Un caractère spécial");
        $(".site-title").text("Portail Admin Supernet");
        /*End:: Common Change Password*/

    } else {
        $("#signin-title").text("Sign In");
        $("#signin-button").text("Sign In");
        $("#lang").text("Select Language");
        $("#uname").text("User Name");
        $("#username").attr("placeholder", "User Name");
        $("#pass").text("Password");
        $("#password").attr("placeholder", "Password");
        $("#forgot").text("Forgot password ?");
        $("#eng").text("English");
        $("#fr").text("French");

        $("#title").text("Forgot Password");
        $("#uname").text("User Name");
        $("#submit-button").text("Submit");
        $("#back").text("Back To Login");

        /*Begin:: Common Change Password*/
        $("#oldPass").text("Old Password");
        $("#newPass").text("New Password");
        $("#input_old_password").attr("placeholder", "Old Password");
        $("#input_new_password").attr("placeholder", "New Password");
        $("#input_confirm_password").attr("placeholder", "Confirm Password");
        $("#confirmPass").text("Confirm Password");
        $("#submit").text("Submit");
        $("#msg").text("Password must contain the following:");
        $("#letter").text("A Lowercase letter");
        $("#capital").text("A Capital(Uppercase)Letter");
        $("#number").text("A Number");
        $("#length").text("Minimum 6 Characters");
        $("#special-char").text("A Special Character");
        /*End:: Common Change Password*/
        $(".site-title").text("Supernet Admin Portal");

    }

    if ($("#LogoutUser").val()) {
        var logout_url = $("#LogoutUser").data("logout-url");
        setTimeout(function() {
            window.location = window.location.origin + "/" + logout_url;
        }, 8000);
    }
});

$("#language").change(function() {
    var x = $(this).val();
    if (x == "English") {
        $("#signin-title").text("Sign In");
        $("#signin-button").text("Sign In");
        $("#lang").text("Select Language");
        $("#uname").text("User Name");
        $("#pass").text("Password");
        $("#forgot").text("Forgot password ?");
        $("#eng").text("English");
        $("#fr").text("French");
        if ($("#username-error").length) {
            $("#username-error").text("Please Enter User Name");
        }
        if ($("#password-error").length) {
            $("#password-error").text("Please Enter Password");
        }
        $("#username").attr("placeholder", "User Name");
        $("#password").attr("placeholder", "Password");

        /*Begin:: Forgot Password*/
        $("#title").text("Forgot Password");
        $("#submit-button").text("Submit");
        $("#back").text("Back To Login");
        /*End::Forgot Password*/

        /*Begin:: Common Change Password*/
        $("#oldPass").text("Old Password");
        $("#newPass").text("New Password");
        $("#input_old_password").attr("placeholder", "Old Password");
        $("#input_new_password").attr("placeholder", "New Password");
        $("#input_confirm_password").attr("placeholder", "Confirm Password");
        $("#confirmPass").text("Confirm Password");
        $("#submit").text("Submit");
        $("#msg").text("Password must contain the following:");
        $("#letter").text("A Lowercase letter");
        $("#capital").text("A Capital(Uppercase)Letter");
        $("#number").text("A Number");
        $("#length").text("Minimum 6 Characters");
        $("#special-char").text("A Special Character");

        $("#input_old_password-error").text("Please Enter Old Password");
        $("#input_new_password-error").text("Please Enter New Password");
        $("#input_confirm_password-error").text("Please Enter Confirm Password");
        /*End:: Common Change Password*/
        $(".site-title").text("Supernet Admin Portal");

    } else if (x == "French") {
        //  $("#demo").text("Portail Admin Supernet");
        $("#signin-title").text("Se connecter");
        $("#signin-button").text("Se connecter");
        $("#lang").text("Choisir la Langue");
        $("#uname").text("Nom Utilisateur");
        $("#pass").text("Mot de passe");
        $("#forgot").text("Mot de Passe Oublié ?");
        $("#eng").text("Anglais");
        $("#fr").text("Français");
        if ($("#username-error").length) {
            $("#username-error").text("Saisir le Utilisateur SVP");
        }
        if ($("#password-error").length) {
            $("#password-error").text("Veuillez saisir le Mot de Passe, SVP");
        }
        $("#username").attr("placeholder", "Nom Utilisateur");
        $("#password").attr("placeholder", "Mot de passe");

        $("#title").text("Mot de Passe Oublié");
        $("#submit-button").text("Envoyer");
        $("#back").text("Retour à la Connection");

        /*Begin:: Common Change Password*/
        $("#oldPass").text("Ancien mot de passe");
        $("#newPass").text("Nouveau mot de passe");
        $("#input_old_password").attr("placeholder", "Ancien mot de passe");
        $("#input_new_password").attr("placeholder", "Nouveau mot de passe");
        $("#input_confirm_password").attr("placeholder", "Confirmer mot de passe");
        $("#confirmPass").text("Confirmer mot de passe");
        $("#submit").text("Valider");
        $("#msg").text("Le Mot de passe doit contenir les éléments suivants:");
        $("#capital").text("Une lettre Majuscule");
        $("#letter").text("Une lettre minuscule");
        $("#number").text("Un chiffre");
        $("#length").text("6 Caractères au minimun");
        $("#special-char").text("Un caractère spécial");

        $("#input_old_password-error").text("Veuillez saisir l'ancien Mot de Passe");
        $("#input_new_password-error").text("Veuillez saisir le Nouveau Mot de Passe");
        $("#input_confirm_password-error").text("Veuillez saisir le Confirmer mot de passe");
        /*End:: Common Change Password*/
        $(".site-title").text("Portail Admin Supernet");
    }
});

$(document).ready(function() {
    var language = $("#language").val();
    if (language == "French" || language == "FR") {
        var select_placeholder = "- - - - Sélectionnez svp - - - -";
    } else {
        var select_placeholder = "- - - -Please Select - - - -";
    }
    $(".select2").select2({
        placeholder: select_placeholder,
        //allowClear: true
    });
    $(".alert-success").fadeTo(4500, 4500).slideUp(4500, function() {
        $(".alert-success").slideUp(4500);
    });

    $(".profile-success").fadeTo(4000, 4000).slideUp(4000, function() {
        $(".profile-success").slideUp(4000);
    });

    //dashboard sidebar data
    /*$.post(current_url + "/transaction_data", { status: status }, function(data, status) {
        console.log(data.user_firstname);
        if (data.total_success_amount == "null" || data.total_success_amount == "") {
            $("#total_success_amount").html("0 XOF");
        } else {
            $("#total_success_amount").html(data.total_success_amount + " XOF");
        }
        $("#total_transaction_count").html(data.total_transaction);
        $("#total_user_count").html(data.total_user_count);
        $("#total_entity_user_count").html(data.total_entity_user_count);
        $("#user_fullname_text1").text(data.user_firstname + " " + data.user_lastname);
        $("#user_role_text1").text(data.user_rolename)
    });*/

    $('#username').bind("cut copy paste", function(e) {
        e.preventDefault();
    });
    $('#password').bind("cut copy paste", function(e) {
        e.preventDefault();
    });

    $(".navigation a").each(function() {
        var href = $(this).attr("href");
        var hrefArr = href.split("/");
        if ((window.location.pathname.indexOf($(this).attr("href"))) > -1) {
            $(this).addClass("active");
            $(this).parents(".navigation-menu-body").find("ul").removeClass("navigation-active");
            $(this).parents("ul").addClass("navigation-active");
            var activeMenuBodyId = $(this).parents("ul").attr("id");
            $(this).parents(".navigation").find(".navigation-icon-menu li").removeClass("active");
            $(this).parents(".navigation").find(".navigation-icon-menu li a[href='#" + activeMenuBodyId + "']").parents("li").addClass("active");
        } else if ((window.location.href.indexOf(hrefArr[1])) > -1 && hrefArr[1] != "Report") {
            $(this).addClass("active");
            $(this).parents(".navigation-menu-body").find("ul").removeClass("navigation-active");
            $(this).parents("ul").addClass("navigation-active");
            var activeMenuBodyId = $(this).parents("ul").attr("id");
            $(this).parents(".navigation").find(".navigation-icon-menu li").removeClass("active");
            $(this).parents(".navigation").find(".navigation-icon-menu li a[href='#" + activeMenuBodyId + "']").parents("li").addClass("active");
        }
    });

});



//profile pic js
$("#image").hide();
var language = $("#language").val();
$("#removePic").click(function() {
    $("#joint-img").attr('src', null);
    $("#image").hide();
    $(".filename").text("");
    if (language == "English") {
        $(".choosefile").text("Choose Photo");
    } else {
        $(".choosefile").text("Choisir la photo");
    }
    $(".joint").val("");
})

if (language == "English") {
    $(".choosefile").text("Choose Photo");
} else {
    $(".choosefile").text("Choisir la photo");
}
$(".joint").change(function() {
    if (this.files && this.files[0]) {
        var fileName = this.files[0].name;
        var reader = new FileReader();
        reader.onload = function(e) {
            $("#joint-img").attr('src', e.target.result);
            $(".choosefile").text(fileName);
            $("#image").show();
        };
        reader.readAsDataURL(this.files[0]);
    }
});

/*function validateImage(element) {
    var file = document.getElementById(element.id);
    var file = document.getElementById(element.id).files[0];
    var filesize = file.size;
    var validFileType = /.+\.(jpg|jpeg|png|gif)$/i;
    if (filesize > 2100000) { // it wont take 1.9MB onwards
        alert("Profile Picture Should be 2MB");
        $("#" + element.id).val('');
        return false;
    } else {
        return true;
    }

}*/

$("#deletePic").hide();

function mouseOverDelete() {
    $("#deletePic").show();
}

function mouseOverDeleteBtnRemove() {
    $("#deletePic").hide();
}

$(".updateprofileTab").hide();
$(".editProfileBtn").click(function() {
    $(".updateprofileTab").show();
    $(".userInfoTab").hide();

})
$(".backProfileBtn").click(function() {
    $(".updateprofileTab").hide();
    $(".userInfoTab").show();

})

$(document).ready(function() {
    if ($(".dataTables_scrollBody .dataTables_empty").length == 1) {
        $(".dataTables_scrollHeadInner").css("width", "100%");
        $(".dataTables_scrollHeadInner .table").css("width", "100%");
    }
    $(".navigation-toggler-icon").click(function() {
        if ($(this).hasClass("toggle-icon")) {
            $(".channelTable").removeClass("table-responsive");
            $(".navigation").css("width", "90px");
            $(".navigation").addClass("navigation-small")
            $(".header-logo").css("width", "90px");
            $(".header-logo .site-logo").css("margin-left", "0px");
            $(".large-logo").css("display", "none");
            $(".small-logo").css("display", "flex").css("width","66px");
            $(".header .header-logo a img").css("margin-right", "0px");
            $(".header .header-logo a img").css("margin-top", "67px");

            $(".navigation-toggler-icon").removeClass("toggle-icon");
            $(".navigation-toggler-icon").addClass("toggle-icon-small");
            // $(".navigation-toggler-icon i").removeClass("ti-menu");
            // $(".navigation-toggler-icon i").addClass("fa fa-arrow-right");
            $(".main-content").css("margin-left", "89px");
            $(".navigation-icon-menu li.active").addClass("active-menu");


        } else {
            $(".channelTable").addClass("table-responsive");
            $(".navigation").css("width", "300px");
            $(".navigation").removeClass("navigation-small")
            $(".header-logo .site-logo").css("margin-left", "0px");
            $(".header-logo").css("width", "300px");
            $(".large-logo").css("display", "flex");
            $(".small-logo").css("display", "none");
            $(".header .header-logo a img").css("margin-right", "0px");
            $(".header .header-logo a img").css("margin-top", "67px");

            $(".navigation-toggler-icon").removeClass("toggle-icon-small");
            $(".navigation-toggler-icon").addClass("toggle-icon");
            // $(".navigation-toggler-icon i").addClass("ti-menu");
            // $(".navigation-toggler-icon i").removeClass("fa fa-arrow-right");
            $(".main-content").css("margin-left", "300px");
            $(".large-logo").css("margin-top", "20px");
            $(".navigation-icon-menu li.active").removeClass("active-menu");
        }
    });

    $(".navigation a").on("click", function() {
        if ($(".navigation").hasClass("navigation-small")) {
            $(".navigation").css("width", "300px");
            // $(".navigation").removeClass("navigation-small")
            $(".header-logo").css("width", "300px");
            $(".large-logo").css("display", "flex");
            $(".small-logo").css("display", "none");
            $(".header .header-logo a img").css("margin-right", "0px");
            $(".header .header-logo a img").css("margin-top", "67px");

            $(".navigation-toggler-icon").removeClass("toggle-icon-small");
            $(".navigation-toggler-icon").addClass("toggle-icon");
            // $(".navigation-toggler-icon i").addClass("ti-menu");
            // $(".navigation-toggler-icon i").removeClass("fa fa-arrow-right");
            $(".main-content").css("margin-left", "300px");
            $(".large-logo").css("margin-top", "20px");
            $(".navigation-icon-menu li.active").removeClass("active-menu");
        }
    });
});

$(document).on("keyup keypress blur change mouseup mouseover mouseout hover", function(e) {
    var container = $(".navigation-small");

    // if the target of the click isn't the container nor a descendant of the container
    if (!container.is(e.target) && container.has(e.target).length === 0) {
        if ($(".navigation").hasClass("navigation-small")) {
            $(".navigation").css("width", "90px");
            // $(".navigation").addClass("navigation-small")
            $(".header-logo").css("width", "90px");
            $(".large-logo").css("display", "none");
            $(".small-logo").css("display", "flex");
            $(".header .header-logo a img").css("margin-right", "0px");
            $(".header .header-logo a img").css("margin-top", "67px");

            $(".navigation-toggler-icon").removeClass("toggle-icon");
            $(".navigation-toggler-icon").addClass("toggle-icon-small");
            // $(".navigation-toggler-icon i").removeClass("ti-menu");
            // $(".navigation-toggler-icon i").addClass("fa fa-arrow-right");
            $(".main-content").css("margin-left", "89px");
            $(".navigation-icon-menu li.active").addClass("active-menu");
        }
    }

    // $(document).bind("contextmenu", function(e) {
    //     return false;
    // });



});

$(document).ready(function() {
    $("body").tooltip({ selector: '[data-toggle=tooltip]' });
});