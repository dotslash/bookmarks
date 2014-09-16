<?php
require_once("./project_config.php");


function getConn(){
  global $config;
  $db = mysqli_connect($config->host,
          $config->username,
          $config->password,
          $config->database,
          $config->port);
  connGood();
  return $db;
}

function connGood(){
  if (mysqli_connect_errno()) {
    error_log("Failed to connect to MySQL: " . mysqli_connect_error());
  }
  return 0;
}

function getExistingBM(){
  $base = "http://$_SERVER[HTTP_HOST]/red/";
  $db = getConn();
  $query = "SELECT *, 1 as action, CONCAT(\"$base\",bookmark) as bm_url from bookmarks";
  $res=mysqli_query($db, $query);
  mysqli_close($db);
  return $res;
}

function bmExists($bm){
  $result = getExistingBM();
  while($row = mysqli_fetch_array($result)) {
    if ($row['bookmark'] === $bm) {
        return true;

    }
  }
  return false;
}

function addBM($long, $short){
  $db = getConn();
  $query = "INSERT INTO bookmarks (bookmark, original_url) VALUES (\"$short\", \"$long\");";
  $status = mysqli_query($db, $query);
  if (!$status) {
    error_log('error adding row: '.  mysqli_error($db));
  }
  mysqli_close($db);
}


function deleteBmById($id){
  $db = getConn();
  $ret = "ok";
  $query = "DELETE FROM bookmarks WHERE id = $id";
  error_log($query);
  $status = mysqli_query($db, $query);
  if (!$status) {
    error_log('error deleting row: '.  mysqli_error($db));
    $ret = "error";
  }
  mysqli_close($db);
  return $ret;
}

function updateBm($id, $col, $nVal){
  $db = getConn();
  $ret = "ok";
  $query = "UPDATE bookmarks SET $col = \"$nVal\" WHERE id = $id";
  $status = mysqli_query($db, $query);
  error_log($query);
  if (!$status) {
    error_log('error updating row: '.  mysqli_error($db));
    $ret = "error";
  }
  mysqli_close($db);
  return $ret;
}
?>
