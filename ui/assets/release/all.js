$(document).ready(function () {
    window.$reviewerList = $('#reviewers');
    $('#type').select2();
    $('.reviewers').hide()


    $('#projects').select2({
        placeholder: 'Select Projects',
        allowClear: true
    });

    $('#projects').on('change', function () {
        var str = [];
        $("#projects option:selected").each(function () {
            str.push(this.value);
        });
        getReviwers(str);
    });




    var getReviwers = function (projectsVal) {
        $.ajax({
            url: "/release/getReviewers?ids="+projectsVal,
            type:"get",
            context: document.body
        })
        .done(function (res) {
            alert(res.data)
            $('#reviewers').val(res.data);
            $('.reviewers').show()
        })
        .fail(function (response) {
            $('#reviewers').val(response.message);
            $('.reviewers').show()
        });
    }
})