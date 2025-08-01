#!/bin/bash
# Generates a changelog using Standard and Gitmoji commit conventions

set -o errexit
set -o nounset
set -o pipefail

if [[ $# -ne 3 ]]; then
  echo "Usage: $0 <module> <tag> <output_file>"
  exit 1
fi

MODULE="$1"
TAG="$2"
OUTPUT="$3"

if ! git rev-parse --git-dir > /dev/null 2>&1; then
  echo "Error: Not in a git repository"
  exit 1
fi

get_previous_tag() {
  git tag -l | sort -V | grep -B1 "^${1}$" | head -1 | grep -v "^${1}$" || echo ""
}

generate_commit_list() {
  if [[ -z "$1" ]]; then
    git log --oneline --no-merges "$2" || echo ""
  else
    git log --oneline --no-merges "${1}..${2}" || echo ""
  fi
}

categorize_commits() {
  local commits="$1"
  breakings=(); features=(); fixes=(); docs=(); perf=(); chore=(); others=()

  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    msg="${line#* }"

    case "$msg" in
      💥*|:boom:*|*BREAKING*CHANGE*) breakings+=("$msg") ;;
      ✨*|:sparkles:*) features+=("$msg") ;;
      🐛*|🐞*|:bug:*) fixes+=("$msg") ;;
      🔧*|:wrench:*) chore+=("$msg") ;;
      🚀*|:rocket:*) perf+=("$msg") ;;
      📦*|:package:*) perf+=("$msg") ;;
      ♻️*|:recycle:*) chore+=("$msg") ;;
      ✅*|:white_check_mark:*) chore+=("$msg") ;;
      📝*|:memo:*) docs+=("$msg") ;;
      🔒*|:lock:*) perf+=("$msg") ;;
      🛠*|:hammer:*) chore+=("$msg") ;;
      *) # fallback to conventional keywords
        if [[ "$msg" =~ ^(feat)(\(.+\))?: ]]; then features+=("$msg")
        elif [[ "$msg" =~ ^(fix)(\(.+\))?: ]]; then fixes+=("$msg")
        elif [[ "$msg" =~ ^(docs)(\(.+\))?: ]]; then docs+=("$msg")
        elif [[ "$msg" =~ ^(perf|refactor)(\(.+\))?: ]]; then perf+=("$msg")
        elif [[ "$msg" =~ ^(chore|style|test)(\(.+\))?: ]]; then chore+=("$msg")
        else others+=("$msg")
        fi
        ;;
    esac
  done <<< "$commits"

  {
    echo "# Release $TAG"
    echo
    echo "## What's Changed"
    echo

    [[ ${#breakings[@]} -gt 0 ]] && {
      echo "### 💥 Breaking Changes"
      for c in "${breakings[@]}"; do echo "- $c"; done
      echo
    }

    [[ ${#features[@]} -gt 0 ]] && {
      echo "### ✨ New Features"
      for c in "${features[@]}"; do echo "- $c"; done
      echo
    }

    [[ ${#fixes[@]} -gt 0 ]] && {
      echo "### 🐛 Bug Fixes"
      for c in "${fixes[@]}"; do echo "- $c"; done
      echo
    }

    [[ ${#docs[@]} -gt 0 ]] && {
      echo "### 📝 Documentation"
      for c in "${docs[@]}"; do echo "- $c"; done
      echo
    }

    [[ ${#perf[@]} -gt 0 ]] && {
      echo "### ⚡ Performance / Deploy"
      for c in "${perf[@]}"; do echo "- $c"; done
      echo
    }

    [[ ${#chore[@]} -gt 0 ]] && {
      echo "### 🔧 Chores / Config / Tests"
      for c in "${chore[@]}"; do echo "- $c"; done
      echo
    }

    [[ ${#others[@]} -gt 0 ]] && {
      echo "### 🔍 Other Changes"
      for c in "${others[@]}"; do echo "- $c"; done
      echo
    }

    contributors=$(generate_contributors "$previous_tag" "$TAG")
    if [[ -n "$contributors" ]]; then
      echo "### 👥 Contributors"
      echo "$contributors"
      echo
    fi

    echo "**Full Changelog**: https://github.com/$MODULE/compare/${previous_tag:-initial}...${TAG}"
  } > "$OUTPUT"
}

generate_contributors() {
  range=${1:+$1..}${2}
  git log --format='%an <%ae>' ${range} 2>/dev/null |
    sort -u | while read -r line; do
      name="${line% <*}"; email="${line#*<}"; email="${email%>}"
      echo "- $name"
    done
}

echo "Generating changelog for $MODULE $TAG"
previous_tag=$(get_previous_tag "$TAG")
echo "Previous tag: ${previous_tag:-none}"
commits=$(generate_commit_list "$previous_tag" "$TAG")
categorize_commits "$commits"
echo "Wrote changelog to $OUTPUT"
