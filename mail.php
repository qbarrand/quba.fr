<?php

$response = array();

$to = "quentin@quba.fr";

$from_name = $_POST['name'];
$from_email = $_POST['email'];
$body = $_POST['body'];

$subject = "New message from " . $from_name . " (" . $from_email . ")";
$headers = "From: " . $from_email . "\r\n";

$response["status"] = mail($to, $subject, $body, $headers);

// Send the response to the client
echo(json_encode($response));

fastcgi_finish_request();

$response["client_data"] = $_POST;

// Also save it locally in case there was an issue with the MTA
$mails = fopen('mails.json', "a");

fwrite($mails, json_encode($response));
fwrite($mails, "\n");

fclose($mails);
