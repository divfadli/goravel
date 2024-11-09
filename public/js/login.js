$("#icon-click").click(function(e) {
    e.preventDefault();
    const passwordInput = $("#password");
    const icon = $("#icon");
    
    if (passwordInput.attr("type") === "password") {
        passwordInput.attr("type", "text");
        icon.removeClass("fa-eye-slash").addClass("fa-eye");
    } else {
        passwordInput.attr("type", "password");
        icon.removeClass("fa-eye").addClass("fa-eye-slash");
    }
});
$("#form-login").submit(function(e) {
    e.preventDefault();

    const email = $("#email").val();
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    
    if (!emailRegex.test(email)) {
        $("#error-message").text("Format email tidak valid");
        $("#email").focus();
        return false;
    }
    $.ajax({
        url: $(this).attr("action"),
        method: "POST",
        data: {
            email: $("#email").val(),
            password: $("#password").val()
        },
        success: function(data) {
            window.location.href = "/dashboard";
        },
        error:function(response){
            $("#password").val("");
            $("#password").focus();
            $("#error-message").text(response.responseJSON.message);
        }
    });
});