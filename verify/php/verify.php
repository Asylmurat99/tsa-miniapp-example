<?php
// #region check
// $fields: decoded values; $signField: 'hash' for the launch context,
// 'sign' for the getPhone envelope.
function checkSignature(array $fields, string $secret, string $signField): bool {
    $received = $fields[$signField] ?? '';
    unset($fields[$signField]);
    ksort($fields, SORT_STRING);
    $parts = [];
    foreach ($fields as $k => $v) {
        $parts[] = $k . '=' . $v;
    }
    $calc = hash_hmac('sha256', implode("\n", $parts), $secret);
    return hash_equals($calc, $received);
}
// #endregion check

$vector = json_decode(file_get_contents(__DIR__ . '/../vector.json'), true);
$failed = false;
foreach ($vector['cases'] as $c) {
    $ok = checkSignature($c['fields'], $vector['secret'], $c['sign_field']);
    echo $c['name'] . ': ' . ($ok ? 'ok' : 'FAIL') . PHP_EOL;
    if (!$ok) {
        $failed = true;
    }
}
exit($failed ? 1 : 0);
