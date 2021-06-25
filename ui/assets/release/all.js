var $releasesTable;
$(document).ready(function () {
    $('#releases_name').select2();
    window.$reviewerList = $('#reviewers');
    $('#type').select2();
    $('.reviewers').hide()
    $('.reviewers_list').hide()


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
            $('#reviewers_values').val(res.data);
            $('.reviewers_list').show()
        })
        .fail(function (response) {
            $('#reviewers_values').val(response.message);
            $('.reviewers_list').show()
        });
    }

    $releasesTable = $('#releases_table');
    if ($releasesTable.length) {
        var additionalOptions = {
            order: [],
            language: {
                searchPlaceholder: "Search Releases"
            },
            bInfo:false,
            pagingType: "simple",
            responsive: true,
            targets: 'no-sort',
            "bSort": false,
        };
        var releaseDt = initDatatable($releasesTable, [
            {data: 'ID', name: 'ID', 'visible': false, searchable: false},
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
                    return "<a href='/release/show/" + data + "'><i class='fa fa-eye'></i></a>";
                }
            },
        ], additionalOptions);
        $releasesTable.data('dt', releaseDt);
    }

    $projectsTable = $('#projects_table');
    if ($projectsTable.length) {
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
        var projectDt = initDatatable($projectsTable, [
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
        $projectsTable.data('dt', projectDt);
    }

    $reviewersTable = $('#reviewers_table');
    if ($reviewersTable.length) {
        var additionalOptions = {
            // order: [[0, "desc"]],
            language: {
                searchPlaceholder: "Search Reviewer"
            },
            // bInfo:false,
            responsive: true,
            columnDefs: [
                { className: 'text-center', targets: [0,1,2,3] },
            ],
        };
        var reviewersDt = initDatatable($reviewersTable, [
            {data: 'name', name: 'Reviewer name', searchable: true,'orderable': false},
            {data: 'Email', name: 'Email', searchable: true,'orderable': false},
            {data: 'username', name: 'User name', searchable: true,'orderable': false},
            {
                "mData": "ID",
                'orderable': false,
                "mRender": function (data, type, row) {
                    return "<a class='btn btn-sm' href='/reviewers/show/" + data + "'><i class='fa fa-edit text-primary'></i></a><button data-id="+row.ID+" class='reviewer_removal btn btn-sm' data-target='#reviewer_removal_modal'><i class='fa fa-trash text-danger'></i></a>";
                }
            },
        ], additionalOptions);
        $reviewersTable.on('click', '.reviewer_removal', function () {
            $('#reviewer_removal_modal').modal('show');
            $('#reviewer_removal_yes').attr("data-id",$(this).data('id'));
        });
        $reviewersTable.data('dt', reviewersDt);
    }

    $dlsTable = $('#dls_table');
    if ($dlsTable.length) {
        var additionalOptions = {
            order: [[0, "desc"]],
            language: {
                searchPlaceholder: "Search DLS"
            },
            bInfo:false,
            responsive: true,
            columnDefs: [
                { className: 'text-center', targets: [0,2] },
            ],
        };
        var dlsDt = initDatatable($dlsTable, [
            {data: 'ID', name: 'ID', searchable: true,'orderable': false},
            {data: 'Email', name: 'Email', searchable: true,'orderable': false},
            {data: 'DlType', name: 'DlType', searchable: true,'orderable': false},
            {
                "mData": "ID",
                'orderable': false,
                "mRender": function (data, type, row) {
                    return "<a class='btn btn-sm' href='/dls/show/" + data + "'><i class='fa fa-edit text-primary'></i></a><button data-id="+row.ID+" class='reviewer_removal btn btn-sm' data-target='#reviewer_removal_modal'><i class='fa fa-trash text-danger'></i></a>";
                }
            },
        ], additionalOptions);
        $dlsTable.on('click', '.reviewer_removal', function () {
            $('#reviewer_removal_modal').modal('show');
            $('#reviewer_removal_yes').attr("data-id",$(this).data('id'));
        });
        $dlsTable.data('dt', dlsDt);
    }
})

$('#reviewer_removal_modal_close').on('click', function () {
    $('#reviewer_removal_modal').modal("hide")
});

$('#reviewer_removal_yes').on('click', function () {
    $('#reviewer_removal_modal').modal("hide")
    let url = $(this).data('url')
    window.location.href = "/"+url+"/delete/" + $(this).data('id')
    
});

$('#reviewers').select2({
    placeholder: 'Select Reviewers',
    allowClear: true
});

$('#reviewers').on('change', function () {
    var str = [];
    $("#reviewers option:selected").each(function () {
        str.push(this.value);
    });
});

let selectedReviewers = $("#selected_reviewers")
if(selectedReviewers.length > 0) {
    let selectedReviewersListString = selectedReviewers.data("selected_reviewers")
    $('#reviewers').val(selectedReviewersListString.split(',')).trigger('change');
}

var initDatatable = function ($table, $columns, additionalOptions) {
    var options = {
        // dom: 'lfrtip',
        processing: true,
        serverSide: true,
        // autoWidth: false,
        "ordering": false,
        ajax: {
            url: $table.data('get_action'),
            type: 'post',
            error: function (xhr, err) {
                if (err === 'parsererror')
                    location.reload();
            }
        },
        "pageLength": $('.admin.dashboard').length ? 5 : 10,
        // "pagingType": "full_numbers",
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

$('#releases_name').on('change', function () {
    var str = $("#releases_name option:selected").text();
    getJiraTicketsByLabel(str)
    alert(str)
});

var getJiraTicketsByLabel = function (releaseName) {
    $.ajax({
        url: "/release/getTickets?release="+releaseName,
        type:"get",
        context: document.body
    })
        .done(function (res) {
            $(function() {
                $.each(res.data, function(i, item) {
                    var $tr = $('<tr>').append(
                        $('<td>').text(item.Id),
                        $('<td>').text(item.Summary),
                        $('<td>').text(item.CreationDate),
                        $('<td>').text(item.CreationTime),
                        $('<td>').text(item.Type),
                        $('<td>').text(item.Project),
                        $('<td>').text(item.Priority),
                        $('<td>').text(item.Status),
                    ).appendTo('.release_tickets');
                });
            });
        })
        .fail(function (response) {
           alert('Failed to load data')
        });
}

