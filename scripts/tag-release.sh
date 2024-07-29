#!/usr/bin/env bash
set -euo pipefail

################################################################################
# This script portably manages tagging of the HEAD commit based on the contents
# of the VERSION file. It performs some heuristic checks to make sure tagging
# will behave as expected.
################################################################################

root="$(git rev-parse --show-toplevel)"

current_git_branch="$(git rev-parse --abbrev-ref HEAD)"
latest_git_tag="$(git tag --list | tail -n1)"
current_listed_version="$(grep -v '#' "${root:-}"/VERSION)"

printf 'Current git branch: %s\n' "${current_git_branch:-}"
printf 'Latest git tag: %s\n' "${latest_git_tag:-}"
printf 'Current version listed in VERSION file: %s\n' "${current_listed_version:-}"

failures=()

if [[ "${current_git_branch}" == 'main' ]] ; then
  # Fail if we forgot to bump VERSION
  if [[ "${latest_git_tag}" == "${current_listed_version}" ]] ; then
    printf 'ERROR: Identifier in VERSION still matches what is tagged on the main branch -- did you forget to update?\n' > /dev/stderr
    failures+=('forgot-to-bump-VERSION')
  fi

  # Fail if we forgot to bump versions across files
  old_git_status="$(git status | grep -i -E 'modified')"
  make -s bump-versions old_version="${current_listed_version}"
  new_git_status="$(git status | grep -i -E 'modified')"
  if [[ "$(diff <(echo "${old_git_status}") <(echo "${new_git_status}") | wc -l)" -gt 0 ]] ; then
    printf 'ERROR: Files modified by version-bump check -- did you forget to update versions across the repo to match VERSION?\n' > /dev/stderr
    failures+=('forgot-to-bump-other-versions')
  fi

  if [[ "${#failures[@]}" -gt 0 ]] ; then
    exit 1
  else
    printf 'All checks passed, tagging & pushing new version: %s --> %s\n' "${latest_git_tag}" "${current_listed_version}"
    git tag --force "v${current_listed_version}"
    git push --tags
  fi

else
  printf 'Not on main branch, nothing to do\n'
  exit 0
fi
