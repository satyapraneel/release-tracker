var $releasesTable;
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

    $releasesTable = $('#releases_table');
    if ($releasesTable.length) {
        var additionalOptions = {
            order: [[0, "desc"]],
            language: {
                searchPlaceholder: "Search Releases"
            },
            bInfo:false,
            // pagingType: "simple",
            responsive: true,
        };
        var releaseDt = initDatatable($releasesTable, [
            {data: 'ID', name: 'ID', 'visible': true, searchable: false},
            {data: 'Name', name: 'Name', searchable: true},
            {data: 'type', name: 'type'},
            {data: 'target_date', name: 'target_date',
                "render": function (data) {
                    var date = new Date(data);
                    var month = date.getMonth() + 1;
                    var day = date.getDate();
                    return date.getFullYear() + "/"
                        + (month.length > 1 ? month : "0" + month) + "/"
                        + ("0" + day).slice(-2);
                    //return date format like 2000/01/01;
                },
            },
            {data: 'owner', name: 'owner'},
            {
                "mData": "ID",
                "mRender": function (data, type, row) {
                    return "<a href='/release/show/" + data + "'><i class='far fa-eye'></i></a>";
                }
            },
        ], additionalOptions);
        $releasesTable.data('dt', releaseDt);
    }

})

var initDatatable = function ($table, $columns, additionalOptions) {
    var options = {
        dom: 'lfrtip',
        processing: true,
        serverSide: true,
        autoWidth: false,
        ajax: {
            url: $table.data('get_action'),
            type: 'post',
            error: function (xhr, err) {
                if (err === 'parsererror')
                    location.reload();
            }
        },
        "pageLength": $('.admin.dashboard').length ? 5 : 10,
        "pagingType": "full_numbers",
        columns: $columns,
        drawCallback: function (settings) {
            var $dtContainer = $($(this).data('dt').table().container());

            var message = $dtContainer.data('dt_message');
            if (message) {
                $dtContainer.data('dt_message', '');
                swal(message.header, message.text, message.type);
            }
        }
    };

    //Es6 merging of options
    var options = Object.assign(options, additionalOptions);

    var $dt = $table.DataTable(options);

    $table.data('dt', $dt);
    return $dt;
}