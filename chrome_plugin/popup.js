function setShortAlias(tabs) {
    var short = $("#short_id").val();
    var secret = $("#secret").val();
    var url = tabs[0].url;
    $.ajax({
        url: 'https://yesteapea.com/bm/actions/add',
        type: 'POST',
        dataType: "html",
        data: {
            short: short,
            url: url,
            secret: secret
        },
        success: function(response) {
            console.log(response);
            if (response == "ok") {
                var s_url = `http://suram.in/r/${short}`;
                var message = `Row added at <a href='${s_url}'>${s_url}</a>`;
                $("#result").html(message);
            } else {
                $("#result").html(response);
            }
        },
        error: function(XMLHttpRequest, textStatus, exception) {
            $("#result").html(textStatus);
        },
        async: true
    });
}

$(function() {
    $('#button').click(
        function() {
            var query = {
                active: true,
                currentWindow: true
            };
            chrome.tabs.query(query, setShortAlias);
        });
});
