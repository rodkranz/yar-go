<?php

$client = new Yar_Client('tcp://192.168.30.4:12345');
$args = ['Urls' => ["http://www.google.com/", "http://www.facebook.com/", "http://www.terra.com.br/"]];
$replay = $client->__call("Fetch.MultipleRequest", $args);

print '<pre>';

if (!is_null($replay['Err'])) {
    printf("Replay: %v \n", $replay['Err']);
    exit(1);
}

foreach ($replay['Responses'] as $value) {
    printf(
        "[%d] elapsed time for request [%s] with [%d]\n",
        $value['Time'],
        $value['Url'],
        $value['Bytes']);
}

printf("[%d] elapsed time.\n", $replay['Time']);
print '</pre>';
