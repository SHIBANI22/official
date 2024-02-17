$(document).ready(function() {

    var language = $("#language").val();
    // alert(language)

    var validationMsg = {};
    if (language == "English" || language == "ENGLISH") {
        validationMsg.UserFirstnamErreMsg = "Please Enter First Name";
        validationMsg.UserLastnameErrMsg = "Please Enter Last Name";
        validationMsg.MobileErrMsg = "Please Enter Mobile Number";
        validationMsg.EmailErrMsg = "Please Enter Email Id";
        validationMsg.ValidEmailErrMsg = "Email incorrecte";
        validationMsg.EmployeeIdErrMsg = "Please Enter Employee Id";
        validationMsg.RoleErrMsg = "Please Select Role";
        validationMsg.DepartmentErrMsg = "Please Enter Department";
        validationMsg.StatusErrMsg = "Please Select Status";
        validationMsg.LanguageErrMsg = "Please Select Language";
        validationMsg.RoleInputErrMsg = "Please Enter Role";
        validationMsg.PlaseEnterConfirmPassAsNewPass = "Enter Confirm Password Same as New Password";
        validationMsg.MobileMinlengthErrMsg = "Please enter valid 10 digit Mobile Number"
        validationMsg.PackagenameErrMsg="Please enter Package Name";
        validationMsg.volumeamtErrMsg="Please Enter Volume";
        validationMsg.CodeErrMsg="Please Enter Code"
        validationMsg.SymbolErrMsg="Please Enter Symbol"
        validationMsg.CountryMsg="Please Enter country"
        validationMsg.TransactionFeeErrMsg="PLease Enter Transaction Fees"
    } else {
        validationMsg.UserFirstnamErreMsg = "Saisir le Prénom SVP";
        validationMsg.UserLastnameErrMsg = "Saisir le 2nd Prénom SVP";
        validationMsg.MobileErrMsg = "Saisir le Numéro de Mobile SVP";
        validationMsg.EmailErrMsg = "Saisir l'Email";
        validationMsg.ValidEmailErrMsg = "Email incorrecte";
        validationMsg.EmployeeIdErrMsg = "Sélectionner ID employé SVP";
        validationMsg.RoleErrMsg = "Sélectionner un rôle SVP";
        validationMsg.DepartmentErrMsg = "Saisir le Département SVP";
        validationMsg.StatusErrMsg = "Sélectionner le Statut  SVP";
        validationMsg.LanguageErrMsg = "Sélectionner la Langue SVP";
        validationMsg.RoleInputErrMsg = "Saisir le nom du rôle SVP";
        validationMsg.CodeErrMsg="Veuillez entrer le code"
        validationMsg.SymbolErrMsg="Veuillez saisir le symbole"
        validationMsg.CountryMsg="Veuillez entrer le pays"
        validationMsg.PlaseEnterConfirmPassAsNewPass = "Saisir le mot de Passe à confirmer, il doit être identique au nouveau mot de Passe";
        validationMsg.MobileMinlengthErrMsg = "Entrez s'il vous plait les 10 chiffres de votre numéro mobile"
    }

    //dynamic onchange login validation message
    validationMsg.UsernameErrEngMsg = "Please Enter Username";
    validationMsg.PasswordErrEngMsg = "Please Enter Password";
    validationMsg.OldPasswordErrEngMsg = "Please Enter Old Password";
    validationMsg.NewPasswordErrEngMsg = "Please Enter New Password";
    validationMsg.ConfirmPasswordErrEngMsg = "Please Enter Confirm Password";
    validationMsg.MatchConfirmPasswordErrEngMsg = "Enter Confirm Password Same as New Password";
    validationMsg.ForgotPasswordErrEngMsg = "Please Enter User Name";

    validationMsg.UsernameErrFreMsg = "Saisir le Utilisateur SVP";
    validationMsg.PasswordErrFreMsg = "Veuillez saisir le Mot de Passe, SVP";
    validationMsg.OldPasswordErrFreMsg = "Veuillez saisir l'ancien Mot de Passe";
    validationMsg.NewPasswordErrFreMsg = "Veuillez saisir le Nouveau Mot de Passe";
    validationMsg.ConfirmPasswordErrFreMsg = "Veuillez saisir le Confirmer mot de passe";
    validationMsg.MatchConfirmPasswordErrFreMsg = "Saisir le mot de Passe à confirmer, il doit être identique au nouveau mot de Passe";
    validationMsg.ForgotPasswordErrFreMsg = "Saisir le Utilisateur SVP";

    //  Form Validation
    window.addEventListener('load', function() {
        // Fetch all the forms we want to apply custom Bootstrap validation styles to
        var forms = document.getElementsByClassName('needs-validation');
        // Loop over them and prevent submission
        var validation = Array.prototype.filter.call(forms, function(form) {
            form.addEventListener('submit', function(event) {
                if (form.checkValidity() === false) {
                    event.preventDefault();
                    event.stopPropagation();
                }
                $(form).addClass('was-validated');
            }, false);
        });
    }, false);


    $('#input_new_password').bind("cut copy paste", function(e) {
        e.preventDefault();
    });
    $('#input_confirm_password').bind("cut copy paste", function(e) {
        e.preventDefault();
    });
    $('#input_old_password').bind("cut copy paste", function(e) {
        e.preventDefault();
    });

    $('#username').bind("cut copy paste", function(e) {
        e.preventDefault();
    });


    $("#input_new_password").on("keyup", function() {
        $("#validateStatus").val(1);
        var lowerCaseLetters = /[a-z]/g;
        if ($(this).val().match(lowerCaseLetters)) {
            $("#letter").removeClass("invalid");
            $("#letter").addClass("valid");
        } else {
            $("#letter").removeClass("valid");
            $("#letter").addClass("invalid");
            $("#validateStatus").val(0);
        }
        var upperCaseLetters = /[A-Z]/g;
        if ($(this).val().match(upperCaseLetters)) {
            $("#capital").removeClass("invalid");
            $("#capital").addClass("valid");
        } else {
            $("#capital").removeClass("valid");
            $("#capital").addClass("invalid");
            $("#validateStatus").val(0);
        }

        var specialCharas = /[-!$%^&*()_+|~=`{}[:;<>?,.@#\]]/g;
        if ($(this).val().match(specialCharas)) {
            $("#special-char").removeClass("invalid");
            $("#special-char").addClass("valid");
        } else {
            $("#special-char").removeClass("valid");
            $("#special-char").addClass("invalid");
            $("#validateStatus").val(0);
        }
        var numbers = /[0-9]/g;
        if ($(this).val().match(numbers)) {
            $("#number").removeClass("invalid");
            $("#number").addClass("valid");
        } else {
            $("#number").removeClass("valid");
            $("#number").addClass("invalid");
            $("#validateStatus").val(0);
        }
        if ($(this).val().length >= 6) {
            $("#length").removeClass("invalid");
            $("#length").addClass("valid");
        } else {
            $("#length").removeClass("valid");
            $("#length").addClass("invalid");
            $("#validateStatus").val(0);
        }
    });


    $("#userForm").validate({
        rules: {
            input_user_first_name: {
                required: true
            },
            input_user_last_name: {
                required: true
            },
            input_user_mobile: {
                required: true,
                minlength: 8,
                maxlength: 10
            },
            input_user_contact_number: {
                minlength: 8,
                maxlength: 10
            },
            input_user_email_id: {
                required: true,
                email: true
            },
            input_user_employee_id: {
                required: true
            },
            input_user_role: {
                required: true
            },
            input_user_department: {
                required: true
            },
            input_user_status: {
                required: true
            },
            input_user_language: {
                required: true
            },
            input_package_name:{
                required:true
            },
            input_Voume_amount:{
                required:true
            },
            input_transaction_fees:{
                required:true
            }

        },
        messages: {
            input_user_first_name: {
                required: validationMsg.UserFirstnamErreMsg
            },
            input_package_name:{
                required: validationMsg.PackagenameErrMsg            
        },
        input_Voume_amount:{
            required: validationMsg.volumeamtErrMsg            
        },
        input_transaction_fees:{
            required: validationMsg.TransactionFeeErrMsg            
        },
             input_user_last_name: {
                required: validationMsg.UserLastnameErrMsg
            },
            input_user_mobile: {
                required: validationMsg.MobileErrMsg,
                minlength: validationMsg.MobileMinlengthErrMsg
            },
            input_user_contact_number: {
                minlength: validationMsg.MobileMinlengthErrMsg
            },
            input_user_email_id: {
                required: validationMsg.EmailErrMsg,
                email: validationMsg.ValidEmailErrMsg
            },
            input_user_employee_id: {
                required: validationMsg.EmployeeIdErrMsg
            },
            input_user_role: {
                required: validationMsg.RoleErrMsg
            },
            input_user_department: {
                required: validationMsg.DepartmentErrMsg
            },
            input_user_status: {
                required: validationMsg.StatusErrMsg
            },
            input_user_language: {
                required: validationMsg.LanguageErrMsg
            }
        },
        errorPlacement: function(error, element) {
            if($(element).hasClass("select2")){
                $(element).parent("div").find('.select2-container').after(error);
            }else{
                 error.insertAfter(element);
            }
        },
        submitHandler: function(form) {
            $(".form-loader").show();
            // do other things for a valid form
            form.submit();
        }
    });

    $("#currencyForm").validate({
        rules: {
            input_code: {
                required: true
            },
            
            input_symbol:{
                required: true
            },
            input_country:{
                required: true
            } 
        },
        messages: {
            input_code: {
                required: validationMsg.CodeErrMsg,
                
            },
            input_symbol:{
                required: validationMsg.SymbolErrMsg
            },
             input_country:{
                required: validationMsg.CountryMsg
            }
             
        },
        errorPlacement: function(error, element) {
            if($(element).hasClass("select2")){
                $(element).parent("div").find('.select2-container').after(error);
                console.log("Validation initialized");
            }else{
                 error.insertAfter(element);
            }
        },
        submitHandler: function(form) {
            $(".form-loader").show();
            console.log("Validation initialized");
            // do other things for a valid form
            form.submit();
        }
    });

    $("#switchForm").validate({
        rules: {
            input_swtichstatus: {
                required: true
            }
        },
        messages: {
            input_swtichstatus: {
                required: validationMsg.StatusErrMsg
            }
        },
        errorPlacement: function(error, element) {
            if($(element).hasClass("select2")){
                $(element).parent("div").find('.select2-container').after(error);
            }else{
                 error.insertAfter(element);
            }
        },
        submitHandler: function(form) {
            $(".form-loader").show();
            // do other things for a valid form
            form.submit();
        }
    });

    $("#roleForm").validate({
        rules: {
            input_role_name: {
                required: true
            }
        },
        messages: {
            input_role_name: {
                required: validationMsg.RoleInputErrMsg
            }
        },
        errorPlacement: function(error, element) {
            if($(element).hasClass("select2")){
                $(element).parent("div").find('.select2-container').after(error);
            }else{
                 error.insertAfter(element);
            }
        },
        submitHandler: function(form) {
            $(".form-loader").show();
            // do other things for a valid form
            form.submit();
        }
    });

    $("#UpdateChannelForm").validate({
        rules: {
            input_gateway_status: {
                required: true
            },
            input_channel_status: {
                required: true
            },
            input_tokennumber: {
                required: true
            }
        },
        messages: {
            input_gateway_status: {
                required: validationMsg.StatusErrMsg
            },
            input_channel_status: {
                required: validationMsg.StatusErrMsg
            },
             input_tokennumber: {
                required: "Enter Token Number"
            }
        },
        errorPlacement: function(error, element) {
            if($(element).hasClass("select2")){
                $(element).parent("div").find('.select2-container').after(error);
            }else{
                 error.insertAfter(element);
            }
        },
        submitHandler: function(form) {
            $(".form-loader").show();
            // do other things for a valid form
            form.submit();
        }
    });

    $("#adminViewProf").validate({
        rules: {
            input_user_language: {
                required: true
            }
        },
        messages: {
            input_user_language: {
                required: validationMsg.LanguageErrMsg
            }
        },
        errorPlacement: function(error, element) {
            if($(element).hasClass("select2")){
                $(element).parent("div").find('.select2-container').after(error);
            }else{
                 error.insertAfter(element);
            }
        },
        submitHandler: function(form) {
            $(".form-loader").show();
            // do other things for a valid form
            form.submit();
        }
    });

    $("#login").validate({
        rules: {
            username: {
                required: true,
            },
            password: {
                required: true,
            }
        },
        messages: {
            username: {
                required: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.UsernameErrEngMsg;
                    } else {
                        return validationMsg.UsernameErrFreMsg;
                    }
                }
            },
            password: {
                required: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.PasswordErrEngMsg;
                    } else {
                        return validationMsg.PasswordErrFreMsg;
                    }
                },
            }
        },
        errorPlacement: function(error, element) {
            if($(element).hasClass("select2")){
                $(element).parent("div").find('.select2-container').after(error);
            }else{
                 error.insertAfter(element);
            }
        },
        submitHandler: function(form) {
            $(".form-loader").show();
            form.submit();
        }

    });

    $("#forgotPass").validate({
        rules: {
            username: {
                required: true,
            },
        },
        messages: {
            username: {
                required: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.ForgotPasswordErrEngMsg;
                    } else {
                        return validationMsg.ForgotPasswordErrFreMsg;
                    }
                }
            },
        },
        errorPlacement: function(error, element) {
            if($(element).hasClass("select2")){
                $(element).parent("div").find('.select2-container').after(error);
            }else{
                 error.insertAfter(element);
            }
        },
        submitHandler: function(form) {
            $(".form-loader").show();
            form.submit();
        }
    });

    $("#changePass").validate({
        rules: {
            input_old_password: {
                required: true,
            },
            input_new_password: {
                required: true,
            },
            input_confirm_password: {
                required: true,
                equalTo: "#input_new_password"
            }
        },
        messages: {
            input_old_password: {
                required: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.OldPasswordErrEngMsg;
                    } else {
                        return validationMsg.OldPasswordErrFreMsg;
                    }
                }
            },
            input_new_password: {
                required: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.NewPasswordErrEngMsg;
                    } else {
                        return validationMsg.NewPasswordErrFreMsg;
                    }
                }
            },
            input_confirm_password: {
                required: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.ConfirmPasswordErrEngMsg;
                    } else {
                        return validationMsg.ConfirmPasswordErrFreMsg;
                    }
                },
                equalTo: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.MatchConfirmPasswordErrEngMsg;
                    } else {
                        return validationMsg.MatchConfirmPasswordErrFreMsg;
                    }
                }
            }

        },
        errorPlacement: function(error, element) {
            if($(element).hasClass("select2")){
                $(element).parent("div").find('.select2-container').after(error);
            }else{
                 error.insertAfter(element);
            }
        },
        submitHandler: function(form) {
            if ($("#validateStatus").val() == 1) {
                $(".form-loader").show();
                form.submit();
            }
        }
    });

    $("#commonChangePass").validate({
        rules: {
            input_old_password: {
                required: true,
            },
            input_new_password: {
                required: true,
            },
            input_confirm_password: {
                required: true,
                equalTo: "#input_new_password"
            },

        },
        messages: {
            input_old_password: {
                required: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.OldPasswordErrEngMsg;
                    } else {
                        return validationMsg.OldPasswordErrFreMsg;
                    }
                }
            },
            input_new_password: {
                required: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.NewPasswordErrEngMsg;
                    } else {
                        return validationMsg.NewPasswordErrFreMsg;
                    }
                }
            },
            input_confirm_password: {
                required: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.ConfirmPasswordErrEngMsg;
                    } else {
                        return validationMsg.ConfirmPasswordErrFreMsg;
                    }
                },
                equalTo: function() {
                    var login_lang = $("#language").val();
                    if (login_lang === "English" || login_lang === "EN") {
                        return validationMsg.MatchConfirmPasswordErrEngMsg;
                    } else {
                        return validationMsg.MatchConfirmPasswordErrFreMsg;
                    }
                }
            }
        },
        errorPlacement: function(error, element) {
            if($(element).hasClass("select2")){
                $(element).parent("div").find('.select2-container').after(error);
            }else{
                 error.insertAfter(element);
            }
        },
        submitHandler: function(form) {
            if ($("#validateStatus").val() == 1) {
                $(".form-loader").show();
                form.submit();
            }
        }
    });

    $("#search_form").validate({
        rules: {
        },
        messages: {
        },
        errorPlacement: function(error, element) {
        },
        submitHandler: function(form) {
            $(".form-loader").show();
            form.submit();
        }
    });
    $("#log_form").validate({
        rules: {
        	input_cnpstxnnumnber: {
				required: true
			}
        },
        messages: {
        },
        errorPlacement: function(error, element) {
        	error.insertAfter(element);
        },
        submitHandler: function(form) {
            form.submit();
        }
    });


});