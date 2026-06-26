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
  deleteProjectSecret
  createBox
  listBoxes
  getBox
  deleteBox
  stopBox
  listBoxSecrets
  createBoxSecret
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
  -f "$spec_path"
  -A run9_portal
  -t .
  --client-package internal/generated/client
  --model-package internal/generated/models
)

for operation in "${operations[@]}"; do
  args+=(-O "$operation")
done

"${args[@]}"
