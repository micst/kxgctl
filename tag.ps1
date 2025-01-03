$tag = (git tag --points-at HEAD)
$hasTag = ($null -ne $tag)
if (-not $hasTag) {
  $tag = "<empty>"
}
$c = "package kxg`n"
$c += "const (`n"
$c += "  Tag = ``$tag```n"
if ($hasTag) {
  $c += "  IsRelease = true`n"
}
else {
  $c += "  IsRelease = false`n"
}
$c += ")`n"
Set-Content -Path "kxg/tag.go" -Value $c
