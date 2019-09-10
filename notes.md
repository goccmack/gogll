# Notes on development of gogll v1

# git
This separate git worktree was created inside the gogll git repo by:

```
git branch gogll105
git worktree add ../gogll105 gogll105
```

Merge back:
```
cd ../gogll
git merge gogll105
rm -rf ../gogll105
git worktree prune
```