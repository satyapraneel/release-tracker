$(document).ready(function () {
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
});