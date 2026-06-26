#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

rm -rf internal/generated

spec_path="swagger/run9_portal.yaml"
if [[ ! -f "$spec_path" ]]; then
  echo "missing swagger snapshot: $spec_path" >&2
  exit 1
fi

normalized_spec="$(mktemp)"
trap 'rm -f "$normalized_spec"' EXIT
go run ./cmd/normalizeswagger -in "$spec_path" -out "$normalized_spec"

operations=(
  updateAccount
  listAccountSSHKeys
  createAccountSSHKey
  deleteAccountSSHKey
  whoami
  updateOrg
  deleteOrg
  listOrgMembers
  updateOrgMember
  removeOrgMember
  listOrgInvitations
  createOrgInvitation
  revokeOrgInvitation
  listAPIKeys
  createAPIKey
  revokeAPIKey
  listOrgHosts
  listProjects
  createProject
  getProject
  updateProject
  deleteProject
  listProjectMembers
  updateProjectMember
  removeProjectMember
  listProjectSecrets
  createProjectSecret
  updateProjectSecret
  deleteProjectSecret
  createBox
  listBoxes
  getBox
  updateBox
  deleteBox
  stopBox
  listBoxSecrets
  createBoxSecret
  updateBoxSecret
  deleteBoxSecret
  importSnap
  listSnaps
  getSnap
  forkSnap
  deleteSnap
  getSnapTree
  execBox
  backgroundExecBox
  getExec
  killBackgroundExec
  listSharedSnaps
  getSharedSnap
  publishSharedSnap
  deleteSharedSnap
  deleteSharedSnapVersion
  consumeSharedSnapToBox
  consumeSharedSnapToSnap
)

args=(
  go run github.com/go-swagger/go-swagger/cmd/swagger@v0.34.1 generate client
  --skip-validation
  -f "$normalized_spec"
  -A run9_portal
  -t .
  --client-package internal/generated/client
  --model-package internal/generated/models
)

for operation in "${operations[@]}"; do
  args+=(-O "$operation")
done

"${args[@]}"
