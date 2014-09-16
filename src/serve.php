<?php
    require("./db_helpers.php");
    $short = $_REQUEST['short'];
    $db = getConn();
    $sql = "SELECT original_url FROM bookmarks WHERE bookmark=\"$short\" LIMIT 1";
    $query = mysqli_query($db,$sql);
    error_log($sql);
    $html = "all good";
    $row = mysqli_fetch_array($query);
    if(!empty($row)) {
      Header("HTTP/1.1 301 Moved Permanently");
      header("Location: ".$row[0]."");
    } else {
      $html = "Error: cannot find short URL";
    }
    mysqli_close($db);
?>

<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>Redirect Page</title>
  </head>
  <body>
    <?= $html ?>
    <br /><br />
    <span class="back"><a href="./">X</a></span>
  </body>
</html>
