# Notes on development of gogll v1

# git
This separate git worktree was created inside the gogll git repo by:

```
git branch gogll1
git worktree add ../gogll1 gogll1 
```

Merge back:
```
cd ../gogll
git merge gogll1
rm -rf ../gogll1
git worktree prune
```