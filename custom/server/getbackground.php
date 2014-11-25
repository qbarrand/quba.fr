<?php

$data = file_get_contents("../../images/bg/backgrounds.json");

$images = json_decode($data);

$chosen = $images[rand(0, count($images) - 1)];

$return = json_encode($chosen);

header("Content-Type: application/json");
header("Content-Length: " . strlen($return));

echo $return;
