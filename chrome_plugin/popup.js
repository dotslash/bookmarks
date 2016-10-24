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
        success: function (response) {
            console.log(response);
            if (response == "ok") {
                var s_url = `http://suram.in/r/${short}`;
                var message = `Added <a href='${s_url}'>${s_url}</a> as alias`;
                $("#result").html(message);
            } else {
                $("#result").html(response);
            }
        },
        error: function (XMLHttpRequest, textStatus, exception) {
            $("#result").html(textStatus);
        }
    });
}


function getExistingShortAliases(tabs) {
    var url = tabs[0].url;
    $.ajax({
        url: 'https://yesteapea.com/bm/actions/revlookup',
        type: 'POST',
        data: {
            long: url
        },
        success: function (response) {
            console.log(response);
            if (response.shorturls.length != 0) {
                var short_to_anchor = function (short) {
                    var s_url = `http://suram.in/r/${short}`;
                    return `<a href='${s_url}'>${s_url}</a>`;
                };
                var short_urls = response.shorturls.map(short_to_anchor).join("<br>");
                var result = 'short urls \n' + short_urls;
                $("#result").html(result);
            }
        }
    });
}

$(function () {
    var query = {
        active: true,
        currentWindow: true
    };

    $('#button').click(
        function () {
            chrome.tabs.query(query, setShortAlias);
        });
    chrome.tabs.query(query, getExistingShortAliases);
});
