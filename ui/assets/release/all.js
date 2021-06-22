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
            {data: 'Name', name: 'Name', searchable: true,'orderable': false},
            {data: 'type', name: 'type','orderable': false},
            {data: 'target_date', name: 'target_date','orderable': false,
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
            {data: 'owner', name: 'owner', 'orderable': false},
            {
                "mData": "ID",
                'orderable': false,
                "mRender": function (data, type, row) {
                    return "<a href='/release/show/" + data + "'><i class='fa fa-eye fa-2x'></i></a>";
                }
            },
        ], additionalOptions);
        $releasesTable.data('dt', releaseDt);
    }

    $releasesTable = $('#projects_table');
    if ($releasesTable.length) {
        var additionalOptions = {
            order: [[0, "desc"]],
            language: {
                searchPlaceholder: "Search Projects"
            },
            bInfo:false,
            responsive: true,
            columnDefs: [
                { className: 'text-center', targets: [0,1,2,3,4,5,6,7] },
            ],
        };
        var releaseDt = initDatatable($releasesTable, [
            {data: 'Name', name: 'Project name', searchable: true,'orderable': false},
            {data: 'repo_name', name: 'Repo name', searchable: false,'orderable': false},
            {data: 'beta_release_date', name: 'Beta Release', searchable: false,'orderable': false},
            {data: 'regression_signor_date', name: 'Regression Signor', searchable: false,'orderable': false},
            {data: 'code_freeze_date', name: 'Code Freeze', searchable: false,'orderable': false},
            {data: 'dev_completion_date', name: 'Dev Completion','orderable': false},
            {data: 'Status', name: 'Status', 'visible': true, searchable: false, "render": function (data) {
                return data == 1 ? "Active" : "Inactive"
            }},
            {
                "mData": "ID",
                'orderable': false,
                "mRender": function (data, type, row) {
                    return "<a href='/projects/show/" + data + "'><i class='fa fa-edit'></i></a>";
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