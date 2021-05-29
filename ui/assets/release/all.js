$(document).ready(function () {

    $('#type').select2();
    $('.reviewers').hide()
    $('#projects').change(function(){
        projectsVal = $(this).val()
        alert(projectsVal)
        getReviwers(projectsVal)
    });


    $('#projects').select2({
        placeholder: 'Select Projects',
        allowClear: true
    });

    var getReviwers = function (projectsVal) {
        $.ajax({
            url: "/release/getReviewers",
            type:"get",
            data: {'ids': projectsVal},
            context: document.body
        })
        .done(function (res) {
            $('#reviewers').append(projectsVal);
            $('.reviewers').show()
        })
        .fail(function (response) {
            $('#reviewers').val("failed");
            $('.reviewers').show()
        });
    }
})