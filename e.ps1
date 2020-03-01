$res = ./e.exe $args
if (!$res.Count) {
    exit($LASTEXITCODE)
}
$first, $rest = $res
if ($first.StartsWith("#env")) {
    foreach ($el in $rest) {
        $s = $el.Split("=")
        $key = $s[0]
        $value = $s[1]
        [Environment]::SetEnvironmentVariable($key, $value)
    }
    exit($LASTEXITCODE)
} else {
    $res | Write-Host
    exit($LASTEXITCODE)
}
