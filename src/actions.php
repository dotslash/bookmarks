<?php
    require_once("libs/EditableGrid.php");
    require_once("./db_helpers.php");

    function checkSecret(){
      $secret = $_REQUEST['secret'];
      global $config;
      $base = "http://$_SERVER[HTTP_HOST]/";
      if ($secret != $config->bm_secret) {
        echo "Error: Wrong Secret - $secret";
        return false;
      }
      return true;

    }
    function getBmsAsJson(){
      $bms = getExistingBM();

      //print_r(json_encode($resp,JSON_PRETTY_PRINT));
      $grid = new EditableGrid();
      $grid->addColumn("bookmark", "BOOKMARK", "string");
      $grid->addColumn("original_url", "ORIGINAL", "website", false);
      $grid->addColumn("bm_url", "BM URL", "website", false);
      $grid->addColumn('action', 'Actions', 'html', NULL, false, 'id');
      $grid->renderJSON($bms);
    }
    function add_bm(){
        $long = $_REQUEST['url'];
        $short = $_REQUEST['short'];

        $sec = checkSecret();
        if(!$sec) {
          return;
        }
        else if(false && !preg_match("/^[a-zA-Z]+[:\/\/]+[A-Za-z0-9\-_\,]+\\.+[A-Za-z0-9\.\/%&=\?\-_]+$/i", $long)) {
          echo "Error: invalid URL";
        }
        else if (bmExists($short)) {
          echo "Error: bookmark alredy taken";
        }
        else {
          addBM($long, $short);
          echo "ok";
          //echo "bookmarked url is ${base}red/$short";
        }
    }

    function delete(){
        $sec = checkSecret();
        if(!$sec) {
          return;
        }
        echo deleteBmById($_POST['id']);
    }

    function update_bm(){
        $id = $_POST['id'];
        $nval = $_POST['newvalue'];
        $col = $_POST['colname'];

        $sec = checkSecret();
        if(!$sec) {
          return;
        }
        echo updateBm($id, $col, $nval);
    }

    $action = $_GET['action'];

    if ($action == 'update') {
        update_bm();
    } else if ($action == 'delete') {
        delete();
    } else if ($action == 'view') {
        getBmsAsJson();
    } else if ($action == 'add' ) {
        add_bm();
    } else {
      echo "error";
    }
?>
