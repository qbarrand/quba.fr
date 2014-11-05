<?php

$to = "quentin@quba.fr";

$from_name = $_POST['name'];
$from_email = $_POST['email'];
$body = $_POST['body'];

$subject = "New message from " . $from_name . " (" . $from_email . ")";
$headers = "From: " . $from_email . "\r\n";

$response = array();

if(mail($to, $subject, $body, $headers)) {
	$response["status"] = true;
}
else {
	$response["status"] = false;
}

header("Content-type:application/json");
echo(json_encode($response));