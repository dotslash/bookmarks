/**
 *  highlightRow and highlight are used to show a visual feedback. If the row has been successfully modified, it will be highlighted in green. Otherwise, in red
 */
function highlightRow(rowId, bgColor, after) {
    var rowSelector = $("#" + rowId);
    rowSelector.css("background-color", bgColor);
    rowSelector.fadeTo("normal", 0.5, function() {
        rowSelector.fadeTo("fast", 1, function() {
            rowSelector.css("background-color", '');
        });
    });
}

function highlight(div_id, style) {
    highlightRow(div_id, style == "error" ? "#e5afaf" : style == "warning" ? "#ffcc00" : "#8dc70a");
}

/**
   updateCellValue calls the PHP script that will update the database.
 */
function updateCellValue(editableGrid, rowIndex, columnIndex, oldValue, newValue, row, onResponse) {
    console.log(oldValue);
    console.log(newValue);
    console.log(editableGrid.getColumnName(columnIndex));
    console.log(editableGrid.getColumnType(columnIndex));
    $.ajax({
        url: 'actions/update',
        type: 'POST',
        dataType: "html",
        data: {
            tablename: editableGrid.name,
            id: editableGrid.getRowId(rowIndex),
            newvalue: editableGrid.getColumnType(columnIndex) == "boolean" ? (newValue ? 1 : 0) : newValue,
            oldvalue: oldValue,
            colname: editableGrid.getColumnName(columnIndex),
            coltype: editableGrid.getColumnType(columnIndex),
            secret : $('#secret').val()
        },
        success: function(response) {
            // reset old value if failed then highlight row
            var success = onResponse ? onResponse(response) : (response == "ok" || !isNaN(parseInt(response))); // by default, a sucessfull reponse can be "ok" or a database id
            if (!success) editableGrid.setValueAt(rowIndex, columnIndex, oldValue);
            highlight(row.id, success ? "ok" : "error");
        },
        error: function(XMLHttpRequest, textStatus, exception) {
            alert("Ajax failure\n" + errortext);
        },
        async: true
    });

}



function DatabaseGrid() {
    this.editableGrid = new EditableGrid("demo", {
        enableSort: true,
        // define the number of row visible by page
        pageSize: 10,
        // Once the table is displayed, we update the paginator state
        tableRendered: function() {
            updatePaginator(this);
        },
        tableLoaded: function() {
            datagrid.initializeGrid(this);
        },
        modelChanged: function(rowIndex, columnIndex, oldValue, newValue, row) {
                updateCellValue(this, rowIndex, columnIndex, oldValue, newValue, row);
            }
    });
    this.fetchGrid();
}

DatabaseGrid.prototype.fetchGrid = function() {
    // call a PHP script to get the data
    //this.editableGrid.loadJSON("loaddata.php");
    this.editableGrid.loadJSON("actions/view");
};

DatabaseGrid.prototype.initializeGrid = function(grid) {

    var self = this;

    // render for the action column
    grid.setCellRenderer("action", new CellRenderer({
        render: function(cell, _id) {
            var rowId = grid.getRowId(cell.rowIndex).trim();
            console.log(rowId)
            cell.innerHTML += "<i onclick=\"datagrid.deleteRow('" + rowId + "');\" class='fa fa-trash-o' ></i>";
        }
    }));


    grid.renderGrid("tablecontent", "testgrid");
};

DatabaseGrid.prototype.deleteRow = function(id) {

    var self = this;
    var check = confirm('Sure?');
    console.log(id);
    if (check) {

        $.ajax({
            url: 'actions/delete',
            type: 'POST',
            dataType: "html",
            data: {
                tablename: self.editableGrid.name,
                id: id,
                secret : $('#secret').val()
            },
            success: function(response) {
                if (response == "ok")
                    self.editableGrid.removeRow(id);
            },
            error: function(XMLHttpRequest, textStatus, exception) {
                alert("Ajax failure\n" + errortext);
            },
            async: true
        });


    }

};



DatabaseGrid.prototype.addRow = function(id) {

    var self = this;
    $.ajax({
        url: 'actions/add',
        type: 'POST',
        dataType: "html",
        data: {
            tablename: self.editableGrid.name,
            short: $("#short").val(),
            url: $("#long").val(),
            secret : $('#secret').val()
        },
        success: function(response) {
            if (response == "ok") {

                // hide form
                showAddForm();
                $("#short").val('');
                $("#long").val('');

                alert("Row added : reload model");
                self.fetchGrid();
            } else
                alert(response);
        },
        error: function(XMLHttpRequest, textStatus, exception) {
            alert("Ajax failure\n" + errortext);
        },
        async: true
    });
};


function updatePaginator(grid, divId) {
    divId = divId || "paginator";
    var paginator = $("#" + divId).empty();
    var nbPages = grid.getPageCount();

    // get interval
    var interval = grid.getSlidingPageInterval(20);
    if (interval == null) return;

    // get pages in interval (with links except for the current page)
    var pages = grid.getPagesInInterval(interval, function(pageIndex, isCurrent) {
        if (isCurrent) return "<span id='currentpageindex'>" + (pageIndex + 1) + "</span>";
        return $("<a>").css("cursor", "pointer").html(pageIndex + 1).click(function(event) {
            grid.setPageIndex(parseInt($(this).html()) - 1);
        });
    });

    // "first" link
    var link = $("<a class='nobg'>").html("<i class='fa fa-fast-backward'></i>");
    if (!grid.canGoBack()) link.css({
        opacity: 0.4,
        filter: "alpha(opacity=40)"
    });
    else link.css("cursor", "pointer").click(function(event) {
        grid.firstPage();
    });
    paginator.append(link);

    // "prev" link
    link = $("<a class='nobg'>").html("<i class='fa fa-backward'></i>");
    if (!grid.canGoBack()) link.css({
        opacity: 0.4,
        filter: "alpha(opacity=40)"
    });
    else link.css("cursor", "pointer").click(function(event) {
        grid.prevPage();
    });
    paginator.append(link);

    // pages
    for (p = 0; p < pages.length; p++) paginator.append(pages[p]).append(" ");

    // "next" link
    link = $("<a class='nobg'>").html("<i class='fa fa-forward'>");
    if (!grid.canGoForward()) link.css({
        opacity: 0.4,
        filter: "alpha(opacity=40)"
    });
    else link.css("cursor", "pointer").click(function(event) {
        grid.nextPage();
    });
    paginator.append(link);

    // "last" link
    link = $("<a class='nobg'>").html("<i class='fa fa-fast-forward'>");
    if (!grid.canGoForward()) link.css({
        opacity: 0.4,
        filter: "alpha(opacity=40)"
    });
    else link.css("cursor", "pointer").click(function(event) {
        grid.lastPage();
    });
    paginator.append(link);
};


function showAddForm() {
    if ($("#addform").is(':visible'))
        $("#addform").hide();
    else
        $("#addform").show();
}

function showSecret() {
    if ($("#inpSecret").is(':visible'))
        $("#inpSecret").hide();
    else
        $("#inpSecret").show();
}
