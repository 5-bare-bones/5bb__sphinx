#!/usr/bin/env bash
#
# Asserts the core SPHINX invariant: a lower-tier binary ships strictly LESS
# code than a higher one — not merely fewer visible commands.
#
# It builds its OWN unstripped binaries into a temp dir, because symbol-table
# inspection (go tool nm) needs the symbols that release builds strip with
# -s -w. Run via `make verify-tiers`.
set -euo pipefail

cd "$(dirname "$0")/.."

out="$(mktemp -d)"
trap 'del -rf "$out"' EXIT

tiers=(apprentice adept scholar keeper master hero tutor)

# Cumulative build tags for a tier (kept bash 3.2 compatible — macOS ships it).
tags_for() {
	case "$1" in
		apprentice) echo "" ;;
		adept)      echo "adept" ;;
		scholar)    echo "adept scholar" ;;
		keeper)     echo "adept scholar keeper" ;;
		master)     echo "adept scholar keeper master" ;;
		hero)       echo "adept scholar keeper master hero" ;;
		tutor)      echo "adept scholar keeper master hero tutor" ;;
	esac
}

echo "== building unstripped tier binaries =="
for t in "${tiers[@]}"; do
	go build -tags "$(tags_for "$t")" -o "$out/sphinx_$t" .
done

fail=0
bin() { echo "$out/sphinx_$1"; }

topcount() {
	"$(bin "$1")" --help 2>&1 \
		| sed -n '/Available Commands/,/^Flags/p' \
		| grep -cE '^  [a-z0-9]'
}
symcount() {
	go tool nm "$(bin "$1")" 2>/dev/null | grep -c "$2" || true
}

echo "== top-level command counts (must strictly increase) =="
prev=-1
for t in "${tiers[@]}"; do
	n=$(topcount "$t")
	printf "  %-11s %s commands\n" "$t" "$n"
	if [ "$n" -le "$prev" ]; then
		echo "FAIL: $t ($n) did not exceed previous tier ($prev)"
		fail=1
	fi
	prev=$n
done

echo "== code elision (command-package symbol counts) =="
assert_absent() {
	c=$(symcount "$1" "$2")
	if [ "$c" -ne 0 ]; then echo "FAIL: $2 code present in $1 ($c symbols)"; fail=1; fi
}
assert_present() {
	c=$(symcount "$1" "$2")
	if [ "$c" -eq 0 ]; then echo "FAIL: $2 code missing from $1"; fail=1; fi
}
assert_absent  apprentice commands/card
assert_absent  adept      commands/card
assert_present scholar    commands/card
assert_absent  keeper     commands/config
assert_present master     commands/config
assert_absent  apprentice commands/vault/backup
assert_present keeper     commands/vault/backup
echo "  ok"

echo "== per-subcommand gating: file move only at scholar+ =="
filesub() { "$(bin "$1")" file --help 2>&1 | sed -n '/Available Commands/,/^Flags/p' | grep -cE "^  $2 "; }
if [ "$(filesub adept move)" -ne 0 ]; then echo "FAIL: 'file move' present at adept"; fail=1; fi
if [ "$(filesub scholar move)" -eq 0 ]; then echo "FAIL: 'file move' missing at scholar"; fail=1; fi
echo "  ok"

if [ "$fail" -ne 0 ]; then
	echo "TIER VERIFICATION FAILED"
	exit 1
fi
echo "ALL TIER INVARIANTS HOLD"
